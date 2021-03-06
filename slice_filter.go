package tricks

import "reflect"

func isValidMapFunc(funcType, sliceType reflect.Type) bool {
	return funcType.NumIn() == 1 && funcType.NumOut() == 1 &&
		funcType.In(0) == sliceType.Elem()
}

func isValidReduceFunc(funcType, sliceType reflect.Type) bool {
	return funcType.NumIn() == 2 && funcType.NumOut() == 1 &&
		funcType.In(0) == funcType.Out(0) &&
		funcType.In(1) == sliceType.Elem()
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
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
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
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
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
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
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
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
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
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
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

// Filter returns a new slice, choosing only the elements for which the given
// function returns true.
func (ts TrickSlice) Filter(fn interface{}) TrickSlice {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !f.IsValid() || !isValidBoolFunc(f.Type(), v.Type()) {
		panic("tricks: slice.Filter: invalid function type")
	}
	typ := v.Type()

	out := reflect.MakeSlice(typ, 0, 0)
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if f.Call([]reflect.Value{val})[0].Bool() {
			out = reflect.Append(out, val)
		}
	}

	return TrickSlice(out)
}

// Map applies the given function to each element of the slice and stores the
// result to a new slice. The cap() of the new slice is set to equal its length.
func (ts TrickSlice) Map(fn interface{}) TrickSlice {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !f.IsValid() || !isValidMapFunc(f.Type(), v.Type()) {
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

// Reduce applies the given function to the values of the slice and reduces them
// down to a single value. fn should be `func(X, T) X`, zero should be type X.
func (ts TrickSlice) Reduce(zero, fn interface{}) interface{} {
	// TODO: Improve those docs above.
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !f.IsValid() || !isValidReduceFunc(f.Type(), v.Type()) {
		panic("tricks: slice.Reduce: invalid function type")
	}
	outType := f.Type().Out(0)
	z := reflect.ValueOf(zero)
	if !z.IsValid() {
		z = reflect.Zero(outType)
	}
	if z.Type() != outType {
		panic("tricks: slice.Reduce: invalid zero type")
	}

	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		z = f.Call([]reflect.Value{z, val})[0]
	}

	return z.Interface()
}

// GroupBy collects the slice values into a map, where the keys are the return
// value of the grouping function and the values are slices of elements that
// correspond to that key.
func (ts TrickSlice) GroupBy(fn interface{}) TrickMap {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !f.IsValid() || !isValidMapFunc(f.Type(), v.Type()) {
		panic("tricks: slice.GroupBy: invalid function type")
	}
	keyType := f.Type().Out(0)
	valType := v.Type()
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
