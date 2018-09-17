package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public class AsString implements Transform {

  @Override
  public String getName() {
    return "toString";
  }

  @Override
  public String transform(Metric metric) {
    return metric.getLastValue().toString();
  }

}
