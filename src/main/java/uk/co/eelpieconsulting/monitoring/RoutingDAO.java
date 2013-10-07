package uk.co.eelpieconsulting.monitoring;

import java.util.List;
import java.util.Map;

import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;

@Component
public class RoutingDAO {

	private Map<String, MetricRouting> routings;
	
	public RoutingDAO() {
		this.routings = Maps.newConcurrentMap();
	}

	public void setRouting(String gauge, String metricName, double scale) {
		routings.put(gauge, new MetricRouting(gauge, metricName, scale));
	}

	public boolean isRoutedMetric(Metric metric) {
		return !getRoutingsForMetric(metric.getName()).isEmpty();
	}

	public List<MetricRouting> getRoutingsForMetric(String metricName) {
		List<MetricRouting> metricRoutings = Lists.newArrayList();
		for (MetricRouting metricRouting : routings.values()) {
			if (metricRouting.getMetricName().equals(metricName)) {
				metricRoutings.add(metricRouting);
			}
		}
		return metricRoutings;
	}

	public Map<String, MetricRouting> getGaugeRoutes() {
		return routings;
	}

	public void clearRouting(String gauge) {
		routings.remove(gauge);
	}

}
