package routing

import (
	"github.com/tonytw1/gauges/model"
	"sync"
)

type MetricsTable struct {
	metrics sync.Map
}

func (svc *MetricsTable) AddMetrics(metric model.Metric) {
	svc.metrics.Store(metric.Name, metric)
}

func (svc *MetricsTable) AllMetrics() []model.Metric {
	var ms = make([]model.Metric, 0)
	svc.metrics.Range(func(k, v interface{}) bool {
		ms = append(ms, v.(model.Metric))
		return true
	})
	return ms
}

func (svc *MetricsTable) GetMetric(metricName string) (*model.Metric, bool) {
	value, ok := svc.metrics.Load(metricName)
	if ok {
		metric := value.(model.Metric)
		return &metric, ok
	}
	return nil, ok
}

func NewMetricsTable() MetricsTable {
	return MetricsTable{metrics: sync.Map{}}
}
