# Gauge controller

Simple router and UI to cherry pick metrics from the metrics bus and send them to display devices; namely the [analog monitoring dashboard](https://github.com/tonytw1/analog-monitoring-system).

Metrics are made available on a metrics MQTT channel. The display devices listen on a separate gauges MQTT channel and periodically announce themselves.

The controller maintains an in memory list of recently announced gauges and a list of available metrics seen in the metrics channel.
The UI can be used to route a given metric to a given display device with an optional scaling.

Routed metrics are scaled and echoed into the gauges channel.

The routing table is persisted to an S3 bucket.

## Build

Spring Boot / Maven project.

```
mvn clean package
```

Which will give a runnable fat jar file

```
java -jar target/gauge-controller-0.0.1-SNAPSHOT.jar --spring.config.location=application.properties
```
