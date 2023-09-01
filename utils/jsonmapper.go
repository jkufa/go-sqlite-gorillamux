package jsonmapper

import (
	"strconv"
)

// Function StrToInt64 Converts a string to int64 and returns it.
// It returns -1 if the conversion fails
func StrToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return i
}
