package uk.co.eelpieconsulting.monitoring;

import java.net.URISyntaxException;

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
	private static final int MQTT_PORT = 32183;

	private final String metricsHost;
	private final String metricsTopic;
	private final String gaugesHost;
	private final String gaugesTopic;

	@Autowired
	public MqttConnectionFactory(
			@Value("${mqtt.metrics.host}") String metricsHost,
			@Value("${mqtt.metrics.topic}") String metricsTopic,
			@Value("${mqtt.gauges.host}") String gaugesHost,
			@Value("${mqtt.gauges.topic}") String gaugesTopic) {
		this.metricsHost = metricsHost;
		this.metricsTopic = metricsTopic;
		this.gaugesHost = gaugesHost;
		this.gaugesTopic = gaugesTopic;
	}
	
	public BlockingConnection connectToGaugesHost() throws URISyntaxException, Exception {
		log.info("Connecting to gauges host: " + gaugesHost);
		MQTT mqtt = new MQTT();
		mqtt.setHost(gaugesHost, MQTT_PORT);
		BlockingConnection connection = mqtt.blockingConnection();
		connection.connect();
		return connection;
	}
	
	public BlockingConnection subscribeToMetricsTopic() throws URISyntaxException, Exception {
		BlockingConnection connection = connectToMetricsHost();
		log.info("Subscribing to topic '" + metricsTopic + "' on host '" + metricsHost + "'");
		connection.subscribe(new Topic[] { new Topic(metricsTopic, QoS.AT_MOST_ONCE) });
		return connection;
	}

	public BlockingConnection subscribeToGaugesTopic() throws Exception {
		BlockingConnection connection = connectToGaugesHost();
		log.info("Subscribing to topic: " + gaugesTopic);
		connection.subscribe(new Topic[] { new Topic(gaugesTopic, QoS.AT_MOST_ONCE) });
		return connection;
	}
	
	private BlockingConnection connectToMetricsHost() throws URISyntaxException, Exception {
		log.info("Connecting to metrics host: " + metricsHost);
		MQTT mqtt = new MQTT();
		mqtt.setHost(metricsHost, MQTT_PORT);
		BlockingConnection connection = mqtt.blockingConnection();
		connection.connect();
		return connection;
	}

}
