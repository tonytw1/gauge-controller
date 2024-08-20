package routing

import (
	"github.com/tonytw1/gauges/model"
	"sync"
)

type RoutesTable struct {
	routes       sync.Map
	routingTable sync.Map
}

func NewRoutesTable() RoutesTable {
	return RoutesTable{routes: sync.Map{}, routingTable: sync.Map{}}
}

func (svc RoutesTable) AddRoute(route model.Route) {
	svc.routes.Store(route.Id, route)
	// Update routing table for effected metric
	effectedRoutes, ok := svc.routingTable.Load(route.FromMetric)
	if ok {
		updated := append(effectedRoutes.([]model.Route), route)
		svc.routingTable.Store(route.FromMetric, updated)
	} else {
		svc.routingTable.Store(route.FromMetric, []model.Route{route})
	}
}

func (svc *RoutesTable) AllRoutes() []model.Route {
	var routesList []model.Route = make([]model.Route, 0)
	svc.routes.Range(func(k, v interface{}) bool {
		routesList = append(routesList, v.(model.Route))
		return true
	})
	return routesList
}

func (svc *RoutesTable) GetRoute(id string) (*model.Route, bool) {
	value, ok := svc.routes.Load(id)
	if ok {
		route := value.(model.Route)
		return &route, ok
	}
	return nil, ok
}

func (svc *RoutesTable) Delete(route *model.Route) {
	svc.routes.Delete(route.Id)
	// Update routing table for effected metric
	effectedMetric := route.FromMetric
	effectedMetricRoutes, ok := svc.routingTable.Load(effectedMetric)
	if ok {
		// Filter out the route that was deleted
		filtered := make([]model.Route, 0)
		for _, route := range effectedMetricRoutes.([]model.Route) {
			if route.Id != route.Id {
				filtered = append(filtered, route)
			}
		}
		svc.routingTable.Store(effectedMetric, filtered)
	}
}

func (svc RoutesTable) GetRoutesForMetric(metricName string) (any, bool) {
	return svc.routingTable.Load(metricName) // TODO pull any down
}
