package common

func Map[Input any, Output any](arr []Input, fn func(val Input) Output) []Output {
	res := []Output{}
	for _, val := range arr {
		res = append(res, fn(val))
	}
	return res
}

func EMap[Input any, Output any](arr []Input, fn func(val Input) (Output, error)) ([]Output, error) {
	res := []Output{}
	for _, val := range arr {
		r, err := fn(val)
		if err != nil {
			return nil, err
		}
		res = append(res, r)

	}
	return res, nil
}
