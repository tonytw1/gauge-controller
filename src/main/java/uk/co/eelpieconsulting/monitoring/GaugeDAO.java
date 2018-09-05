package uk.co.eelpieconsulting.monitoring;

import com.google.common.collect.Maps;
import com.google.common.collect.Ordering;
import org.apache.log4j.Logger;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Gauge;

import javax.annotation.Nullable;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Component
public class GaugeDAO {

	private final static Logger log = Logger.getLogger(GaugeDAO.class);
	
	private Map<String, Gauge> gauges;

	private Ordering<Gauge> byGaugeName = new Ordering<Gauge>() {
		@Override
		public int compare(@Nullable Gauge gauge, @Nullable Gauge gauge2) {
			return gauge.getName().compareTo(gauge2.getName());
		}
	};

	public GaugeDAO() {
		this.gauges = Maps.newConcurrentMap();
	}
	
	public void registerGauge(Gauge gauge) {
		log.info("Registering gauge: " + gauge);
		gauges.put(gauge.getName(), gauge);
	}
	
	public List<Gauge> getGauges() {
		List<Gauge> gauges = new ArrayList<>(this.gauges.values());
		gauges.sort(byGaugeName);
		return gauges;
	}

	public boolean isKnownGaugeName(String gaugeName) {
		return gauges.keySet().contains(gaugeName);
	}
	
}
