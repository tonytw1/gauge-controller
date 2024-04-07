package transforms

import "testing"

func TestGetTransformByName(t *testing.T) {
	transform, ok := GetTransformByName("to_int")
	println(ok)
	println(transform("123"))
}
