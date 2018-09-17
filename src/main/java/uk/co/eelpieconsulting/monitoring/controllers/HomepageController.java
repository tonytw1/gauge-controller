package uk.co.eelpieconsulting.monitoring.controllers;

import com.google.common.base.Strings;
import org.joda.time.DateTime;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.servlet.ModelAndView;
import org.springframework.web.servlet.view.RedirectView;
import uk.co.eelpieconsulting.monitoring.*;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.transforms.Transform;

@Controller
public class HomepageController {

  private final GaugeDAO gaugeDAO;
  private final MetricsDAO metricsDAO;
  private final RoutingDAO routingDAO;
  private final MetricPublisher metricPublisher;
  private final TransformsDAO transformsDAO;

  @Autowired
  public HomepageController(GaugeDAO gaugeDAO, MetricsDAO metricsDAO, RoutingDAO routingDAO, MetricPublisher metricPublisher, TransformsDAO transformsDAO) {
    this.gaugeDAO = gaugeDAO;
    this.metricsDAO = metricsDAO;
    this.routingDAO = routingDAO;
    this.metricPublisher = metricPublisher;
    this.transformsDAO = transformsDAO;
  }

  @RequestMapping(value = "/", method = RequestMethod.GET)
  public ModelAndView homepage() {
    ModelAndView mv = new ModelAndView("templates/homepage").
            addObject("gauges", gaugeDAO.getGauges()).
            addObject("availableMetrics", metricsDAO.getMetrics()).
            addObject("transforms", transformsDAO.getAll()).
            addObject("gaugeRoutes", routingDAO.getGaugeRoutes());
    return mv;
  }

  @RequestMapping(value = "/", method = RequestMethod.POST)
  public ModelAndView set(@RequestParam String gauge,
                          @RequestParam String metric,
                          @RequestParam String transform) throws Exception {
    if (!gaugeDAO.isKnownGaugeName(gauge)) {
      throw new RuntimeException("Unknown gauge: " + gauge);
    }

    if (Strings.isNullOrEmpty(metric)) {
      routingDAO.clearRouting(gauge);
      final Metric nullMetric = new Metric(null, null, "0", DateTime.now());
      metricPublisher.publishOntoGaugesChannel(gauge, nullMetric, TransformsDAO.AS_STRING);
      return new ModelAndView(new RedirectView("/"));
    }

    if (!metricsDAO.isKnownMetricName(metric)) {
      throw new RuntimeException("Unknown metric: " + metric);
    }

    Transform t = transformsDAO.transformByName(transform);

    routingDAO.setRouting(gauge, metric, t);
    metricPublisher.publishOntoGaugesChannel(gauge, metricsDAO.getByName(metric), t);
    return new ModelAndView(new RedirectView("/"));
  }

}
