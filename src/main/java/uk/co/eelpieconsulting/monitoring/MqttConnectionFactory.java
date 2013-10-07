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
	
	private final static Logger log = Logger.getLogger(MqttConnectionFactory.class);
	
	private final String metricsHost;
	private final String metricsTopic;
	private final String gaugesHost;
	private final String gaugesTopic;

	@Autowired
	public MqttConnectionFactory(
			@Value(value = "#{config['mqtt.metrics.host']}") String metricsHost,
			@Value(value = "#{config['mqtt.metrics.topic']}") String metricsTopic,
			@Value(value = "#{config['mqtt.gauges.host']}") String gaugesHost,
			@Value(value = "#{config['mqtt.gauges.topic']}") String gaugesTopic) {
		this.metricsHost = metricsHost;
		this.metricsTopic = metricsTopic;
		this.gaugesHost = gaugesHost;
		this.gaugesTopic = gaugesTopic;
	}
	
	public BlockingConnection connectToGaugesHost() throws URISyntaxException, Exception {
		log.info("Connecting to gauges host: " + gaugesHost);
		MQTT mqtt = new MQTT();
		mqtt.setHost(gaugesHost, 1883);
		BlockingConnection connection = mqtt.blockingConnection();
		connection.connect();
		return connection;
	}
	
	public BlockingConnection subscribeToMetricsTopic() throws URISyntaxException, Exception {
		BlockingConnection connection = connectToMetricsHost();
		log.info("Subscribing to topic '" + metricsTopic + "' on host '" + metricsHost + "'");
		connection.subscribe(new Topic[] { new Topic(metricsTopic, QoS.AT_LEAST_ONCE) });
		return connection;
	}

	public BlockingConnection subscribeToGaugesTopic() throws Exception {
		BlockingConnection connection = connectToGaugesHost();
		log.info("Subscribing to topic: " + gaugesTopic);
		connection.subscribe(new Topic[] { new Topic(gaugesTopic, QoS.AT_LEAST_ONCE) });
		return connection;
	}
	
	private BlockingConnection connectToMetricsHost() throws URISyntaxException, Exception {
		log.info("Connecting to mertics host: " + metricsHost);
		MQTT mqtt = new MQTT();
		mqtt.setHost(metricsHost, 1883);
		BlockingConnection connection = mqtt.blockingConnection();
		connection.connect();
		return connection;
	}

}
