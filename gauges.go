package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
	"github.com/tonytw1/gauges/messaging"
	"github.com/tonytw1/gauges/model"
	"github.com/tonytw1/gauges/persistence"
	"github.com/tonytw1/gauges/transforms"
	"github.com/tonytw1/gauges/views"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
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

	gauges := sync.Map{}
	metrics := sync.Map{}

	gaugesMessageHandler := messaging.GaugesMessageHandler(&gauges)
	metricsMessageHandler := messaging.MetricsMessageHandler(&metrics, &routingTable, gaugesTopic, metricsTopic)

	log.Print("Connecting to MQTT")
	mqttClient := messaging.SetupMqttClient(configuration.MqttUrl, "gauges-ui-"+uuid.New().String(),
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
