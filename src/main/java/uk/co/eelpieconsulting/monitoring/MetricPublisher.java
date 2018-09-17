package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.QoS;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.transforms.Transform;

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

  public void publishOntoGaugesChannel(String gauge, Metric metric, Transform transform) throws Exception {
    log.debug("Publishing routed metric: " + metric + " to gauge: " + gauge + " transformed by: " + transform);
    String transformed = transformValue(metric, transform);
    if (transformed != null) {
      final String metricMessage = gauge + ":" + transformed;
      connection.publish(gaugesTopic, metricMessage.getBytes(), QoS.AT_MOST_ONCE, false);
    }
  }

  private String transformValue(Metric metric, Transform transform) {
    try {
      if (metric.getLastValue() != null) {
        return transform.transform(metric);
      }
    } catch (Exception e){
      log.warn("Failed to transform metric");
    }
    return null;
  }

}
