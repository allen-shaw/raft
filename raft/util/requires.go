package util

import "fmt"

func RequireTrue(expression bool, format string, args ...interface{}) {
	if !expression {
		panic(fmt.Sprintf(format, args))
	}
}
