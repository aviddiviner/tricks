package tricks

import (
	"reflect"
)

type TrickSlice reflect.Value

var trickSliceType = reflect.TypeOf((*TrickSlice)(nil)).Elem()

func Slice(sliceOrElement interface{}, moreElements ...interface{}) TrickSlice {
	v := reflect.ValueOf(sliceOrElement)
	if !v.IsValid() { // nil
		v = reflect.ValueOf([]interface{}{})
		return TrickSlice(v)
	}
	if len(moreElements) == 0 {
		if v.Kind() == reflect.Slice { // TrickSlice is a Kind of reflect.Struct
			return TrickSlice(v)
		}
		if v.Type() == trickSliceType {
			return sliceOrElement.(TrickSlice)
		}
	}
	typ := reflect.SliceOf(v.Type())
	slice := reflect.MakeSlice(typ, len(moreElements)+1, len(moreElements)+1)
	slice.Index(0).Set(v)
	for i := 0; i < len(moreElements); i++ {
		el := reflect.ValueOf(moreElements[i])
		// Nicer messages on our panics.
		if !el.IsValid() {
			// panic: reflect: call of reflect.Value.Set on zero Value
			panic("tricks: Slice: unexpected nil value as variadic element")
		}
		if el.Type() != v.Type() {
			// panic: reflect.Set: value of type string is not assignable to type int
			panic("tricks: Slice: value of type " + el.Type().String() +
				" is not assignable to slice of type " + typ.String())
		}
		slice.Index(i + 1).Set(el)
	}
	return TrickSlice(slice)
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

// Len returns the length of the slice.
func (ts TrickSlice) Len() int {
	return reflect.Value(ts).Len()
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
		funcType.In(0) == sliceType.Elem()
}

func isValidBoolFunc(funcType, sliceType reflect.Type) bool {
	return isValidMapFunc(funcType, sliceType) &&
		funcType.Out(0).Kind() == reflect.Bool
}

// Any returns true if the given function returns true for any element in the
// slice. Otherwise, it returns false.
func (ts TrickSlice) Any(fn interface{}) bool {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.Any: invalid function type")
	}

	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if f.Call([]reflect.Value{val})[0].Bool() {
			return true
		}
	}
	return false
}

// All returns true if the given function returns true for every element in the
// slice. Otherwise, it returns false.
func (ts TrickSlice) All(fn interface{}) bool {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.All: invalid function type")
	}

	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if !f.Call([]reflect.Value{val})[0].Bool() {
			return false
		}
	}
	return true
}

// None returns true if the given function returns false for every element in the
// slice. Otherwise, it returns false.
func (ts TrickSlice) None(fn interface{}) bool {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.None: invalid function type")
	}

	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if f.Call([]reflect.Value{val})[0].Bool() {
			return false
		}
	}
	return true
}

// One returns true if the given function returns true for exactly one element
// in the slice. Otherwise, it returns false.
func (ts TrickSlice) One(fn interface{}) bool {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.One: invalid function type")
	}

	found := false
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if f.Call([]reflect.Value{val})[0].Bool() {
			if found {
				return false
			}
			found = true
		}
	}
	return found
}

// Many returns true if the given function returns true for more than one element
// in the slice. Otherwise, it returns false.
func (ts TrickSlice) Many(fn interface{}) bool {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.Many: invalid function type")
	}

	found := false
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if f.Call([]reflect.Value{val})[0].Bool() {
			if found {
				return true
			}
			found = true
		}
	}
	return false
}

// GroupBy collects the slice values into a map, where the keys are the return
// value of the grouping function and the values are slices of elements that
// correspond to that key.
func (ts TrickSlice) GroupBy(fn interface{}) TrickMap {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !isValidMapFunc(f.Type(), v.Type()) {
		panic("tricks: slice.GroupBy: invalid function type")
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
	if !isValidMapFunc(f.Type(), v.Type()) {
		panic("tricks: slice.Map: invalid function type")
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
