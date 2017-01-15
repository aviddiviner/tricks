package tricks

import (
	"reflect"
	"strings"
)

// Join joins a slice of strings into a single string, separated by glue.
func (ts TrickSlice) Join(glue string) string {
	v := reflect.Value(ts)
	switch v.Type() {
	case typeStringSlice:
		return strings.Join(v.Interface().([]string), glue)
	default:
		panic("tricks: slice.Join: not a slice of strings")
	}
}
