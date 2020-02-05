package util

import "fmt"

// WrapError wraps an error to return as a byteslice
func WrapError(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
}
