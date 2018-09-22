package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

import java.math.BigDecimal;
import java.math.RoundingMode;

public class Scaled implements Transform {

    private BigDecimal scale;

    public Scaled(double scale) {
        this.scale = BigDecimal.valueOf(scale);
    }

    @Override
    public String getName() {
        return "x " + scale;
    }

    @Override
    public String transform(Metric metric) {
        BigDecimal metricValue = BigDecimal.valueOf(Double.parseDouble(metric.getLastValue()));
        BigDecimal scaled = metricValue.multiply(scale).setScale(0, RoundingMode.HALF_UP);
        return scaled.toPlainString();
    }

}