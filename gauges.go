package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"github.com/tonytw1/gauges/model"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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
		load, ok := routingTable.Load(metric.Name)
		if ok {
			route := load.(model.Route)
			log.Print("Routing " + metric.Name + " to " + route.ToGauge)
			publish(client, "test", payload)
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
		name := split[0]
		value := split[1]
		gauge := model.Gauge{Name: name, MaxValue: value}
		gauges.Store(name, gauge)
	}

	log.Print("Connecting to MQTT")
	mqttClient := setupMqttClient(configuration.MqttUrl, "gauges",
		"metrics", metricsMessageHandler,
		"gauges", gaugesMessageHandler)

	defer mqttClient.Disconnect(250)

	getHomepage := func(w http.ResponseWriter, r *http.Request) {
		var gaugesCount = 0
		gauges.Range(func(k, v interface{}) bool {
			gaugesCount++
			return true
		})
		var metricsCount = 0
		metrics.Range(func(k, v interface{}) bool {
			metricsCount++
			return true
		})

		output := "Gauges: " + fmt.Sprint(gaugesCount)
		output += "Metrics: " + fmt.Sprint(metricsCount)
		io.WriteString(w, output)
	}

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

		w.Header().Set("Access-Control-Allow-Origin", "*")
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

		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, string(asJson))
	}

	getRoutes := func(w http.ResponseWriter, r *http.Request) {
		asJson := routesAsJson(routingTable)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		io.WriteString(w, string(asJson))
	}

	getRoute := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // TODO null check
		route, ok := routingTable.Load(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		asJson, err := json.Marshal(route)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		io.WriteString(w, string(asJson))
	}

	type routeRequest struct {
		Metric string
		Gauge  string
	}

	optionsRoutes := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		io.WriteString(w, string("ok"))
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
		routingTable.Store(id, model.Route{
			Id:         id,
			FromMetric: rr.Metric,
			ToGauge:    rr.Gauge,
		})

		asJson := routesAsJson(routingTable)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		io.WriteString(w, string(asJson))
	}

	log.Print("Starting HTTP server")
	r := mux.NewRouter()
	r.HandleFunc("/", getHomepage)
	r.HandleFunc("/gauges", getGauges)
	r.HandleFunc("/metrics", getMetrics)
	r.HandleFunc("/routes", getRoutes).Methods("GET")
	r.HandleFunc("/routes/{id}", getRoute).Methods("GET")
	r.HandleFunc("/routes", optionsRoutes).Methods("OPTIONS")
	r.HandleFunc("/routes", postRoutes).Methods("POST")
	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
	log.Print("Done")
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
