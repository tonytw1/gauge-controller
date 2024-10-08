package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"github.com/tonytw1/gauges/messaging"
	"github.com/tonytw1/gauges/model"
	"github.com/tonytw1/gauges/persistence"
	"github.com/tonytw1/gauges/routing"
	"github.com/tonytw1/gauges/transforms"
	"github.com/tonytw1/gauges/views"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
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

	routePersistence := persistence.NewS3RoutePersistence(configuration.Bucket, "routes.json")

	routesTable := routing.NewRoutesTable()
	gaugesTable := routing.NewGaugesTable()
	metricsTable := routing.NewMetricsTable()

	// Reload persisted routes
	persistedRoutes := routePersistence.LoadPersistedRoutes()
	for _, route := range persistedRoutes {
		routesTable.AddRoute(route)
	}

	gaugesTopic := "gauges"
	metricsTopic := "metrics"

	log.Print("Connecting to MQTT")
	mqttClient := messaging.SetupMqttClient(configuration.MqttUrl,
		metricsTopic,
		gaugesTopic,
		&gaugesTable, &metricsTable, &routesTable)
	defer mqttClient.Disconnect(250)

	getGauges := func(w http.ResponseWriter, r *http.Request) {
		var gs = gaugesTable.AllGauges()
		sort.Slice(gs, func(i, j int) bool {
			return strings.Compare(gs[i].Name, gs[j].Name) < 0
		})
		asJson, _ := json.Marshal(gs)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getMetrics := func(w http.ResponseWriter, r *http.Request) {
		ms := metricsTable.AllMetrics()
		sort.Slice(ms, func(i, j int) bool {
			return strings.Compare(ms[i].Name, ms[j].Name) < 0
		})
		asJson, _ := json.Marshal(ms)

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getRoutes := func(w http.ResponseWriter, r *http.Request) {
		asJson := views.RoutesAsJson(routesTable.AllRoutes())

		setCORSHeadersOn(w)
		io.WriteString(w, string(asJson))
	}

	getRoute := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // TODO null check
		route, ok := routesTable.GetRoute(id)
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
		id := vars["id"]
		if id == "" {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}

		route, ok := routesTable.GetRoute(id)
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		routesTable.Delete(route)

		asJson := views.RoutesAsJson(routesTable.AllRoutes())

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

		metric, ok := metricsTable.GetMetric(rr.Metric)
		if !ok {
			http.Error(w, "Invalid metric name", http.StatusBadRequest)
			return
		}
		transform, ok := transforms.GetTransformByName(rr.Transform)
		if !ok {
			http.Error(w, "Invalid transform name", http.StatusBadRequest)
			return
		}
		gauge, ok := gaugesTable.GetGauge(rr.Gauge)
		if !ok {
			http.Error(w, "Invalid gauge name", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()
		route := model.Route{
			Id:         id,
			FromMetric: metric.Name,
			ToGauge:    gauge.Name,
			Transform:  transform.Name,
		}
		routesTable.AddRoute(route)

		asJson := views.RoutesAsJson(routesTable.AllRoutes())

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

func setCORSHeadersOn(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}
