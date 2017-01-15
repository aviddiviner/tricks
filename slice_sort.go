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
type sortableIface reflect.Value

func (p sortableInts) Len() int    { return reflect.Value(p).Len() }
func (p sortableFloats) Len() int  { return reflect.Value(p).Len() }
func (p sortableStrings) Len() int { return reflect.Value(p).Len() }

func (p sortableIface) Len() int {
	f := reflect.Value(p).MethodByName("Len")
	return f.Call(nil)[0].Interface().(int)
}

func (p sortableInts) Swap(i, j int)    { swapValue(reflect.Value(p), i, j) }
func (p sortableFloats) Swap(i, j int)  { swapValue(reflect.Value(p), i, j) }
func (p sortableStrings) Swap(i, j int) { swapValue(reflect.Value(p), i, j) }

func (p sortableIface) Swap(i, j int) {
	f := reflect.Value(p).MethodByName("Swap")
	args := []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)}
	f.Call(args)
}

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

func (p sortableIface) Less(i, j int) bool {
	f := reflect.Value(p).MethodByName("Less")
	args := []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)}
	return f.Call(args)[0].Bool()
}

// We declare these types because we want to check that the slice is of this
// exact type, and not just having elements of the same reflect.Kind. In this
// way, we won't sort elements that look like ints, unless that slice type is
// either exactly []int or it explicitly implements sort.Interface.
var (
	typeIntSlice      = reflect.SliceOf(reflect.TypeOf((*int)(nil)).Elem())     // []int
	typeInt8Slice     = reflect.SliceOf(reflect.TypeOf((*int8)(nil)).Elem())    // []int8
	typeInt16Slice    = reflect.SliceOf(reflect.TypeOf((*int16)(nil)).Elem())   // []int16
	typeInt32Slice    = reflect.SliceOf(reflect.TypeOf((*int32)(nil)).Elem())   // []int32
	typeInt64Slice    = reflect.SliceOf(reflect.TypeOf((*int64)(nil)).Elem())   // []int64
	typeStringSlice   = reflect.SliceOf(reflect.TypeOf((*string)(nil)).Elem())  // []string
	typeFloat32Slice  = reflect.SliceOf(reflect.TypeOf((*float32)(nil)).Elem()) // []float32
	typeFloat64Slice  = reflect.SliceOf(reflect.TypeOf((*float64)(nil)).Elem()) // []float64
	typeSortInterface = reflect.TypeOf((*sort.Interface)(nil)).Elem()           // sort.Interface
)

func getSortable(context string, v reflect.Value) sort.Interface {
	switch v.Type() {
	case typeStringSlice:
		return sortableStrings(v)
	case typeIntSlice, typeInt8Slice, typeInt16Slice, typeInt32Slice, typeInt64Slice:
		return sortableInts(v)
	case typeFloat32Slice, typeFloat64Slice:
		return sortableFloats(v)
	default:
		if v.Type().Implements(typeSortInterface) {
			return sortableIface(v)
		}
		panic("tricks: " + context + ": slice doesn't implement sort.Interface")
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

type sortableBy struct {
	val    reflect.Value
	sortFn interface{}
}

func (s *sortableBy) Len() int {
	return s.val.Len()
}

func (s *sortableBy) Swap(i, j int) {
	swapValue(s.val, i, j)
}

func (s *sortableBy) Less(i, j int) bool {
	fn := reflect.ValueOf(s.sortFn)
	return fn.Call([]reflect.Value{s.val.Index(i), s.val.Index(j)})[0].Bool()
}

func isValidSortByFunc(funcType, sliceType reflect.Type) bool {
	return funcType.NumIn() == 2 && funcType.NumOut() == 1 &&
		funcType.In(0) == sliceType.Elem() &&
		funcType.In(1) == sliceType.Elem() &&
		funcType.Out(0).Kind() == reflect.Bool
}

// SortBy sorts the slice values by some by some `func(a, b T) bool` that
// returns whether element `a < b`.
func (ts TrickSlice) SortBy(fn interface{}) TrickSlice {
	v := reflect.Value(ts)
	f := reflect.ValueOf(fn)
	if !f.IsValid() || !isValidSortByFunc(f.Type(), v.Type()) {
		panic("tricks: slice.SortBy: invalid function type")
	}
	sort.Sort(&sortableBy{v, fn})
	return ts
}
