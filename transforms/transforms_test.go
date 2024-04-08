package transforms

import "testing"

func TestGetTransformByName_CanFindTransformByNameAndCallIt(t *testing.T) {
	transform, ok := GetTransformByName("to_int")
	if !ok {
		t.Fatalf("Expected to find transform by name")
	}
	transformed, _ := transform("123")
	if !(transformed == 123) {
		t.Fatalf("Expected to transform '123' to 123")
	}
}
