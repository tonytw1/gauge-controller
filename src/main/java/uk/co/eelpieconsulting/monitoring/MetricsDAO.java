package uk.co.eelpieconsulting.monitoring;

import java.util.List;
import java.util.concurrent.TimeUnit;

import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Metric;

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;
import com.google.common.collect.ImmutableSortedSet;

@Component
public class MetricsDAO {

	private final Cache<String, Metric> metrics;

	public MetricsDAO() {
		this.metrics = CacheBuilder.newBuilder().
				maximumSize(100000).
				expireAfterWrite(1, TimeUnit.DAYS).
				build();		   
	}
	
	public void registerMetric(Metric metric) {
		metrics.put(metric.getName(), metric);
	}
	
	public List<Metric> getMetrics() {
		return ImmutableSortedSet.copyOf(metrics.asMap().values()).asList();
	}
	
	public Metric getByName(String metricName) {
		return metrics.getIfPresent(metricName);
	}

	public boolean isKnownMetricName(String metricName) {
		return metrics.getIfPresent(metricName) != null;
	}
	
}
