package uk.co.eelpieconsulting.monitoring.model.transforms;

import org.joda.time.DateTime;
import uk.co.eelpieconsulting.monitoring.model.Metric;

public class LastChanged implements Transform {

    @Override
    public String getName() {
        return "lastChangeDuration";
    }

    @Override
    public String transform(Metric metric) {
        long i = 0;
        if (metric.getChanges().size() >= 2) {
            DateTime a = metric.getChanges().get(metric.getChanges().size() - 1);
            DateTime b = metric.getChanges().get(metric.getChanges().size() - 2);
            i = a.getMillis() - b.getMillis();
        }
        return Long.toString(i);
    }

}
