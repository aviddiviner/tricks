package tricks

import (
	"reflect"
)

type TrickMap reflect.Value

func Map(anyMap interface{}) TrickMap {
	v := reflect.ValueOf(anyMap)
	if v.Kind() != reflect.Map {
		panic("input is not a map")
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
