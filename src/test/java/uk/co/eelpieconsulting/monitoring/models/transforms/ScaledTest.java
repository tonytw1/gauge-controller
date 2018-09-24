package uk.co.eelpieconsulting.monitoring.models.transforms;

import com.google.common.collect.Lists;
import org.joda.time.DateTime;
import org.junit.Test;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricType;
import uk.co.eelpieconsulting.monitoring.model.transforms.Scaled;

import static org.junit.Assert.assertEquals;

public class ScaledTest {

    @Test
    public void shouldPassThroughUnscaledNumbers() {
        Metric metric = new Metric("test", MetricType.NUMBER, "23173", DateTime.now(), Lists.newArrayList());

        String scaled = new Scaled(1).transform(metric);

        assertEquals("23173", scaled);
    }

    @Test
    public void scaledValuesShouldBeCorrectlyRounded() {
        Metric metric = new Metric("test", MetricType.NUMBER, "23173", DateTime.now(), Lists.newArrayList());

        String scaled = new Scaled(0.01).transform(metric);

        assertEquals("232", scaled);
    }

    @Test
    public void upScaledValuesShouldBeCorrectlyRounded() {
        Metric metric = new Metric("test", MetricType.NUMBER, "-23173", DateTime.now(), Lists.newArrayList());

        String scaled = new Scaled(10).transform(metric);

        assertEquals("-231730", scaled);
    }

}