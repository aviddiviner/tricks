package tricks

import (
	"math"
	"reflect"
	"sort"
)

type sortableInts reflect.Value
type sortableFloats reflect.Value
type sortableStrings reflect.Value

func (p sortableInts) Len() int    { return reflect.Value(p).Len() }
func (p sortableFloats) Len() int  { return reflect.Value(p).Len() }
func (p sortableStrings) Len() int { return reflect.Value(p).Len() }

func swapValue(v reflect.Value, i, j int) {
	vi := v.Index(i)
	vii := reflect.ValueOf(vi.Interface())
	vj := v.Index(j)
	vi.Set(vj)
	vj.Set(vii)
}

func (p sortableInts) Swap(i, j int)    { swapValue(reflect.Value(p), i, j) }
func (p sortableFloats) Swap(i, j int)  { swapValue(reflect.Value(p), i, j) }
func (p sortableStrings) Swap(i, j int) { swapValue(reflect.Value(p), i, j) }

func (p sortableInts) Less(i, j int) bool {
	v := reflect.Value(p)
	return v.Index(i).Int() < v.Index(j).Int()
}

func (p sortableFloats) Less(i, j int) bool {
	v := reflect.Value(p)
	vi := v.Index(i).Float()
	vj := v.Index(j).Float()
	return vi < vj || math.IsNaN(vi) && !math.IsNaN(vj)
}

func (p sortableStrings) Less(i, j int) bool {
	v := reflect.Value(p)
	return v.Index(i).String() < v.Index(j).String()
}

type sortableIface reflect.Value

func (p sortableIface) Len() int {
	f := reflect.Value(p).MethodByName("Len")
	return f.Call(nil)[0].Interface().(int)
}

func (p sortableIface) Swap(i, j int) {
	f := reflect.Value(p).MethodByName("Swap")
	args := []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)}
	f.Call(args)
}

func (p sortableIface) Less(i, j int) bool {
	f := reflect.Value(p).MethodByName("Less")
	args := []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)}
	return f.Call(args)[0].Interface().(bool)
}

var sortInterfaceType = reflect.TypeOf((*sort.Interface)(nil)).Elem()

// Sort the contents of the slice in place. Slices of type string, int, or float
// are handled automatically, otherwise, the underlying slice must implement
// sort.Interface or this method panics.
func (ts TrickSlice) Sort() TrickSlice {
	v := reflect.Value(ts)
	t := v.Type().Elem()

	switch t.Kind() {
	case reflect.String:
		sort.Sort(sortableStrings(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sort.Sort(sortableInts(v))
	case reflect.Float32, reflect.Float64:
		sort.Sort(sortableFloats(v))
	default:
		if !v.Type().Implements(sortInterfaceType) {
			panic("slice doesn't implement sort.Interface")
		}
		sort.Sort(sortableIface(v))
	}
	return ts
}

func findIndexMax(s sort.Interface) (max int) {
	if s.Len() == 1 {
		return
	}
	for i := 1; i < s.Len(); i++ {
		if s.Less(max, i) {
			max = i
		}
	}
	return
}

func findIndexMin(s sort.Interface) (min int) {
	if s.Len() == 1 {
		return
	}
	for i := 1; i < s.Len(); i++ {
		if s.Less(i, min) {
			min = i
		}
	}
	return
}

// Max returns the element of the slice with the maximum value. Slices of type
// string, int, or float are handled automatically, otherwise, the underlying
// slice must implement sort.Interface or this method panics.
// If the slice is empty, this method returns the nil interface{}.
func (ts TrickSlice) Max() interface{} {
	v := reflect.Value(ts)
	t := v.Type().Elem()

	if v.Len() == 0 {
		return nil
	}

	var max int
	switch t.Kind() {
	case reflect.String:
		max = findIndexMax(sortableStrings(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max = findIndexMax(sortableInts(v))
	case reflect.Float32, reflect.Float64:
		max = findIndexMax(sortableFloats(v))
	default:
		if !v.Type().Implements(sortInterfaceType) {
			panic("slice doesn't implement sort.Interface")
		}
		max = findIndexMax(sortableIface(v))
	}
	return v.Index(max).Interface()
}

// Min returns the element of the slice with the minimum value. Slices of type
// string, int, or float are handled automatically, otherwise, the underlying
// slice must implement sort.Interface or this method panics.
// If the slice is empty, this method returns the nil interface{}.
func (ts TrickSlice) Min() interface{} {
	v := reflect.Value(ts)
	t := v.Type().Elem()

	if v.Len() == 0 {
		return nil
	}

	var min int
	switch t.Kind() {
	case reflect.String:
		min = findIndexMin(sortableStrings(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min = findIndexMin(sortableInts(v))
	case reflect.Float32, reflect.Float64:
		min = findIndexMin(sortableFloats(v))
	default:
		if !v.Type().Implements(sortInterfaceType) {
			panic("slice doesn't implement sort.Interface")
		}
		min = findIndexMin(sortableIface(v))
	}
	return v.Index(min).Interface()
}
