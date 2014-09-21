package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.Message;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Gauge;
import uk.co.eelpieconsulting.monitoring.model.GaugeType;

@Component
public class GaugesListener {
	
	private final static Logger log = Logger.getLogger(GaugesListener.class);
	
	private static final String GAUGE_PREFIX = "gauge:";
	
	private final MqttConnectionFactory mqttConnectionFactory;
	
	@Autowired
	public GaugesListener(GaugeDAO gaugeDAO, MqttConnectionFactory mqttConnectionFactory) throws Exception {	   
		this.mqttConnectionFactory = mqttConnectionFactory;		
		new Thread(new Listener(gaugeDAO)).start();
	}

	private class Listener implements Runnable {
		
		private final GaugeDAO gaugeDAO;

		public Listener(GaugeDAO gaugeDAO) {
			this.gaugeDAO = gaugeDAO;
		}

		@Override
		public void run() {
			try {
				BlockingConnection connection = mqttConnectionFactory.subscribeToGaugesTopic();
		        
		        while (true) {        	
		        	final Message message = connection.receive();
		        	byte[] payload = message.getPayload();
		        	String messageString = new String(payload, "UTF-8");
		        	log.debug("Got gauges message: " + messageString);
		        	try {
		        		if (messageString.startsWith(GAUGE_PREFIX)) {
		        			String gaugeDescripition = messageString.split(GAUGE_PREFIX)[1];
		        			gaugeDAO.registerGauge(parseGaugeDescription(gaugeDescripition));
		        		}		        		
		        	} catch (Exception e) {
		        		log.error(e);
		        	}
		        }
		        
			} catch (Exception e) {
				log.error(e);
			}
		}

		private Gauge parseGaugeDescription(String gaugeDescripition) {
			return new Gauge(gaugeDescripition.split(",")[0], GaugeType.valueOf(gaugeDescripition.split(",")[1].toUpperCase()));
		}
		
	}
}
