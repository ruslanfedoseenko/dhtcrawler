package Utils

import (
	"fmt"
	"strings"
)

func SliceToPathString(values []interface{}) string {
	s := make([]string, len(values)) // Pre-allocate the right size
	for index := range values {
		s[index] = fmt.Sprintf("%v", values[index])
	}
	return strings.Join(s, "/")
}

func IndexOf(slice []string, item string) int {

	for index, s := range slice {
		if strings.Trim(s, "\n\r ") == strings.Trim(item, "\n\r ") {
			return index
		}
	}

	return -1
}
