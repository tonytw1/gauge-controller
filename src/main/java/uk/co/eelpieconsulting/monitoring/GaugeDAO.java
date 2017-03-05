package uk.co.eelpieconsulting.monitoring;

import java.util.List;
import java.util.Map;

import org.apache.log4j.Logger;
import org.springframework.stereotype.Component;

import uk.co.eelpieconsulting.monitoring.model.Gauge;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;

@Component
public class GaugeDAO {

	private final static Logger log = Logger.getLogger(GaugeDAO.class);
	
	private Map<String, Gauge> gauges;

	public GaugeDAO() {
		this.gauges = Maps.newConcurrentMap();
	}
	
	public void registerGauge(Gauge gauge) {
		log.info("Registering gauge: " + gauge);
		gauges.put(gauge.getName(), gauge);
	}
	
	public List<Gauge> getGauges() {
		return Lists.newArrayList(gauges.values());
	}

	public boolean isKnownGaugeName(String gaugeName) {
		return gauges.keySet().contains(gaugeName);
	}
	
}
