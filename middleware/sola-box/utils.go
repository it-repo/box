package box

func toStringArray(x interface{}) []string {
	y := x.([]interface{})
	result := make([]string, 0, len(y))
	for _, o := range y {
		result = append(result, o.(string))
	}
	return result
}
