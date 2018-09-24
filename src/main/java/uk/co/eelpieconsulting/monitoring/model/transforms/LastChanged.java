package uk.co.eelpieconsulting.monitoring.model.transforms;

import org.apache.log4j.Logger;
import org.joda.time.DateTime;
import uk.co.eelpieconsulting.monitoring.model.Metric;

public class LastChanged implements Transform {

    private final static Logger log = Logger.getLogger(LastChanged.class);

    @Override
    public String getName() {
        return "lastChangeDuration";
    }

    @Override
    public String transform(Metric metric) {
        double i = 0.0;
        if (metric.getChanges().size() >= 2) {
            DateTime a = metric.getChanges().get(metric.getChanges().size() - 1);
            DateTime b = metric.getChanges().get(metric.getChanges().size() - 2);
            double s = (a.getMillis() - b.getMillis()) * 0.001;
            // log.info("Delta: " + s);

            int kiloWattHour = 60 * 60 * 1000;
            int ticksPerKiloWattHour = 800;
            int joulesPerTick = kiloWattHour / ticksPerKiloWattHour;

            double ticksPerSecond = 1 / s;
            //log.info("Ticks per second: " + ticksPerSecond);

            i = ticksPerSecond * joulesPerTick;
        }
        return Double.toString(i);
    }

}
