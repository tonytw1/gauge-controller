package uk.co.eelpieconsulting.monitoring.controllers;

import com.google.common.collect.Lists;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.ModelAndView;
import uk.co.eelpieconsulting.monitoring.MetricsDAO;
import uk.co.eelpieconsulting.monitoring.model.Metric;

import java.util.List;
import java.util.regex.Pattern;

@Controller
public class MetricsController {

  private static final Pattern NO_ODD_CHARACTERS = Pattern.compile("[^a-zA-Z0-9]+");

  private final MetricsDAO metricsDAO;

  @Autowired
	public MetricsController(MetricsDAO metricsDAO) {
		this.metricsDAO = metricsDAO;
	}
	
	@RequestMapping(value="/browse", method=RequestMethod.GET)
	public ModelAndView browse() {
		return new ModelAndView("templates/browse").addObject("availableMetrics", metricsDAO.getMetrics());
	}

	@RequestMapping(value="/metrics", method=RequestMethod.GET)
	public ModelAndView metrics() {
		List<Metric> metrics = metricsDAO.getMetrics();

    List<Metric> cleaned = Lists.newArrayList();
    for (Metric m: metrics) {
		  String safeName = NO_ODD_CHARACTERS.matcher(m.getName()).replaceAll("");

		  try {
		    Double.parseDouble(m.getLastValue());
        cleaned.add(new Metric(safeName, m.getType(), m.getLastValue(), m.getDate()));
      } catch (Exception e) {
      }
    }

		return new ModelAndView("templates/export").addObject("availableMetrics", cleaned);
	}

}
