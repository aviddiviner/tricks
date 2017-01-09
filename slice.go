package tricks

import (
	"reflect"
)

type TrickSlice reflect.Value

func Slice(anySlice interface{}) TrickSlice {
	v := reflect.ValueOf(anySlice)
	if v.Kind() != reflect.Slice {
		panic("input is not a slice")
	}
	return TrickSlice(v)
}

func (ts TrickSlice) Value() interface{} {
	return reflect.Value(ts).Interface()
}

// First reslices to only include the first N elements.
// If N > len(slice), it returns the same slice unchanged.
func (ts TrickSlice) First(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		return TrickSlice(v)
	} else {
		return TrickSlice(v.Slice(0, n))
	}
}

// Last reslices to only include the last N elements.
// If N > len(slice), it returns the same slice unchanged.
func (ts TrickSlice) Last(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		return TrickSlice(v)
	} else {
		return TrickSlice(v.Slice(v.Len()-n, v.Len()))
	}
}
