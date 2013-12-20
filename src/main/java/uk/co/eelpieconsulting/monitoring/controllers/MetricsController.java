package uk.co.eelpieconsulting.monitoring.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.ModelAndView;

import uk.co.eelpieconsulting.monitoring.MetricsDAO;

@Controller
public class MetricsController {

	private final MetricsDAO metricsDAO;
	
	@Autowired
	public MetricsController(MetricsDAO metricsDAO) {
		this.metricsDAO = metricsDAO;
	}
	
	@RequestMapping(value="/metrics", method=RequestMethod.GET)
	public ModelAndView homepage() {
		return new ModelAndView("metrics").addObject("availableMetrics", metricsDAO.getMetrics());
	}	
	
}
