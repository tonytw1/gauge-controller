package uk.co.eelpieconsulting.monitoring;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import io.minio.MinioClient;
import io.minio.errors.InvalidEndpointException;
import io.minio.errors.InvalidPortException;
import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.Metric;
import uk.co.eelpieconsulting.monitoring.model.MetricRouting;

import java.io.InputStream;
import java.io.StringBufferInputStream;
import java.util.List;
import java.util.Map;

@Component
public class RoutingDAO {

  private final static Logger log = Logger.getLogger(RoutingDAO.class);

  private final MinioClient minioClient;
  private final Map<String, MetricRouting> routings;
  private final ObjectMapper objectMapper;

  private final String bucketName;
  private final String filename = "routings.json";

  @Autowired
  public RoutingDAO(
          @Value("${state.s3.endpoint}") String endPoint,
          @Value("${state.s3.accesskey}") String accessKey,
          @Value("${state.s3.secretkey}") String secretKey,
          @Value("${state.s3.bucket}") String bucketName
  ) throws InvalidPortException, InvalidEndpointException {
    this.objectMapper = new ObjectMapper();
    this.minioClient = new MinioClient(endPoint, accessKey, secretKey);
    this.routings = loadRoutings();
    this.bucketName = bucketName;
  }

  public boolean isRoutedMetric(Metric metric) {
    return !getRoutingsForMetric(metric.getName()).isEmpty();
  }

  public List<MetricRouting> getRoutingsForMetric(String metricName) {
    List<MetricRouting> metricRoutings = Lists.newArrayList();
    for (MetricRouting metricRouting : routings.values()) {
      if (metricRouting.getMetricName().equals(metricName)) {
        metricRoutings.add(metricRouting);
      }
    }
    return metricRoutings;
  }

  public Map<String, MetricRouting> getGaugeRoutes() {
    return routings;
  }

  public void setRouting(String gauge, String metricName, double scale) {
    routings.put(gauge, new MetricRouting(gauge, metricName, scale));
    persistRoutings(routings);
  }

  public void clearRouting(String gauge) {
    routings.remove(gauge);
    persistRoutings(routings);
  }

  private void persistRoutings(Map<String, MetricRouting> routings) {
    try {
      final String asJson = objectMapper.writeValueAsString(routings);
      minioClient.putObject(bucketName, filename, new StringBufferInputStream(asJson), "appplication/json");
    } catch (Exception e) {
      log.error("Failed to persist routes", e);
      throw new RuntimeException(e);
    }
  }

  private Map<String, MetricRouting> loadRoutings() {
    try {
      InputStream object = minioClient.getObject(bucketName, filename);
      String json = IOUtils.toString(object);
      return objectMapper.readValue(json, new TypeReference<Map<String, MetricRouting>>() {
      });
    } catch (Exception e) {
      log.error("Failed to load routes; returning empty", e);
      return Maps.newConcurrentMap();
    }
  }

}