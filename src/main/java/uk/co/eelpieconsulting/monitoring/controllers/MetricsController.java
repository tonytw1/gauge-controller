package uk.co.eelpieconsulting.monitoring.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.ModelAndView;

import uk.co.eelpieconsulting.monitoring.MetricsDAO;
import uk.co.eelpieconsulting.monitoring.model.Metric;

import java.util.List;

@Controller
public class MetricsController {

	private final MetricsDAO metricsDAO;
	
	@Autowired
	public MetricsController(MetricsDAO metricsDAO) {
		this.metricsDAO = metricsDAO;
	}
	
	@RequestMapping(value="/browse", method=RequestMethod.GET)
	public ModelAndView browse() {
		return new ModelAndView("templates/metrics").addObject("availableMetrics", metricsDAO.getMetrics());
	}

	@RequestMapping(value="/metrics", method=RequestMethod.GET)
	public ModelAndView metrics() {
		List<Metric> metrics = metricsDAO.getMetrics();
		return new ModelAndView("templates/export").addObject("availableMetrics", metrics);
	}

}
