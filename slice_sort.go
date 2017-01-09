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

// Sort the contents of the slice in place. Slices of type string, int, or float
// are handled automatically, otherwise, the underlying slice must implement
// sort.Interface or this method panics.
func (ts TrickSlice) Sort() TrickSlice {
	v := reflect.Value(ts)
	t := v.Type().Elem()

	switch t.Kind() {
	case reflect.String:
		s := sortableStrings(v)
		sort.Sort(s)
		return TrickSlice(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := sortableInts(v)
		sort.Sort(s)
		return TrickSlice(s)
	case reflect.Float32, reflect.Float64:
		s := sortableFloats(v)
		sort.Sort(s)
		return TrickSlice(s)
	default:
		// if !keyType.Implements(reflect.TypeOf((*sort.Interface)(nil))) { // .Elem()?
		panic("slice not sortable")
		// }
	}
}
