package views

import (
	"encoding/json"
	"github.com/tonytw1/gauges/model"
	"sort"
	"strings"
	"sync"
)

func RoutesAsJson(routes sync.Map) []byte {
	var routesList = make([]model.Route, 0)
	routes.Range(func(k, v interface{}) bool {
		routesList = append(routesList, v.(model.Route))
		return true
	})
	sort.Slice(routesList, func(i, j int) bool {
		return strings.Compare(routesList[i].FromMetric, routesList[j].FromMetric) < 0
	})
	asJson, _ := json.Marshal(routesList)
	return asJson
}
