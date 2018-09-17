package uk.co.eelpieconsulting.monitoring;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.collect.Maps;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;

import java.io.IOException;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

@Component
public class RoutingJsonService {

  private final ObjectMapper objectMapper;
  private TransformsDAO transformsDAO;

  @Autowired
  public RoutingJsonService(TransformsDAO transformsDAO) {
    this.transformsDAO = transformsDAO;
    this.objectMapper = new ObjectMapper();
  }

  public String toString(Map<String, MetricRouting> routings) throws JsonProcessingException {
    return objectMapper.writeValueAsString(routings);
  }

  public Map<String, MetricRouting> fromString(String json) throws IOException {
    JsonNode jsonNode = objectMapper.readTree(json);

    HashMap<String, MetricRouting> result = Maps.newHashMap();

    Iterator<String> fields = jsonNode.fieldNames();
    while(fields.hasNext()) {
      String field = fields.next();
      JsonNode next = jsonNode.get(field);
      result.put(field, new MetricRouting(next.get("gauge").asText(), next.get("metricName").asText(), transformsDAO.transformByName(next.get("transform").get("name").asText())));
    }
    return result;
  }

}
