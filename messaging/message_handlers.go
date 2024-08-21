package messaging

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tonytw1/gauges/model"
	"github.com/tonytw1/gauges/routing"
	"github.com/tonytw1/gauges/transforms"
	"log"
	"strconv"
	"strings"
	"time"
)

func GaugesMessageHandler(gaugesTable *routing.GaugesTable) func(mqtt.Client, mqtt.Message) {
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
		gaugesTable.AddGauge(gauge)
	}

	return gaugesMessageHandler
}

func MetricsMessageHandler(metricsTable *routing.MetricsTable, routingTable *routing.RoutesTable, gaugesTopic string, metricsTopic string) func(mqtt.Client, mqtt.Message) {
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
		metricsTable.AddMetrics(metric)

		// Route metrics
		routes, ok := routingTable.GetRoutesForMetric(metric.Name)
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

	return metricsMessageHandler
}

func publish(c mqtt.Client, topic string, message string) {
	token := c.Publish(topic, 0, false, message)
	token.WaitTimeout(time.Second * 1)
}
