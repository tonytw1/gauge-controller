package routing

import (
	"github.com/tonytw1/gauges/model"
	"sync"
)

type RoutesTable struct {
	Routes       sync.Map
	RoutingTable sync.Map
}

func (svc RoutesTable) AddRoute(route model.Route) {
	svc.Routes.Store(route.Id, route)
	// Update routing table for effected metric
	effectedRoutes, ok := svc.RoutingTable.Load(route.FromMetric)
	if ok {
		updated := append(effectedRoutes.([]model.Route), route)
		svc.RoutingTable.Store(route.FromMetric, updated)
	} else {
		svc.RoutingTable.Store(route.FromMetric, []model.Route{route})
	}
}
