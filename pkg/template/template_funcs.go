package template

func Sum(args ...float64) float64 {
	if len(args) == 0 {
		return 0
	}

	var res float64
	for _, arg := range args {
		res += arg
	}

	return res
}

func Gt(arg1, arg2 float64) bool {
	return arg1 > arg2
}

func Gte(arg1, arg2 float64) bool {
	return arg1 >= arg2
}

func Lt(arg1, arg2 float64) bool {
	return arg1 < arg2
}

func Lte(arg1, arg2 float64) bool {
	return arg1 <= arg2
}