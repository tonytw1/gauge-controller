package messaging

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/tonytw1/gauges/routing"
	"log"
	"os"
)

func SetupMqttClient(mqttURL string, metricsTopic string, gaugesAnnouncementsTopic string, gaugesTable *routing.GaugesTable,
	metricsTable *routing.MetricsTable, routesTable *routing.RoutesTable, gaugesTopic string) mqtt.Client {

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)

	var subscribeToTopics mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Print("Connected to " + mqttURL)
		log.Print("Subscribing to " + gaugesAnnouncementsTopic)
		client.Subscribe(gaugesAnnouncementsTopic, 0, GaugesMessageHandler(gaugesTable))
		log.Print("Subscribing to " + metricsTopic)
		client.Subscribe(metricsTopic, 0, MetricsMessageHandler(metricsTable, routesTable, gaugesTopic, metricsTopic))
	}
	var logConnectionLost mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Print("Connection lost")
	}
	var logReconnecting mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Print("Reconnecting")
	}

	clientId := "gauges-ui-" + uuid.New().String()
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
