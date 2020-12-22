package stream

import (
	"reflect"
)

// Stream is a wrapper to use golin functions
type Stream struct {
	entity interface{}
}

// NewStream function creates a Stream and return the pointer
func NewStream(entity interface{}) *Stream {
	return &Stream{entity}
}

// End function unwrap the result from Stream wrapper as interface{}
func (s *Stream) End() interface{} {
	return s.entity
}

// Let function is like the let function in Kotlin.
// Let receives a function fn that makes a value to another value
// and applys it to stream's entity.
func (s *Stream) Let(fn interface{}) *Stream {
	entityv, fnv := reflect.ValueOf(s.entity), reflect.ValueOf(fn)

	assertFnAndFnIn(fnv, []reflect.Type{reflect.TypeOf(s.entity)})

	newEntity := fnv.Call([]reflect.Value{entityv})[0].Interface()
	return &Stream{newEntity}
}

// Apply function is like the apply function in Kotlin.
// Apply receives a function fn that works to its argument
// and applys it with stream's entity as its argument.
func (s *Stream) Apply(fn interface{}) *Stream {
	entityv, fnv := reflect.ValueOf(s.entity), reflect.ValueOf(fn)

	assertFnAndFnInAndFnOut(fnv, []reflect.Type{entityv.Type()}, []reflect.Type{})

	fnv.Call([]reflect.Value{entityv})
	return s
}

// ForEach function applys fn to each of entities.
func (s *Stream) ForEach(fn interface{}) *Stream {
	entityv, fnv := reflect.ValueOf(s.entity), reflect.ValueOf(fn)

	if reflect.TypeOf(s.entity).Kind() == reflect.Ptr {
		assertFnAndFnInAndFnOut(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem().Elem()}, []reflect.Type{})
	} else {
		assertFnAndFnInAndFnOut(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem()}, []reflect.Type{})
	}

	for i := 0; i < entityv.Len(); i++ {
		fnv.Call([]reflect.Value{entityv.Index(i)})
	}

	return s
}

// Map function is like a map function in a functional programming.
func (s *Stream) Map(fn interface{}) *Stream {
	entityv, fnv := reflect.Indirect(reflect.ValueOf(s.entity)), reflect.ValueOf(fn)

	// validate type of fn.
	if reflect.TypeOf(s.entity).Kind() == reflect.Ptr {
		assertFnAndFnIn(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem().Elem()})
	} else {
		assertFnAndFnIn(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem()})
	}

	resType := reflect.SliceOf(fnv.Type().Out(0))
	res, indirect := getSlicePtrOf(resType, entityv.Len(), entityv.Cap())

	for i := 0; i < entityv.Len(); i++ {
		indirect.Index(i).Set(fnv.Call([]reflect.Value{entityv.Index(i)})[0])
	}

	return &Stream{res.Interface()}
}

// Filter function is like a filter function in a functional programming.
func (s *Stream) Filter(fn interface{}) *Stream {
	entityv, fnv := reflect.Indirect(reflect.ValueOf(s.entity)), reflect.ValueOf(fn)

	// validate fn type
	if reflect.TypeOf(s.entity).Kind() == reflect.Ptr {
		assertFnAndFnIn(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem().Elem()})
	} else {
		assertFnAndFnIn(fnv, []reflect.Type{reflect.TypeOf(s.entity).Elem()})
	}

	resType := reflect.SliceOf(fnv.Type().In(0))
	res, indirect := getSlicePtrOf(resType, 0, entityv.Cap())
	idxIndirect := 0
	for i := 0; i < entityv.Len(); i++ {
		if fnv.Call([]reflect.Value{entityv.Index(i)})[0].Bool() {
			indirect.SetLen(idxIndirect + 1)
			indirect.Index(idxIndirect).Set(entityv.Index(i))
			idxIndirect++
		}
	}

	return &Stream{res.Interface()}
}

// Inject function is like a inject or reduce function in a functional programming.
func (s *Stream) Inject(init interface{}, fn interface{}) *Stream {
	entityv, initv, fnv := reflect.Indirect(reflect.ValueOf(s.entity)), reflect.ValueOf(init), reflect.ValueOf(fn)

	if reflect.TypeOf(s.entity).Kind() == reflect.Ptr {
		assertFnAndFnInAndFnOut(
			fnv,
			[]reflect.Type{initv.Type(), reflect.TypeOf(s.entity).Elem().Elem()},
			[]reflect.Type{initv.Type()},
		)
	} else {
		assertFnAndFnInAndFnOut(
			fnv,
			[]reflect.Type{initv.Type(), reflect.TypeOf(s.entity).Elem()},
			[]reflect.Type{initv.Type()},
		)
	}

	res := initv
	for i := 0; i < entityv.Len(); i++ {
		res = fnv.Call([]reflect.Value{res, entityv.Index(i)})[0]
	}

	return &Stream{res.Interface()}
}

// getSlicePtrOf creates a slice which has length len and capacity cap and whose type is t.
func getSlicePtrOf(t reflect.Type, len int, cap int) (ptr reflect.Value, indirect reflect.Value) {
	ptr = reflect.New(t)
	ptr.Elem().Set(reflect.MakeSlice(t, len, cap))
	indirect = reflect.Indirect(ptr)

	return ptr, indirect
}
