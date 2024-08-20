package views

import (
	"github.com/tonytw1/gauges/model"
	"strings"
	"testing"
)

func TestRoutesAsJson_CanHandleEmptyRoutes(t *testing.T) {
	var routes = []model.Route{}

	asJson := RoutesAsJson(routes)

	if !(strings.Compare(string(asJson), "[]") == 0) {
		t.Fatalf("Expected empty json array")
	}
}
