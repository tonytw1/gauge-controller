package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"github.com/tonytw1/gauges/model"
	"github.com/tonytw1/gauges/persistence"
	"github.com/tonytw1/gauges/transforms"
	"github.com/tonytw1/gauges/views"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	type Configuration struct {
		MqttUrl string
		Bucket  string
	}
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	// Setup S3 client
	region := "eu-west-2"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
	routePersistence := persistence.S3RoutePersistence{S3Client: s3Client, Bucket: configuration.Bucket, Key: "routes.json"}

	var routes = sync.Map{}
	var routingTable = sync.Map{}

	// Reload persisted routes
	persistedRoutes := routePersistence.LoadPersistedRoutes()
	for _, route := range persistedRoutes {
		addRoute(&routes, route, &routingTable)
	}

	gaugesTopic := "gauges"
	metricsTopic := "metrics"

	var metrics = sync.Map{}
	metricsMessageHandler := func(client mqtt.Client, message mqtt.Message) {
		payload := strings.TrimSpace(string(message.Payload()))
		//log.Print("Received: " + payload + " on " + message.Topic())

		topic := message.Topic()
		payloadFields := strings.Split(payload, ":")
		if len(payloadFields) != 2 {
			log.Print("Rejected malformed metrics message: '" + payload + "'")
			return
		}

		subtopic := strings.TrimPrefix(topic, metricsTopic+"/")
		name := subtopic + "/" + payloadFields[0]
		value := payloadFields[1]
		metric := model.Metric{Name: name, Value: value}
		metrics.Store(name, metric)

		// Route metrics
		routes, ok := routingTable.Load(metric.Name)
		if ok {
			routes := routes.([]model.Route)
			if ok {
				for _, route := range routes {
					log.Print("Routing " + metric.Name + " to " + route.ToGauge)
					transform, ok := transforms.GetTransformByName(route.Transform)
					if ok {
						transformedValue, err := transform.Transform(value)
						if err != nil {
							log.Print("Transform error: " + err.Error())
							return
						}
						gaugesMessage := route.ToGauge + ":" + strconv.Itoa(transformedValue)
						log.Print("Sending gauge message: " + gaugesMessage)
						publish(client, gaugesTopic+"/signals", gaugesMessage)

					} else {
						log.Print("Unknown transform: " + route.Transform)
					}
				}
			}
		}
	}

	var gauges = sync.Map{}
	gaugesMessageHandler := func(client mqtt.Client, message mqtt.Message) {
		payload := strings.TrimSpace(string(message.Payload()))

		split := strings.Split(payload, ":")
		if len(split) != 2 {
			log.Print("Rejected malformed gauges message: '" + payload + "'")
			return
		}
		description := split[1]
		fields := strings.Split(description, "[")
		if len(fields) != 2 {
			log.Print("Rejected malformed gauges message: '" + payload + "'")
			return
		}

		name := fields[0]
		value := strings.TrimSuffix(fields[1], "]")

		gauge := model.Gauge{Name: name, MaxValue: value}
		gauges.Store(name, gauge)
	}

	log.Print("Connecting to MQTT")
	mqttClient := setupMqttClient(configuration.MqttUrl, "gauges-ui-"+uuid.New().String(),
		metricsTopic+"/#", metricsMessageHandler,
		gaugesTopic+"/announcements", gaugesMessageHandler)

	defer mqttClient.Disconnect(250)

	getGauges := func(w http.ResponseWriter, r *http.Request) {
		var gs = make([]model.Gauge, 0)
		gauges.Range(func(k, v interface{}) bool {
			gs = append(gs, v.(model.Gauge))
			return true
		})
		sort.Slice(gs, func(i, j int) bool {
			return strings.Compare(gs[i].Name, gs[j].Name) < 0
		})
		asJson, _ := json.Marshal(gs)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getMetrics := func(w http.ResponseWriter, r *http.Request) {
		var ms = make([]model.Metric, 0)
		metrics.Range(func(k, v interface{}) bool {
			ms = append(ms, v.(model.Metric))
			return true
		})
		sort.Slice(ms, func(i, j int) bool {
			return strings.Compare(ms[i].Name, ms[j].Name) < 0
		})
		asJson, _ := json.Marshal(ms)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getRoutes := func(w http.ResponseWriter, r *http.Request) {
		asJson := views.RoutesAsJson(routes)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getRoute := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // TODO null check
		route, ok := routes.Load(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		asJson, err := json.Marshal(route)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	optionsRoute := func(w http.ResponseWriter, r *http.Request) {
		setCORSHeadersOn(w)
		io.WriteString(w, "ok")
	}

	deleteRoute := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // TODO null check
		route, ok := routes.Load(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}

		routes.Delete(route.(model.Route).Id)
		// Update routing table for effected metric
		effectedMetric := route.(model.Route).FromMetric
		effectedMetricRoutes, ok := routingTable.Load(effectedMetric)
		if ok {
			// Filter out the route that was deleted
			filtered := make([]model.Route, 0)
			for _, route := range effectedMetricRoutes.([]model.Route) {
				if route.Id != id {
					filtered = append(filtered, route)
				}
			}
			routingTable.Store(effectedMetric, filtered)
		}

		asJson := views.RoutesAsJson(routes)

		_, err = routePersistence.PersistRoutes(asJson)
		if err != nil {
			http.Error(w, "Failed to store updated routes", http.StatusInternalServerError)
			return
		}

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	type routeRequest struct {
		Metric    string
		Gauge     string
		Transform string
	}

	optionsRoutes := func(w http.ResponseWriter, r *http.Request) {
		setCORSHeadersOn(w)
		io.WriteString(w, "ok")
	}

	postRoutes := func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var rr routeRequest
		err := decoder.Decode(&rr)
		if err != nil {
			log.Print("Decode error", err)
			panic(err)
		}

		metric, ok := metrics.Load(rr.Metric)
		if !ok {
			http.Error(w, "Invalid metric name", http.StatusBadRequest)
			return
		}
		transform, ok := transforms.GetTransformByName(rr.Transform)
		if !ok {
			http.Error(w, "Invalid transform name", http.StatusBadRequest)
			return
		}
		gauge, ok := gauges.Load(rr.Gauge)
		if !ok {
			http.Error(w, "Invalid gauge name", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()
		route := model.Route{
			Id:         id,
			FromMetric: metric.(model.Metric).Name,
			ToGauge:    gauge.(model.Gauge).Name,
			Transform:  transform.Name,
		}
		addRoute(&routes, route, &routingTable)

		asJson := views.RoutesAsJson(routes)

		_, err = routePersistence.PersistRoutes(asJson)
		if err != nil {
			http.Error(w, "Failed to store updated routes", http.StatusInternalServerError)
			return
		}

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getTransforms := func(w http.ResponseWriter, r *http.Request) {
		type DisplayTransform struct {
			Name string
		}

		var displayTransforms = make([]DisplayTransform, 0)
		for t := range transforms.Transforms() {
			displayTransforms = append(displayTransforms, DisplayTransform{Name: t})
		}
		sort.Slice(displayTransforms, func(i, j int) bool {
			return strings.Compare(displayTransforms[i].Name, displayTransforms[j].Name) < 0
		})
		asJson, _ := json.Marshal(displayTransforms)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	log.Print("Starting HTTP server")
	r := mux.NewRouter()
	r.HandleFunc("/gauges", getGauges)
	r.HandleFunc("/metrics", getMetrics)
	r.HandleFunc("/routes", getRoutes).Methods("GET")
	r.HandleFunc("/routes/{id}", getRoute).Methods("GET")
	r.HandleFunc("/routes/{id}", optionsRoute).Methods("OPTIONS")
	r.HandleFunc("/routes/{id}", deleteRoute).Methods("DELETE")
	r.HandleFunc("/routes", optionsRoutes).Methods("OPTIONS")
	r.HandleFunc("/routes", postRoutes).Methods("POST")
	r.HandleFunc("/transforms", getTransforms).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("client/dist")))

	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
	log.Print("Done")
}

func addRoute(routes *sync.Map, route model.Route, routingTable *sync.Map) {
	routes.Store(route.Id, route)
	// Update routing table for effected metric
	effectedRoutes, ok := routingTable.Load(route.FromMetric)
	if ok {
		updated := append(effectedRoutes.([]model.Route), route)
		routingTable.Store(route.FromMetric, updated)
	} else {
		routingTable.Store(route.FromMetric, []model.Route{route})
	}
}

func setCORSHeadersOn(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func setupMqttClient(mqttURL string, clientId string, metricsTopic string, metricsHandler mqtt.MessageHandler,
	gaugesAnnouncementsTopic string, gaugesHandler mqtt.MessageHandler) mqtt.Client {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)

	var subscribeToTopics mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Print("Connected to " + mqttURL)
		log.Print("Subscribing to " + metricsTopic)
		client.Subscribe(metricsTopic, 0, metricsHandler)
		log.Print("Subscribing to " + gaugesAnnouncementsTopic)
		client.Subscribe(gaugesAnnouncementsTopic, 0, gaugesHandler)
	}
	var logConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Print("Connection lost")
	}
	var logReconnecting mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Print("Reconnecting")
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	opts.SetClientID(clientId)
	opts.SetOnConnectHandler(subscribeToTopics)
	opts.SetConnectionLostHandler(logConnectionLost)
	opts.SetReconnectingHandler(logReconnecting)

	log.Print("Connecting to: ", mqttURL)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func publish(c mqtt.Client, topic string, message string) {
	token := c.Publish(topic, 0, false, message)
	token.WaitTimeout(time.Second * 1)
}
