package uk.co.eelpieconsulting.monitoring;

import java.util.Collections;
import java.util.List;
import java.util.Map;

import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Metric;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;

@Component
public class MetricsDAO {

	private Map<String, Metric> metrics;

	public MetricsDAO() {
		this.metrics = Maps.newConcurrentMap();
	}
	
	public void registerMetric(Metric metric) {
		metrics.put(metric.getName(), metric);
	}
	
	public List<Metric> getMetrics() {
		List<Metric> allMetrics = Lists.newArrayList(metrics.values());
		Collections.sort(allMetrics);	// TODO guava way todo this?
		return allMetrics;	// TODO immutable copy please
	}
	
	public Metric getByName(String metricName) {
		return metrics.get(metricName);
	}

	public boolean isKnownMetricName(String metricName) {
		return metrics.containsKey(metricName);
	}
	
}
