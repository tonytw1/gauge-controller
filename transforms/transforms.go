package transforms

import "strconv"

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

var transforms = map[string]func(string) (int, error){
	"to_int":         toInt,
	"boolean_to_int": booleanToInt,
}

func Transforms() map[string]func(string) (int, error) {
	return transforms
}

func GetTransformByName(name string) (func(string) (int, error), bool) {
	transform, ok := transforms[name]
	return transform, ok
}
