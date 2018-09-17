package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public class Scaled implements Transform {

  private double scale;

  public Scaled(double scale) {
    this.scale = scale;
  }

  @Override
  public String getName() {
    return "x " + scale;
  }

  @Override
  public String transform(Metric metric) {
    return Double.toString(Double.parseDouble(metric.getLastValue()) * scale);
  }

}