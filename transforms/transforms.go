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

var toIntTimesZeroDotTen = func(input string) (int, error) {
	i, err := toInt(input)
	if err != nil {
		return 0, err
	}
	return i / 10, nil
}

var transforms = map[string]model.Transform{
	"to_int":                    {Name: "to_int", Transform: toInt},
	"to_int_times_zero_dot_ten": {Name: "to_int * 0.1", Transform: toIntTimesZeroDotTen},
	"boolean_to_int":            {Name: "boolean_to_int", Transform: booleanToInt},
	"inverted_boolean_to_int":   {Name: "inverted_boolean_to_int", Transform: invertedBooleanToInt},
}

func Transforms() map[string]model.Transform {
	return transforms
}

func GetTransformByName(name string) (model.Transform, bool) {
	transform, ok := transforms[name]
	return transform, ok
}
