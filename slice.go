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

// First reslices to only include the first n elements.
// If n > len(slice), it returns the same slice unchanged.
func (ts TrickSlice) First(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		return TrickSlice(v)
	} else {
		return TrickSlice(v.Slice(0, n))
	}
}

// Last reslices to only include the last n elements.
// If n > len(slice), it returns the same slice unchanged.
func (ts TrickSlice) Last(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		return TrickSlice(v)
	} else {
		return TrickSlice(v.Slice(v.Len()-n, v.Len()))
	}
}

func isValidMapFunc(funcType, sliceType reflect.Type) bool {
	return funcType.NumIn() == 1 && funcType.NumOut() == 1 &&
		funcType.In(0) == sliceType
}

// GroupBy collects the slice values into a map, where the keys are the return
// value of the grouping function and the values are slices of elements that
// correspond to that key.
func (ts TrickSlice) GroupBy(fn interface{}) TrickMap {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidMapFunc(f.Type(), v.Type().Elem()) {
		panic("invalid grouping function")
	}
	keyType := f.Type().Out(0)
	valType := reflect.SliceOf(v.Type().Elem())
	mapType := reflect.MapOf(keyType, valType)

	out := reflect.MakeMap(mapType)
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		key := f.Call([]reflect.Value{val})[0]
		group := out.MapIndex(key)
		if !group.IsValid() {
			group = reflect.MakeSlice(valType, 0, 1)
		}
		out.SetMapIndex(key, reflect.Append(group, val))
	}

	return TrickMap(out)
}
