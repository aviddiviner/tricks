package tricks

import (
	"reflect"
)

type TrickMap reflect.Value

func Map(anyMap interface{}) TrickMap {
	v := reflect.ValueOf(anyMap)
	if v.Kind() != reflect.Map {
		panic("Map: input is not a map")
	}
	return TrickMap(v)
}

func (tm TrickMap) Value() interface{} {
	return reflect.Value(tm).Interface()
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

// Only returns a map containing only those keys given by the variadic argument.
func (tm TrickMap) Only(keys ...interface{}) TrickMap {
	v := reflect.Value(tm)

	keyType := v.Type().Key()
	mapType := reflect.MapOf(keyType, v.Type().Elem())

	out := reflect.MakeMap(mapType)
	for i := 0; i < len(keys); i++ {
		key := reflect.ValueOf(keys[i])
		if !key.Type().AssignableTo(keyType) {
			panic("map.Only: key doesn't match map's key type")
		}
		out.SetMapIndex(key, v.MapIndex(key))
	}

	return TrickMap(out)
}
