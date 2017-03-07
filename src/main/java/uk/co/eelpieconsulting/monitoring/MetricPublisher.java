package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.QoS;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricType;

import java.net.URISyntaxException;

@Component
public class MetricPublisher {
	
	private final static Logger log = Logger.getLogger(MetricPublisher.class);

	private final BlockingConnection connection;
	
	private final String gaugesTopic;

	@Autowired
	public MetricPublisher(MqttConnectionFactory mqttConnectionFactory,
			@Value(value = "${mqtt.gauges.topic}") String gaugesTopic) throws URISyntaxException, Exception {
		this.connection = mqttConnectionFactory.connectToGaugesHost();	

		this.gaugesTopic = gaugesTopic;
	}
	
	public void publishOntoGaugesChannel(String gauge, Metric metric, double scale) throws Exception {
		log.debug("Publishing routed metric: " + metric + " to gauge: " + gauge + " scaled by: " + scale);			
		final String metricMessage = gauge + ":" + scaleValue(metric, scale);
		connection.publish(gaugesTopic, metricMessage.getBytes(), QoS.AT_MOST_ONCE, false);
	}

	private String scaleValue(Metric metric, double scale) {
		if (metric.getType() == MetricType.BOOLEAN) {
			if (scale == -1) {
				return booleanToIntString(!Boolean.parseBoolean(metric.getLastValue()));
			}
			return booleanToIntString(Boolean.parseBoolean(metric.getLastValue()));
		}
		
		return Double.toString(Double.parseDouble(metric.getLastValue()) * scale);		
	}

	private String booleanToIntString(boolean b) {
		int v = b ? 1 : -1;
		return Integer.toString(v);
	}

}
