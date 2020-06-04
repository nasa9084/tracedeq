package tracedeq

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	TagKey         = "tracedeq"
	IgnoreTagValue = "ignore"
)

type Result struct {
	IsEqual bool
	Trace   Trace
	X, Y    interface{}
}

type Trace []string

func (t Trace) Join(sep string) string {
	return strings.Join(t, sep)
}

var expected = Result{
	IsEqual: true,
}

type visit struct {
	x, y uintptr
	typ  reflect.Type
}

func DeepEqual(x, y interface{}) Result {
	if x == nil || y == nil {
		return Result{
			IsEqual: x == y,
			X:       x,
			Y:       y,
		}
	}
	return deepValueEqual(reflect.ValueOf(x), reflect.ValueOf(y), make(map[visit]bool), nil)
}

func deepValueEqual(x, y reflect.Value, visited map[visit]bool, trace []string) Result {
	if x.Type() != y.Type() {
		return Result{
			IsEqual: false,
			Trace:   append(trace, "TYPE"),
			X:       x.Type(),
			Y:       y.Type(),
		}
	}

	// basically copied from reflect/deepequal.go
	// https://github.com/golang/go/blob/go1.14.3/src/reflect/deepequal.go
	hard := func(v1, v2 reflect.Value) bool {
		switch v1.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			// Nil pointers cannot be cyclic. Avoid putting them in the visited map.
			return !v1.IsNil() && !v2.IsNil()
		}
		return false
	}

	if hard(x, y) && x.CanAddr() && y.CanAddr() {
		addr1 := x.UnsafeAddr()
		addr2 := y.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := x.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return expected
		}

		// Remember for later.
		visited[v] = true
	}

	switch kind := x.Type().Kind(); kind {
	case reflect.Array:
		for i := 0; i < x.Len(); i++ {
			result := deepValueEqual(x.Index(i), y.Index(i), visited, append(trace, strconv.Itoa(i)))
			if !result.IsEqual {
				return result
			}
		}
	case reflect.Slice:
		if x.IsNil() != y.IsNil() {
			return newUnexpected(x, y, trace)
		}
		if x.Len() != y.Len() {
			return newUnexpectedLengthResult(x, y, trace)
		}
		if x.Pointer() == y.Pointer() {
			return expected
		}
		for i := 0; i < x.Len(); i++ {
			result := deepValueEqual(x.Index(i), y.Index(i), visited, append(trace, strconv.Itoa(i)))
			if !result.IsEqual {
				return result
			}
		}
		return expected
	case reflect.Interface:
		if (x.IsNil() || y.IsNil()) && x.IsNil() != y.IsNil() {
			return newUnexpected(x, y, trace)
		}
		return deepValueEqual(x.Elem(), y.Elem(), visited, trace)
	case reflect.Ptr:
		if x.Pointer() == y.Pointer() {
			return expected
		}
		return deepValueEqual(x.Elem(), y.Elem(), visited, trace)
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			if shouldIgnore(x.Type().Field(i), y.Type().Field(i)) {
				continue
			}
			if x.Field(i).IsZero() != y.Field(i).IsZero() {
				return newUnexpected(x.Field(i), y.Field(i), append(trace, x.Type().Field(i).Name))
			}
			if x.Field(i).IsZero() && y.Field(i).IsZero() {
				continue
			}
			result := deepValueEqual(x.Field(i), y.Field(i), visited, append(trace, x.Type().Field(i).Name))
			if !result.IsEqual {
				return result
			}
		}
		return expected
	case reflect.Map:
		if x.IsNil() != y.IsNil() {
			return newUnexpected(x, y, trace)
		}
		if x.Len() != y.Len() {
			return newUnexpectedLengthResult(x, y, trace)
		}
		if x.Pointer() == y.Pointer() {
			return expected
		}
		for _, key := range x.MapKeys() {
			v1 := x.MapIndex(key)
			v2 := y.MapIndex(key)
			if !v1.IsValid() || !v2.IsValid() {
				return newUnexpected(v1, v2, append(trace, key.String()))
			}
			result := deepValueEqual(v1, v2, visited, append(trace, key.String()))
			if !result.IsEqual {
				return result
			}
		}
		return expected
	case reflect.Func:
		if x.IsNil() && y.IsNil() {
			return expected
		}
		return newUnexpected(x, y, append(trace, "FUNC"))
	case reflect.Bool:
		v1 := x.Bool()
		v2 := y.Bool()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v1 := x.Int()
		v2 := y.Int()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v1 := x.Uint()
		v2 := y.Uint()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	case reflect.Float32, reflect.Float64:
		v1 := x.Float()
		v2 := y.Float()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	case reflect.Complex64, reflect.Complex128:
		v1 := x.Complex()
		v2 := y.Complex()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	case reflect.String:
		v1 := x.String()
		v2 := y.String()
		if v1 != v2 {
			return newUnexpected(v1, v2, trace)
		}
		return expected
	default:
		log.Print(x.String(), y.String())
		panic(fmt.Errorf("kind %s is not supported", kind))
	}
	return expected
}

func newUnexpected(x, y interface{}, trace []string) Result {
	return Result{
		IsEqual: false,
		Trace:   trace,
		X:       x,
		Y:       y,
	}
}

func newUnexpectedLengthResult(x, y reflect.Value, trace []string) Result {
	return newUnexpected(x.Len(), y.Len(), append(trace, "LENGTH"))
}

func shouldIgnore(x, y reflect.StructField) bool {
	xTags := strings.Split(x.Tag.Get(TagKey), ",")

	for _, tag := range xTags {
		if tag == IgnoreTagValue {
			return true
		}
	}

	yTags := strings.Split(x.Tag.Get(TagKey), ",")

	for _, tag := range yTags {
		if tag == IgnoreTagValue {
			return true
		}
	}

	return false
}
