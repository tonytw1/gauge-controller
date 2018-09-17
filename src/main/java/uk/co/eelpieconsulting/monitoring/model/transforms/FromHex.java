package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public class FromHex implements Transform {
  @Override
  public String getName() {
    return "from hex";
  }

  @Override
  public String transform(Metric metric) {
    return parseHexNumber(metric.getLastValue());
  }

  private String parseHexNumber(String i) {
    return Long.toString(Long.parseLong(i, 16));
  }

}
