package tricks

import "reflect"

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
