package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.MQTT;
import org.fusesource.mqtt.client.QoS;
import org.fusesource.mqtt.client.Topic;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class MqttConnectionFactory {

    private static final Logger log = Logger.getLogger(MqttConnectionFactory.class);

    private final String metricsHost;
    private final int metricsPort;
    private final String metricsTopic;
    private final String gaugesHost;
    private final int gaugesPort;
    private final String gaugesTopic;

    @Autowired
    public MqttConnectionFactory(
            @Value("${mqtt.metrics.host}") String metricsHost,
            @Value("${mqtt.metrics.port}") int metricsPort,
            @Value("${mqtt.metrics.topic}") String metricsTopic,
            @Value("${mqtt.gauges.host}") String gaugesHost,
            @Value("${mqtt.gauges.port}") int gaugesPort,
            @Value("${mqtt.gauges.topic}") String gaugesTopic) {
        this.metricsHost = metricsHost;
        this.metricsPort = metricsPort;
        this.metricsTopic = metricsTopic;
        this.gaugesHost = gaugesHost;
        this.gaugesPort = gaugesPort;
        this.gaugesTopic = gaugesTopic;
    }

    public BlockingConnection connectToGaugesHost() throws Exception {
        log.info("Connecting to gauges host: " + gaugesHost);
        MQTT mqtt = new MQTT();
        mqtt.setHost(gaugesHost, gaugesPort);
        BlockingConnection connection = mqtt.blockingConnection();
        connection.connect();
        return connection;
    }

    public BlockingConnection subscribeToMetricsTopic() throws Exception {
        BlockingConnection connection = connectToMetricsHost();
        log.info("Subscribing to topic '" + metricsTopic + "' on host '" + metricsHost + "'");
        connection.subscribe(new Topic[]{new Topic(metricsTopic, QoS.AT_MOST_ONCE)});
        return connection;
    }

    public BlockingConnection subscribeToGaugesTopic() throws Exception {
        BlockingConnection connection = connectToGaugesHost();
        log.info("Subscribing to topic: " + gaugesTopic);
        connection.subscribe(new Topic[]{new Topic(gaugesTopic, QoS.AT_MOST_ONCE)});
        return connection;
    }

    private BlockingConnection connectToMetricsHost() throws Exception {
        log.info("Connecting to metrics host: " + metricsHost);
        MQTT mqtt = new MQTT();
        mqtt.setHost(metricsHost, metricsPort);
        mqtt.setClientId("gauge-controller");

        BlockingConnection connection = mqtt.blockingConnection();
        connection.connect();
        return connection;
    }

}
