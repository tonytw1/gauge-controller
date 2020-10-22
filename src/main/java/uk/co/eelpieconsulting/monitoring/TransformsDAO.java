package uk.co.eelpieconsulting.monitoring;

import com.google.common.base.Predicate;
import com.google.common.collect.Iterables;
import com.google.common.collect.Lists;
import org.springframework.stereotype.Component;
import uk.co.eelpieconsulting.monitoring.model.transforms.*;

import java.util.List;

@Component
public class TransformsDAO {

  public static final Transform AS_STRING = new AsString();

  private final List<Transform> transforms = Lists.newArrayList(
          AS_STRING,
          new Scaled(-1),
          new Scaled(0.001),
          new Scaled(0.01),
          new Scaled(0.1),
          new Scaled(10.0),
          new Scaled(100.0),
          new Scaled(1000.0),
          new Scaled(10000.0),
          new FromBoolean(),
          new FromInvertedBoolean(),
          new FromHex(),
          new LastChanged()
  );

  public Transform transformByName(String transformName) {
    Transform transform = Iterables.find(transforms, new Predicate<Transform>() {
      public boolean apply(Transform transform) {
        return transform.getName().equals(transformName);
      }
    });
    if (transform == null) {
      throw new RuntimeException("Could not find transform: " + transformName); // TODO
    }

    return transform;
  }

  public  List<Transform> getAll() {
    return transforms;
  }

}
