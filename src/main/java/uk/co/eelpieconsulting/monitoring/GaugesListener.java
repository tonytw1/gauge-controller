package uk.co.eelpieconsulting.monitoring;

import org.apache.log4j.Logger;
import org.fusesource.mqtt.client.BlockingConnection;
import org.fusesource.mqtt.client.Message;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Gauge;

@Component
public class GaugesListener {

  private final static Logger log = Logger.getLogger(GaugesListener.class);

  private final MqttConnectionFactory mqttConnectionFactory;
  private final GaugeAnnouncementParser gaugeAnnouncementParser;

  @Autowired
  public GaugesListener(GaugeDAO gaugeDAO, MqttConnectionFactory mqttConnectionFactory, GaugeAnnouncementParser gaugeAnnouncementParser) throws Exception {
    this.mqttConnectionFactory = mqttConnectionFactory;
    this.gaugeAnnouncementParser = gaugeAnnouncementParser;
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
          String messageString = new String(payload, "UTF-8").trim();
          log.debug("Got gauges message: '" + messageString + "'");
          Gauge parsedGauge = gaugeAnnouncementParser.parse(messageString);
          if (parsedGauge != null) {
            log.debug("Registering gauge: " + parsedGauge);
            gaugeDAO.registerGauge(parsedGauge);
          }
        }

      } catch (Exception e) {
        log.error(e);
      }
    }
  }

}
