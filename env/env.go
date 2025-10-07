package env

import (
	"os"
	"strconv"
)

func MustGetStr(env string) string {
	val := os.Getenv(env)
	if val == "" {
		panic(env + " environment variable not set")
	}

	return val
}

func MustGetInt(env string) int {
	val := MustGetStr(env)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(env + " environment variable is not a valid integer")
	}

	return i
}
