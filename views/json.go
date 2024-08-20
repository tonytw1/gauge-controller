package views

import (
	"encoding/json"
	"github.com/tonytw1/gauges/model"
	"sort"
	"strings"
)

func RoutesAsJson(routes []model.Route) []byte {
	sort.Slice(routes, func(i, j int) bool {
		return strings.Compare(routes[i].FromMetric, routes[j].FromMetric) < 0
	})
	asJson, _ := json.Marshal(routes)
	return asJson
}
