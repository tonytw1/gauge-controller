package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public class FromBoolean implements Transform {
  @Override
  public String getName() {
    return "invert boolean";
  }

  @Override
  public String transform(Metric metric) {
    return booleanToIntString(Boolean.parseBoolean(metric.getLastValue()));
  }

  private String booleanToIntString(boolean b) {
    return Integer.toString(b ? 0 : 1);
  }

}
