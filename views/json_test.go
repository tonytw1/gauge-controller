package views

import (
	"strings"
	"sync"
	"testing"
)

func TestRoutesAsJson_CanHandleEmptyRoutes(t *testing.T) {
	var routes = sync.Map{}

	asJson := RoutesAsJson(routes)

	if !(strings.Compare(string(asJson), "[]") == 0) {
		t.Fatalf("Expected empty json array")
	}
}
