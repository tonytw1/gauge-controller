package transforms

import "strconv"

func Transforms() map[string]func(string) (int, error) {
	toInt := func(input string) (int, error) {
		i, err := strconv.Atoi(input)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	booleanToInt := func(input string) (int, error) {
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

	transforms := make(map[string]func(string) (int, error))
	transforms["to_int"] = toInt
	transforms["boolean_to_int"] = booleanToInt

	return transforms
}
