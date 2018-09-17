package uk.co.eelpieconsulting.monitoring.model;

import uk.co.eelpieconsulting.monitoring.model.transforms.Transform;

import java.io.Serializable;

public class MetricRouting implements Serializable {

  private static final long serialVersionUID = 1L;

  private String gauge;
  private String metricName;
  private Transform transform;

  public MetricRouting() {
  }

  public MetricRouting(String gauge, String metricName, Transform transform) {
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

  public Transform getTransform() {
    return transform;
  }

  @Override
  public String toString() {
    return "MetricRouting{" +
            "gauge='" + gauge + '\'' +
            ", metricName='" + metricName + '\'' +
            ", transform=" + transform +
            '}';
  }

}
