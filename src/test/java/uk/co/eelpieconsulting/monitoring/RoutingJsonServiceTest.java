package uk.co.eelpieconsulting.monitoring;

import com.google.common.collect.Maps;
import org.junit.Test;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;
import uk.co.eelpieconsulting.monitoring.model.transforms.AsString;

import java.util.Map;

import static org.junit.Assert.assertEquals;

public class RoutingJsonServiceTest {

  private RoutingJsonService routingJsonService = new RoutingJsonService(new TransformsDAO());

  @Test
  public void emptyMeansEmpty() throws Exception {
    assertEquals("{}", routingJsonService.toString(Maps.newHashMap()));
  }

  @Test
  public void canDumpRoutesToJsonString() throws Exception {
    Map<String, MetricRouting> routings = Maps.newHashMap();
    routings.put("test", new MetricRouting("agauge", "temp", new AsString()));

    String json = routingJsonService.toString(routings);

    assertEquals("{\"test\":{\"gauge\":\"agauge\",\"metricName\":\"temp\",\"transform\":{\"name\":\"toString\"}}}", json);
  }

  @Test
  public void canRoundTripRoutingMapToJson() throws Exception {
    Map<String, MetricRouting> routings = Maps.newHashMap();
    routings.put("test", new MetricRouting("agauge", "temp", new AsString()));

    String json = routingJsonService.toString(routings);

    Map<String, MetricRouting> stringMetricRoutingMap = routingJsonService.fromString(json);

    System.out.println(stringMetricRoutingMap);

    assertEquals(1, stringMetricRoutingMap.size());
  }

}
