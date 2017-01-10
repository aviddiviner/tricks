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

// Copy returns a new slice containing the same values.
func (ts TrickSlice) Copy() TrickSlice {
	in := reflect.Value(ts)
	out := reflect.MakeSlice(in.Type(), in.Len(), in.Len()) // TODO: v.Cap()?
	reflect.Copy(out, in)
	return TrickSlice(out)
}

// First reslices to only include the first n elements. If n > len(slice), it
// reslices to include all elements. In both cases, cap() of the new slice is
// set to equal its length.
func (ts TrickSlice) First(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		n = v.Len()
	}
	return TrickSlice(v.Slice3(0, n, n))
}

// Last reslices to only include the last n elements. If n > len(slice), it
// reslices to include all elements. In both cases, cap() of the new slice is
// set to equal its length.
func (ts TrickSlice) Last(n int) TrickSlice {
	v := reflect.Value(ts)
	if n > v.Len() {
		n = v.Len()
	}
	return TrickSlice(v.Slice3(v.Len()-n, v.Len(), v.Len()))
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
		panic("invalid group-by function")
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

// Map applies the given function to each element of the slice and stores the
// result to a new slice. The cap() of the new slice is set to equal its length.
func (ts TrickSlice) Map(fn interface{}) TrickSlice {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidMapFunc(f.Type(), v.Type().Elem()) {
		panic("invalid map function")
	}
	outType := f.Type().Out(0)
	typ := reflect.SliceOf(outType)

	out := reflect.MakeSlice(typ, v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		result := f.Call([]reflect.Value{val})[0]
		out.Index(i).Set(result)
	}

	return TrickSlice(out)
}
