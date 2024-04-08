package transforms

import "strconv"
import "github.com/tonytw1/gauges/model"

var toInt = func(input string) (int, error) {
	return strconv.Atoi(input)
}

var booleanToInt = func(input string) (int, error) {
	i, err := strconv.ParseBool(input)
	if err != nil {
		return 0, err
	}
	if i {
		return 1, nil
	} else {
		return 0, nil
	}
}

var invertedBooleanToInt = func(input string) (int, error) {
	i, err := booleanToInt(input)
	if err != nil {
		return 0, err
	}
	if i == 1 {
		return 0, nil
	} else {
		return 1, nil
	}
}

var transforms = map[string]model.Transform{
	"to_int":                  {Name: "to_int", Transform: toInt},
	"boolean_to_int":          {Name: "boolean_to_int", Transform: booleanToInt},
	"inverted_boolean_to_int": {Name: "inverted_boolean_to_int", Transform: invertedBooleanToInt},
}

func Transforms() map[string]model.Transform {
	return transforms
}

func GetTransformByName(name string) (model.Transform, bool) {
	transform, ok := transforms[name]
	return transform, ok
}
