package uk.co.eelpieconsulting.monitoring.model;

public class MetricRouting {
	
	private final String gauge;
	private final String metricName;
	private final double scale;
	
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
