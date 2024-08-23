package routing

import (
	"github.com/tonytw1/gauges/model"
	"sync"
)

type GaugesTable struct {
	gauges sync.Map
}

func NewGaugesTable() GaugesTable {
	return GaugesTable{gauges: sync.Map{}}
}

func (svc *GaugesTable) AddGauge(gauge model.Gauge) {
	svc.gauges.Store(gauge.Name, gauge)
}

func (svc *GaugesTable) AllGauges() []model.Gauge {
	var gs = make([]model.Gauge, 0)
	svc.gauges.Range(func(k, v interface{}) bool {
		gs = append(gs, v.(model.Gauge))
		return true
	})
	return gs
}

func (svc *GaugesTable) GetGauge(gaugeName string) (*model.Gauge, bool) {
	value, ok := svc.gauges.Load(gaugeName)
	if ok {
		gauge := value.(model.Gauge)
		return &gauge, ok
	}
	return nil, ok
}
