package tricks

import (
	"reflect"
	"sort"
)

type sortableStrings reflect.Value

func (p sortableStrings) Len() int {
	return reflect.Value(p).Len()
}

func (p sortableStrings) Swap(i, j int) {
	v := reflect.Value(p)
	vi := v.Index(i)
	v2 := reflect.ValueOf(vi.Interface())
	vj := v.Index(j)
	vi.Set(vj)
	vj.Set(v2)
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
		sortable := sortableStrings(v)
		sort.Sort(sortable)
		return TrickSlice(sortable)
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 	sortable := sortableInts(m)
	// case reflect.Float32, reflect.Float64:
	// 	sortable := sortableFloats(m)
	default:
		// if !keyType.Implements(reflect.TypeOf((*sort.Interface)(nil))) { // .Elem()?
		panic("slice not sortable")
		// }
	}
}
