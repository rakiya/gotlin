package stream

import r "reflect"

func valueAt(v interface{}, i int) interface{} {
	return r.Indirect(r.ValueOf(v)).Index(i).Interface()
}

func fieldNameOf(v interface{}, i int) string {
	return r.ValueOf(v).Type().Field(i).Name
}

func fieldOf(v interface{}, i int) interface{} {
	return r.ValueOf(v).Field(i).Interface()
}

func dereference(v interface{}) interface{} {
	return r.Indirect(r.ValueOf(v)).Interface()
}

func lenOf(v interface{}) int {
	return r.Indirect(r.ValueOf(v)).Len()
}
