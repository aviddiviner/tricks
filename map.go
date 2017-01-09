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
