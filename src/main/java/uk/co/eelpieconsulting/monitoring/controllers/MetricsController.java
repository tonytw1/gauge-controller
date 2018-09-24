package uk.co.eelpieconsulting.monitoring.controllers;

import com.google.common.collect.Lists;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.servlet.ModelAndView;
import uk.co.eelpieconsulting.monitoring.MetricsDAO;
import uk.co.eelpieconsulting.monitoring.RoutingDAO;
import uk.co.eelpieconsulting.monitoring.TransformsDAO;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;
import uk.co.eelpieconsulting.monitoring.model.transforms.Transform;

import java.util.List;
import java.util.regex.Pattern;

@Controller
public class MetricsController {

  private static final Pattern NO_ODD_CHARACTERS = Pattern.compile("[^a-zA-Z0-9]+");

  private final MetricsDAO metricsDAO;
  private final RoutingDAO routingDAO;
  private final TransformsDAO transformsDAO;

  @Autowired
  public MetricsController(MetricsDAO metricsDAO, RoutingDAO routingDAO, TransformsDAO transformsDAO) {
    this.metricsDAO = metricsDAO;
    this.routingDAO = routingDAO;
    this.transformsDAO = transformsDAO;
  }

  @RequestMapping(value = "/browse", method = RequestMethod.GET)
  public ModelAndView browse() {
    return new ModelAndView("templates/browse").addObject("availableMetrics", metricsDAO.getMetrics());
  }

  @RequestMapping(value = "/metrics", method = RequestMethod.GET)
  public ModelAndView metrics() {
    List<Metric> metrics = metricsDAO.getMetrics();

    List<Metric> cleaned = Lists.newArrayList();
    for (Metric m : metrics) {
      try {
        Double.parseDouble(m.getLastValue());
        cleaned.add(new Metric(makeSafeName(m), m.getType(), m.getLastValue(), m.getDate()));
      } catch (Exception e) {
        // TODO
      }
    }

    return new ModelAndView("templates/export").addObject("availableMetrics", cleaned);
  }

  @RequestMapping(value = "/transformed", method = RequestMethod.GET)
  public ModelAndView routed() {

    List<Metric> cleaned = Lists.newArrayList();

    for (MetricRouting route : routingDAO.getGaugeRoutes().values()) {
      if (route.getTransform() != null) {
        Metric metric = metricsDAO.getByName(route.getMetricName());
        if (metric != null) {
          Transform transform = transformsDAO.transformByName(route.getTransform().getName());
          String value = transform.transform(metric);
          Double.parseDouble(value);
          cleaned.add(new Metric(makeSafeName(metric), metric.getType(), value, metric.getDate()));
        }
      }
    }

    return new ModelAndView("templates/export").addObject("availableMetrics", cleaned);
  }

  private String makeSafeName(Metric m) {
    return NO_ODD_CHARACTERS.matcher(m.getName()).replaceAll("");
  }

}
