package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public class InvertedBooleanInt implements Transform {
  @Override
  public String getName() {
    return "inverted boolean int";
  }

  @Override
  public String transform(Metric metric) {
    int i = Integer.parseInt(metric.getLastValue());
    return Boolean.toString(i == 0 ? true : false);
  }
  
}
