package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.Message;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;
import uk.co.eelpieconsulting.monitoring.model.MetricType;

@Component
public class MetricsListener {
	
	private final static Logger log = Logger.getLogger(MetricRouting.class);
	
	private final MqttConnectionFactory mqttConnectionFactory;
	
	@Autowired
	public MetricsListener(MqttConnectionFactory mqttConnectionFactory, 
			MetricsDAO metricsDAO,
			RoutingDAO routingDAO,
			MetricPublisher metricPublisher) throws Exception {
		
		this.mqttConnectionFactory = mqttConnectionFactory;
		new Thread(new Listener(metricsDAO, routingDAO, metricPublisher)).start();
	}

	private class Listener implements Runnable {
		
		private final MetricsDAO metricsDAO;
		private final RoutingDAO routingDAO;
		private final MetricPublisher metricPublisher;
		
		public Listener(MetricsDAO metricsDAO, RoutingDAO routingDAO, MetricPublisher metricPublisher) {
			this.metricsDAO = metricsDAO;
			this.routingDAO = routingDAO;
			this.metricPublisher = metricPublisher;	      
		}

		@Override
		public void run() {
			try {
				log.info("Starting metrics listener");
		        final BlockingConnection connection = mqttConnectionFactory.subscribeToMetricsTopic();
		        while (true) {
		        	processNextMessageFrom(connection);
		        }
		        
			} catch (Exception e) {
				log.error(e);
			}
		}

		private void processNextMessageFrom(BlockingConnection connection) {
			try {
				Message message = connection.receive();
	        	byte[] payload = message.getPayload();
	        			        	
	        	String metricMessage = new String(payload, "UTF-8");
	        	log.debug("Got metric message: " + metricMessage);
	        	String[] fields = metricMessage.split(":");
	        	String lastValue = fields[1];
	        	MetricType type = lastValue.equals("true") || lastValue.equals("false") ? MetricType.BOOLEAN : MetricType.NUMBER;
				Metric metric = new Metric(fields[0], type, lastValue);
	        	metricsDAO.registerMetric(metric);
	        	
	        	final boolean isRoutedMetric = routingDAO.isRoutedMetric(metric);
				if (isRoutedMetric) {
					for (MetricRouting routing : routingDAO.getRoutingsForMetric(metric.getName())) {
						metricPublisher.publishOntoGaugesChannel(routing.getGauge(), metric, routing.getScale());						
					}
	        	}
				
				message.ack();
				
			} catch (Exception e) {
				log.error(e);
			}
		}
		
	}
	
}
