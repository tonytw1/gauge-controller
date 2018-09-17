package uk.co.eelpieconsulting.monitoring.model;

import java.io.Serializable;

public class PersistedMetricRouting implements Serializable {

  private static final long serialVersionUID = 1L;

  private String gauge;
  private String metricName;
  private String transform;

  public PersistedMetricRouting() {
  }

  public PersistedMetricRouting(String gauge, String metricName, String transform) {
    this.gauge = gauge;
    this.metricName = metricName;
    this.transform = transform;
  }

  public String getGauge() {
    return gauge;
  }

  public String getMetricName() {
    return metricName;
  }

  public String getTransform() {
    return transform;
  }

}
