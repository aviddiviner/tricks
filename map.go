package tricks

import (
	"reflect"
)

type TrickMap reflect.Value

func Map(anyMap interface{}) TrickMap {
	v := reflect.ValueOf(anyMap)
	if !v.IsValid() { // nil
		v = reflect.ValueOf(map[interface{}]interface{}{})
	}
	if v.Kind() != reflect.Map {
		panic("tricks: Map: input is not a map")
	}
	return TrickMap(v)
}

func (tm TrickMap) Value() interface{} {
	return reflect.Value(tm).Interface()
}

// Copy returns a new map containing the same values.
func (tm TrickMap) Copy() TrickMap {
	in := reflect.Value(tm)
	out := reflect.MakeMap(in.Type())
	keys := in.MapKeys()
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		out.SetMapIndex(key, in.MapIndex(key))
	}
	return TrickMap(out)
}

// Len returns the length of the map (number of keys).
func (tm TrickMap) Len() int {
	return reflect.Value(tm).Len()
}

// Keys returns a slice of the map's keys.
func (tm TrickMap) Keys() TrickSlice {
	v := reflect.Value(tm)
	typ := reflect.SliceOf(v.Type().Key())

	k := v.MapKeys()
	out := reflect.MakeSlice(typ, len(k), len(k))
	for i := 0; i < len(k); i++ {
		out.Index(i).Set(k[i])
	}

	return TrickSlice(out)
}

// Only returns a map containing only those keys in the given slice. Also
// accepts a single key, or nil (return empty map). [2]
func (tm TrickMap) Only(keys interface{}) TrickMap {
	v := reflect.Value(tm)
	k := reflect.Value(Slice(keys))

	keyType := v.Type().Key()
	mapType := reflect.MapOf(keyType, v.Type().Elem())

	out := reflect.MakeMap(mapType)
	for i := 0; i < k.Len(); i++ {
		key := k.Index(i)
		if !key.Type().AssignableTo(keyType) {
			panic("tricks: map.Only: key doesn't match map's key type")
		}
		out.SetMapIndex(key, v.MapIndex(key))
	}

	return TrickMap(out)
}
