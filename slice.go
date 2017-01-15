package tricks

import "reflect"

type TrickSlice reflect.Value

var (
	typeTrickSlice = reflect.TypeOf((*TrickSlice)(nil)).Elem()  // TrickSlice
	typeInterface  = reflect.TypeOf((*interface{})(nil)).Elem() // interface{}
)

func Slice(sliceOrElements ...interface{}) TrickSlice {
	s := reflect.ValueOf(sliceOrElements) // []interface{}

	if len(sliceOrElements) == 0 {
		return TrickSlice(s)
	}

	if len(sliceOrElements) == 1 {
		v := reflect.ValueOf(sliceOrElements[0])
		if !v.IsValid() { // nil
			return TrickSlice(s)
		}
		if v.Kind() == reflect.Slice { // TrickSlice is a Kind of reflect.Struct
			return TrickSlice(v)
		}
		if v.Type() == typeTrickSlice {
			return sliceOrElements[0].(TrickSlice)
		}
	}

	// Try to make a typed slice from variadic args. First identify the type.
	var typ reflect.Type
	for i := 0; i < len(sliceOrElements); i++ {
		el := reflect.ValueOf(sliceOrElements[i])
		if el.IsValid() {
			switch typ {
			case nil:
				typ = el.Type()
				if typ == typeInterface {
					return TrickSlice(s)
				}
			case el.Type():
			default: // mixed types; fall back to []interface{}
				return TrickSlice(s)
			}
		}
	}
	if typ == nil { // no IsValid (non-nil) values found
		return TrickSlice(s)
	}
	// Now make the slice.
	sliceTyp := reflect.SliceOf(typ)
	slice := reflect.MakeSlice(sliceTyp, len(sliceOrElements), len(sliceOrElements))
	for i := 0; i < len(sliceOrElements); i++ {
		el := reflect.ValueOf(sliceOrElements[i])
		if el.IsValid() {
			slice.Index(i).Set(el)
		}
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

// IsEmpty returns true if the slice has no length, else false.
func (ts TrickSlice) IsEmpty() bool {
	return ts.Len() == 0
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

// Push appends a single element to the end of the slice.
func (ts *TrickSlice) Push(element interface{}) {
	in := reflect.Value(*ts)
	out := reflect.Append(in, reflect.ValueOf(element))
	*ts = TrickSlice(out)
}

// Pop removes the last element from the slice and returns it. If the slice is
// empty, this method returns a nil value.
func (ts *TrickSlice) Pop() interface{} {
	if ts.IsEmpty() {
		return nil
	}
	v := reflect.Value(*ts)
	*ts = TrickSlice(v.Slice(0, v.Len()-1))
	return v.Index(v.Len() - 1).Interface()
}

// Shift removes the first element from the slice and returns it. If the slice
// is empty, this method returns a nil value.
func (ts *TrickSlice) Shift() interface{} {
	if ts.IsEmpty() {
		return nil
	}
	v := reflect.Value(*ts)
	*ts = TrickSlice(v.Slice(1, v.Len()))
	return v.Index(0).Interface()
}

// Unshift prepends a single element to the start of the slice.
func (ts *TrickSlice) Unshift(element interface{}) {
	ts.Insert(element, 0) // will never panic
}

// TODO: Refactor to reuse code Push/Unshift -> Insert, Pop/Shift -> Delete

// Insert inserts an element at the given position in the slice. Slices are
// indexed from 0.
func (ts *TrickSlice) Insert(element interface{}, n int) {
	in := reflect.Value(*ts)
	if n < 0 || n > in.Len() {
		panic("tricks: slice.Insert: index out of bounds")
	}
	v := reflect.ValueOf(element)
	out := reflect.Append(in, v) // Grow as required
	// Shift everything up
	reflect.Copy(out.Slice(n+1, out.Len()), out.Slice(n, out.Len()-1))
	out.Index(n).Set(v)
	*ts = TrickSlice(out)
}

// Delete removes an element at the given position in the slice. Note that this
// does an internal copy to preserve the order of elements in the slice. Slices
// are indexed from 0.
func (ts *TrickSlice) Delete(n int) {
	v := reflect.Value(*ts)
	if n < 0 || n >= v.Len() {
		panic("tricks: slice.Delete: index out of bounds")
	}
	if n == 0 { // Special case, no copy needed.
		ts.Shift()
		return
	}
	// TODO: Copy from the shorter half and reslice from either end.
	reflect.Copy(v.Slice(n, v.Len()-1), v.Slice(n+1, v.Len()))
	last := v.Index(v.Len() - 1)
	last.Set(reflect.Zero(last.Type()))
	*ts = TrickSlice(v.Slice(0, v.Len()-1))
}

// Flatten returns a new slice of values, recursively extracting the elements
// from any nested slices. This new slice tries to take on the type of the first
// non-slice element encountered. If the elements are of mixed types, it falls
// back to []interface{}. nil values are treated as zeroes of the common type.
func (ts TrickSlice) Flatten() TrickSlice {
	in := reflect.Value(ts)

	var typ reflect.Type
	var vals []reflect.Value

	var extract, extractSlice func(reflect.Value)
	extract = func(el reflect.Value) {
		if el.Kind() == reflect.Slice {
			extractSlice(el)
			return
		}
		if el.IsValid() {
			switch el.Type() {
			case typeTrickSlice:
				extractSlice(reflect.Value(el.Interface().(TrickSlice)))
				return
			case typeInterface:
				extract(reflect.ValueOf(el.Interface()))
				return
			default:
				if typ == nil {
					typ = el.Type()
				} else if typ != typeInterface && typ != el.Type() {
					typ = typeInterface // fall back to []interface{}
				}
			}
		}
		vals = append(vals, el)
	}
	extractSlice = func(slice reflect.Value) {
		// Invariant: slice.Type().Kind() == reflect.Slice
		for i := 0; i < slice.Len(); i++ {
			extract(slice.Index(i))
		}
	}
	extractSlice(in)

	if typ == nil { // no IsValid (non-nil) values found
		typ = typeInterface
	}

	out := reflect.MakeSlice(reflect.SliceOf(typ), len(vals), len(vals))
	for i := 0; i < out.Len(); i++ {
		if vals[i].IsValid() {
			out.Index(i).Set(vals[i])
		}
	}
	return TrickSlice(out)
}
