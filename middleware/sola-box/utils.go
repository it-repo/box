package box

import (
	"strconv"

	"github.com/ddosakura/sola/v2"
)

func toStringArray(x interface{}) []string {
	y := x.([]interface{})
	result := make([]string, 0, len(y))
	for _, o := range y {
		result = append(result, o.(string))
	}
	return result
}

func getPageSize(c sola.Context) (page, size int) {
	r := c.Request()
	var err error
	if page, err = strconv.Atoi(r.URL.Query().Get("page")); err != nil {
		page = 1
	}
	if size, err = strconv.Atoi(r.URL.Query().Get("size")); err != nil {
		size = 5
	}
	return
}
