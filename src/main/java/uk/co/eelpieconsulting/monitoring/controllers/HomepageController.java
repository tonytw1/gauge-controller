package uk.co.eelpieconsulting.monitoring.controllers;

import com.google.common.base.Function;
import com.google.common.base.Strings;
import com.google.common.collect.*;
import org.checkerframework.checker.nullness.qual.Nullable;
import org.joda.time.DateTime;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import uk.co.eelpieconsulting.monitoring.GaugeDAO;
import uk.co.eelpieconsulting.monitoring.MetricPublisher;
import uk.co.eelpieconsulting.monitoring.MetricsDAO;
import uk.co.eelpieconsulting.monitoring.RoutingDAO;
import uk.co.eelpieconsulting.monitoring.model.Metric;

import java.util.List;

@Controller
public class HomepageController {

	private final GaugeDAO gaugeDAO;
	private final MetricsDAO metricsDAO;
	private final RoutingDAO routingDAO;
	private final MetricPublisher metricPublisher;

	private final Iterable<Double> powersOfTen = Iterables.transform(ContiguousSet.create(Range.closed(-2, 2), DiscreteDomain.integers()).asList(),
									new Function<Integer, Double>() {
										@Nullable
										@Override
										public Double apply(@Nullable Integer i) {
											return Math.pow(10, i);
										}
									});

	private final Iterable<Double>  negativePowersOfTen = Iterables.transform(powersOfTen, new Function<Double, Double>() {
		@Nullable
		@Override
		public Double apply(@Nullable Double aDouble) {
			return -aDouble;
		}
	});

	@Autowired
	public HomepageController(GaugeDAO gaugeDAO, MetricsDAO metricsDAO, RoutingDAO routingDAO, MetricPublisher metricPublisher) {
		this.gaugeDAO = gaugeDAO;
		this.metricsDAO = metricsDAO;
		this.routingDAO = routingDAO;
		this.metricPublisher = metricPublisher;
	}
	
	@RequestMapping(value="/", method=RequestMethod.GET)
	public ModelAndView homepage() {
		List<Double> scales = Lists.newArrayList(Iterables.concat(negativePowersOfTen, powersOfTen));
		scales.sort(Ordering.natural());

		ModelAndView mv = new ModelAndView("templates/homepage").
			addObject("gauges", gaugeDAO.getGauges()).
			addObject("availableMetrics", metricsDAO.getMetrics()).
			addObject("scales", scales).
			addObject("gaugeRoutes", routingDAO.getGaugeRoutes());
		return mv;
	}
	
	@RequestMapping(value="/", method=RequestMethod.POST)
	public ModelAndView set(@RequestParam String gauge,
			@RequestParam String metric,
			@RequestParam double scale) throws Exception {
		if (!gaugeDAO.isKnownGaugeName(gauge)) {
			throw new RuntimeException("Unknown gauge: " + gauge);
		}

		if (Strings.isNullOrEmpty(metric)) {
			routingDAO.clearRouting(gauge);
			final Metric nullMetric = new Metric(null, null, "0", DateTime.now());
			metricPublisher.publishOntoGaugesChannel(gauge, nullMetric, 1);
			return new ModelAndView(new RedirectView("/"));		
		}
		
		if (!metricsDAO.isKnownMetricName(metric)) {
			throw new RuntimeException("Unknown metric: " + metric);
		}
		
		routingDAO.setRouting(gauge, metric, scale);
		metricPublisher.publishOntoGaugesChannel(gauge, metricsDAO.getByName(metric), scale);	// TODO race condition
		return new ModelAndView(new RedirectView("/"));		
	}

}
