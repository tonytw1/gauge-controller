package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"github.com/tonytw1/gauges/model"
	"github.com/tonytw1/gauges/transforms"
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

func main() {
	type Configuration struct {
		MqttUrl string
	}
	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		panic(err)
	}

	var routes = sync.Map{}
	var routingTable = sync.Map{}

	var metrics = sync.Map{}

	metricsMessageHandler := func(client mqtt.Client, message mqtt.Message) {
		payload := strings.TrimSpace(string(message.Payload()))
		//log.Print("Received: " + payload + " on " + message.Topic())

		split := strings.Split(payload, ":")
		if len(split) != 2 {
			log.Print("Rejected malformed metrics message: '" + payload + "'")
			return
		}

		name := split[0]
		value := split[1]
		metric := model.Metric{Name: name, Value: value}
		metrics.Store(name, metric)

		// Route metrics
		route, ok := routingTable.Load(metric.Name)
		if ok {
			route := route.(model.Route)
			log.Print("Routing " + metric.Name + " to " + route.ToGauge)
			transform, ok := transforms.Transforms()[route.Transform]
			if ok {
				transformedValue, err := transform(value)
				if err != nil {
					gaugesMessage := route.ToGauge + ":" + strconv.Itoa(transformedValue)
					log.Print("Sending gauge message: " + gaugesMessage)
					publish(client, "gauges", gaugesMessage)
				} else {
					log.Print("Transform error: " + err.Error())
				}
			} else {
				log.Print("Unknown transform: " + route.Transform)
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
	mqttClient := setupMqttClient(configuration.MqttUrl, "gauges",
		"metrics/#", metricsMessageHandler,
		"gauges", gaugesMessageHandler)

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
		asJson := routesAsJson(routingTable)

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
		routingTable.Delete(route.(model.Route).FromMetric)

		asJson := routesAsJson(routes)
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

		id := uuid.New().String()
		route := model.Route{
			Id:         id,
			FromMetric: rr.Metric,
			ToGauge:    rr.Gauge,
			Transform:  rr.Transform, // TODO validate transform
		}
		routes.Store(id, route)
		routingTable.Store(rr.Metric, route)

		asJson := routesAsJson(routes)
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

func setCORSHeadersOn(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func routesAsJson(routingTable sync.Map) []byte {
	var routes = make([]model.Route, 0)
	routingTable.Range(func(k, v interface{}) bool {
		routes = append(routes, v.(model.Route))
		return true
	})
	sort.Slice(routes, func(i, j int) bool {
		return strings.Compare(routes[i].FromMetric, routes[j].FromMetric) < 0
	})
	asJson, _ := json.Marshal(routes)
	return asJson
}

func setupMqttClient(mqttURL string, clientId string, metricsTopic string, metricsHandler mqtt.MessageHandler,
	gaugesTopic string, gaugesHandler mqtt.MessageHandler) mqtt.Client {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)

	var subscribeToTopic mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Print("Connected to " + mqttURL)
		log.Print("Subscribing to " + metricsTopic)
		client.Subscribe(metricsTopic, 0, metricsHandler)
		log.Print("Subscribing to " + gaugesTopic)
		client.Subscribe(gaugesTopic, 0, gaugesHandler)
	}
	var logConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Print("Connection lost")
	}
	var logReconnecting mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Print("Reconnecting")
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttURL)
	opts.SetClientID(clientId)
	opts.SetOnConnectHandler(subscribeToTopic)
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
