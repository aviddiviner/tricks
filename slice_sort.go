package tricks

import (
	"math"
	"reflect"
	"sort"
)

func swapValue(v reflect.Value, i, j int) {
	vi := v.Index(i)
	vii := reflect.ValueOf(vi.Interface())
	vj := v.Index(j)
	vi.Set(vj)
	vj.Set(vii)
}

// Reverse reverses the order of elements of the slice in place.
func (ts TrickSlice) Reverse() TrickSlice {
	v := reflect.Value(ts)
	for i, j := 0, v.Len()-1; i < j; i, j = i+1, j-1 {
		swapValue(v, i, j)
	}
	return ts
}

// Multiple types of sortable reflect.Values so that we can have independant
// Less() implementations, rather than one big switch statement inside Less().

type sortableInts reflect.Value
type sortableFloats reflect.Value
type sortableStrings reflect.Value

func (p sortableInts) Len() int    { return reflect.Value(p).Len() }
func (p sortableFloats) Len() int  { return reflect.Value(p).Len() }
func (p sortableStrings) Len() int { return reflect.Value(p).Len() }

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

func getSortable(context string, v reflect.Value) sort.Interface {
	switch v.Type().Elem().Kind() {
	case reflect.String:
		return sortableStrings(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return sortableInts(v)
	case reflect.Float32, reflect.Float64:
		return sortableFloats(v)
	default:
		if !v.Type().Implements(sortInterfaceType) {
			panic("tricks: " + context + ": slice doesn't implement sort.Interface")
		}
		return sortableIface(v)
	}
}

// Sort the contents of the slice in place. Slices of type string, int, or float
// are handled automatically, otherwise, the underlying slice must implement
// sort.Interface or this method panics.
func (ts TrickSlice) Sort() TrickSlice {
	v := reflect.Value(ts)
	sort.Sort(getSortable("slice.Sort", v))
	return ts
}

// Find the maximum value in O(n) time.
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

func findIndexMin(s sort.Interface) int {
	return findIndexMax(sort.Reverse(s))
}

// Max returns the element of the slice with the maximum value. Slices of type
// string, int, or float are handled automatically, otherwise, the underlying
// slice must implement sort.Interface or this method panics.
// If the slice is empty, this method returns the nil interface{}.
func (ts TrickSlice) Max() interface{} {
	v := reflect.Value(ts)
	if v.Len() == 0 {
		return nil
	}
	return v.Index(findIndexMax(getSortable("slice.Max", v))).Interface()
}

// Min returns the element of the slice with the minimum value. Slices of type
// string, int, or float are handled automatically, otherwise, the underlying
// slice must implement sort.Interface or this method panics.
// If the slice is empty, this method returns the nil interface{}.
func (ts TrickSlice) Min() interface{} {
	v := reflect.Value(ts)
	if v.Len() == 0 {
		return nil
	}
	return v.Index(findIndexMin(getSortable("slice.Min", v))).Interface()
}
