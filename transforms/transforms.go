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

func Transforms() map[string]func(string) (int, error) {
	transforms := make(map[string]func(string) (int, error))
	transforms["to_int"] = toInt
	transforms["boolean_to_int"] = booleanToInt
	return transforms
}

func GetTransformByName(name string) (func(string) (int, error), bool) {
	transforms := make(map[string]func(string) (int, error))
	transforms["to_int"] = toInt
	transforms["boolean_to_int"] = booleanToInt
	transform, ok := transforms[name]
	return transform, ok
}
