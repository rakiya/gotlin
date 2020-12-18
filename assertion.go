package stream

import (
	"fmt"
	"reflect"
)

// assertFn validates that the type of argument fnv is Func
func assertFn(fnv reflect.Value) {
	if fnv.Kind() != reflect.Func {
		panic(fmt.Sprintf("The argument, fn should be a function, but %v", fnv.Kind().String()))
	}
}

// assertFnIn validates arguments' type of fnv whose type is function.
func assertFnIn(fnv reflect.Value, inTypes []reflect.Type) {
	if fnv.Type().NumIn() != len(inTypes) {
		panic(fmt.Sprintf("The number of arguments of fn should be %v, but %v", len(inTypes), fnv.Type().NumIn()))
	}

	for i := 0; i < fnv.Type().NumIn(); i++ {
		if fnv.Type().In(i) != inTypes[i] {
			panic(fmt.Sprintf("The argument type at %v of fn should be %v, but %v", i, inTypes[i], fnv.Type().In(i)))
		}
	}
}

// assertFnOut validates return's type of fnv whose type is function.
func assertFnOut(fnv reflect.Value, outTypes []reflect.Type) {
	if fnv.Type().NumOut() != len(outTypes) {
		panic(fmt.Sprintf("The number of arguments of fn should be %v, but %v", len(outTypes), fnv.Type().NumOut()))
	}

	for i := 0; i < fnv.Type().NumOut(); i++ {
		if fnv.Type().Out(i) != outTypes[i] {
			panic(fmt.Sprintf("The return type at %v of fn should be %v, but %v", i, outTypes[i], fnv.Type().Out(i)))
		}
	}
}

// assertFnAndFnIn executes assertFn and assertFnIn
func assertFnAndFnIn(fnv reflect.Value, inTypes []reflect.Type) {
	assertFn(fnv)
	assertFnIn(fnv, inTypes)
}

// assertFnAndFnIn executes assertFn and assertFnOut
func assertFnAndFnOut(fnv reflect.Value, outTypes []reflect.Type) {
	assertFn(fnv)
	assertFnOut(fnv, outTypes)
}

// assertFnAndFnIn executes assertFn, assertFnIn and assertFnOut
func assertFnAndFnInAndFnOut(fnv reflect.Value, inTypes []reflect.Type, outTypes []reflect.Type) {
	assertFn(fnv)
	assertFnIn(fnv, inTypes)
	assertFnOut(fnv, outTypes)
}
