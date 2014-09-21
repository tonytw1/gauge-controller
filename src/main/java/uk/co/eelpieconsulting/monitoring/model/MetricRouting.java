package uk.co.eelpieconsulting.monitoring.model;

import java.io.Serializable;

public class MetricRouting implements Serializable {
	
	private static final long serialVersionUID = 1L;
	
	private String gauge;
	private String metricName;
	private double scale;
	
	public MetricRouting() {
	}
	
	public MetricRouting(String gauge, String metricName, double scale) {
		this.gauge = gauge;
		this.metricName = metricName;
		this.scale = scale;
	}
	
	public String getGauge() {
		return gauge;
	}
	
	public String getMetricName() {
		return metricName;
	}

	public double getScale() {
		return scale;
	}

	@Override
	public String toString() {
		return "MetricRouting [gauge=" + gauge + ", metricName=" + metricName + ", scale=" + scale + "]";
	}

}
