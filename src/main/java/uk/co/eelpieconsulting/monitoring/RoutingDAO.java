package uk.co.eelpieconsulting.monitoring;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.List;
import java.util.Map;

import org.apache.commons.io.IOUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.base.Strings;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;

@Component
public class RoutingDAO {

	private Map<String, MetricRouting> routings;
	private ObjectMapper objectMapper;

	private String stateFile;

	@Autowired
	public RoutingDAO(@Value("${state.file}") String stateFile) throws IOException {
		this.routings = Maps.newConcurrentMap();
		objectMapper = new ObjectMapper();
		
		final String conf = IOUtils.toString(new FileInputStream(new File(stateFile)));
		if (Strings.isNullOrEmpty(conf)) {
			this.routings = Maps.newConcurrentMap();
			return;
		}		
		this.routings = objectMapper.readValue(conf, new TypeReference<Map<String, MetricRouting>>() {});
	}

	public void setRouting(String gauge, String metricName, double scale) throws IOException {
		routings.put(gauge, new MetricRouting(gauge, metricName, scale));
		
		final String conf = objectMapper.writeValueAsString(routings);
		IOUtils.write(conf, new FileOutputStream(new File(stateFile)));
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
