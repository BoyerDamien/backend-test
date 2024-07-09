package common

import "strings"

func GenString(n int) string {
	name := []string{}

	for i := 0; i < n; i++ {
		name = append(name, "a")
	}
	return strings.Join(name, "")
}
