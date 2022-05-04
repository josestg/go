package flatten

import (
	"fmt"
	"reflect"
	"strings"
)

// Any aliases of interface{}.
type Any = interface{}

// Flatten transforms Slice and Map into flat Map.
func Flatten(obj Any, sep string) Any {
	rv := reflect.ValueOf(obj)
	switch rv.Kind() {
	default:
		return obj
	case reflect.Map, reflect.Slice, reflect.Array:
		var dst = make(map[string]Any)
		flatten(rv, dst, []string{}, sep)
		return dst
	}
}

func flatten(rv reflect.Value, dst map[string]Any, prefixes []string, sep string) {
	switch rv.Kind() {
	case reflect.Map:
		for _, field := range rv.MapKeys() {
			child := rv.MapIndex(field)
			if !child.IsZero() {
				nextPrefixes := append(prefixes, fmt.Sprint(field.String()))
				flatten(reflect.ValueOf(child.Interface()), dst, nextPrefixes, sep)
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			child := rv.Index(i)
			if !child.IsZero() {
				nextPrefixes := append(prefixes, fmt.Sprint(i))
				flatten(reflect.ValueOf(child.Interface()), dst, nextPrefixes, sep)
			}
		}
	default:
		flatKey := strings.Join(prefixes, sep)
		dst[flatKey] = rv.Interface()
	}
}
