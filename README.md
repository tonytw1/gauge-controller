# Gauge controller

Simple router and UI to cherry pick metrics from the metrics bus and send them to display devices; namely the [analog monitoring dashboard](https://github.com/tonytw1/analog-monitoring-system).

Metrics are made available on a metrics MQTT channel. The display devices listen on a seperate gauges MQTT channel and periodically announce themselves.

The application maintains an in memory list of gauges which have been announced on the gauges channel and available metrics which have been seen on the metrics channel.
The UI can then be used to route a given metrics to a given display device with an optional scaling.
Routed metrics are echoed into the gauges channel.

The routing table is persisted to an S3 bucket.


