package uk.co.eelpieconsulting.monitoring.model.transforms;

import uk.co.eelpieconsulting.monitoring.model.Metric;

public interface Transform {

  public String getName();

  public String transform(Metric metric);

}
