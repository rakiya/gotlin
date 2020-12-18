package stream

import (
	r "reflect"
	"strconv"
	"testing"
)

func Test_NewStreamBasic(t *testing.T) {
	testPatterns := []struct {
		testCase string
		data     interface{}
	}{
		{testCase: "int", data: 32},
		{testCase: "float", data: 3.14},
		{testCase: "bool", data: false},
	}

	for _, p := range testPatterns {
		t.Run(p.testCase, func(t *testing.T) {
			s := NewStream(p.data)

			if s.entity != p.data {
				t.Fatalf("Expected %v, but %v", p.data, s.entity)
			}

			if s.entity == &(p.data) {
				t.Fatalf("The pointers store same address, %v", p.data)
			}
		})
	}
}

func Test_NewStreamPointer(t *testing.T) {
	testPatterns := []struct {
		testCase string
		data     interface{}
	}{
		{testCase: "*int", data: 11},
		{testCase: "*float", data: 2.71},
		{testCase: "*[]string", data: []string{"fiz", "buz"}},
	}

	for _, p := range testPatterns {
		t.Run(p.testCase, func(t *testing.T) {
			s := NewStream(&(p.data))

			if s.entity != &(p.data) {
				t.Fatalf("The pointers do NOT store same address. Expecte %v, but %v", s.entity, &(p.data))
			}
		})
	}
}

func Test_NewStreamArray(t *testing.T) {
	testPatterns := []struct {
		testCase string
		data     interface{}
	}{
		{testCase: "[1]int", data: [1]int{1}},
		{testCase: "[2]float32", data: [2]float32{3.3333, 0.5}},
		{testCase: "[3]string", data: [3]string{"fiz", "buz", "fizbuz"}},
	}

	for _, p := range testPatterns {
		t.Run(p.testCase, func(t *testing.T) {
			s := NewStream(p.data)

			for i, l := 0, r.ValueOf(p.data).Len(); i < l; i++ {
				if valueAt(s.entity, i) != valueAt(p.data, i) {
					t.Fatalf("data[%v] Expected %v, but %v", i, valueAt(p.data, i), valueAt(s.entity, i))
				}
			}
		})
	}
}

func Test_NewStreamStruct(t *testing.T) {
	type Example1 struct {
		A int
		B float32
		C *int
	}

	type Example2 struct {
		A *[]int
	}

	ex1C := 2

	testPatterns := []struct {
		testCase string
		data     interface{}
	}{
		{testCase: "Example1", data: Example1{A: 1, B: 3.14, C: &ex1C}},
		{testCase: "Example2", data: Example2{A: &[]int{1, 2, 3}}},
	}

	for _, p := range testPatterns {
		t.Run(p.testCase, func(t *testing.T) {
			s := NewStream(p.data)

			for i, l := 0, r.ValueOf(p.data).NumField(); i < l; i++ {
				if fieldOf(s.entity, i) != fieldOf(p.data, i) {
					t.Fatalf("data.%va Expected %v, but %v", fieldNameOf(p.data, i), fieldOf(p.data, i), fieldOf(s.entity, i))
				}
			}
		})
	}
}

func Test_LetSuccessBasic(t *testing.T) {
	testCase := []struct {
		caseName string
		data     interface{} // input of Let function
		fn       interface{} // fn of Let function
		expected interface{} // output of Let function
	}{
		{
			caseName: "int->int",
			data:     32,
			fn:       func(it int) int { return 2 * it },
			expected: 64,
		},
		{
			caseName: "string->string",
			data:     "Hello",
			fn:       func(it string) string { return "[" + it + "]" },
			expected: "[Hello]",
		},
		{
			caseName: "string->int",
			data:     "101",
			fn:       func(it string) int { tmp, _ := strconv.Atoi(it); return tmp },
			expected: 101,
		},
	}

	for _, p := range testCase {
		t.Run(p.caseName, func(t *testing.T) {
			s := NewStream(p.data).Let(p.fn)

			if s.entity != p.expected {
				t.Errorf("Expected %v, but %v", p.expected, s.entity)
			}
		})
	}
}

func Test_LetFailure_InvalidFnIn(t *testing.T) {
	testCase := []struct {
		caseName string
		data     interface{} // input of Let function
		fn       interface{} // fn of Let function
		expected interface{} // output of Let function
	}{
		{
			caseName: "int/float64",
			data:     32,
			fn:       func(it float64) int { return int(2.3 * it) },
			expected: 64,
		},
		{
			caseName: "string/[]int",
			data:     "Hello",
			fn:       func(it []int) string { return strconv.Itoa(it[0]) },
			expected: "[Hello]",
		},
		{
			caseName: "string/*int",
			data:     "101",
			fn:       func(it *int) int { return *it },
			expected: 101,
		},
	}

	for _, p := range testCase {
		t.Run(p.caseName, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()

			NewStream(&(p.data)).Let(p.fn)
		})
	}
}

func Test_Apply_Success(t *testing.T) {
	data := 32
	s := NewStream(&data).Apply(func(it *int) { *it *= 2 })

	if dereference(s.entity) != 64 {
		t.Errorf("Expected %v, but %v", 64, s.entity)
	}
}

func Test_Apply_InvalidFnIn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	data := 32
	NewStream(&data).Apply(func(it *string) { *it = "[" + *it + "]" })
}

func Test_Apply_InvalidFnOut(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	data := 32
	NewStream(&data).Apply(func(it *int) int { return 2 * *it })
}

func Test_Map_Success_Basic(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}

	s := NewStream(data).
		Map(func(it int) int { return 2 * it })

	for i := range data {
		if data[i]*2 != valueAt(s.entity, i) {
			t.Errorf("Expected %v, but %v", data[i]*2, valueAt(s.entity, i))
		}
	}
}

func Test_Map_Success_Basic_Ptr(t *testing.T) {
	d1 := 1
	d2 := 2
	data := []*int{&d1, &d2}

	s := NewStream(data).Map(func(it *int) int { return 2 * *it })

	for i := range data {
		if *data[i]*2 != valueAt(s.entity, i) {
			t.Errorf("Expected %v, but %v", *data[i]*2, valueAt(s.entity, i))
		}
	}
}

func Test_Map_InvalidEntityType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewStream(1).Map(func(it int) int { return it * 2 })
}

func Test_Map_Failure_InvalidFnIn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	NewStream([]int{1, 2, 3, 4, 5}).Map(func(it string) string { return it + "." })
}

func Test_Filter_Success(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4}

	s := NewStream(data).Filter(func(it int) bool { return it%2 == 0 })

	if len(expected) != lenOf(s.entity) {
		t.Errorf("Expected %v, but %v", len(expected), len(s.entity.([]int)))
	}

	for i := range expected {
		if valueAt(s.entity, i) != expected[i] {
			t.Errorf("At s.entity[%v], expected %v but %v", i, expected[i], valueAt(s.entity, i))
		}
	}
}

func Test_Filter_InvalidFnIn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		}
	}()

	NewStream([]int{1, 2, 3, 4, 5}).Filter(func(it *int) bool { return *it%2 == 0 })
}

func Test_Filter_InvalidEntityType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		}
	}()

	NewStream(1).Filter(func(it *int) bool { return *it%2 == 0 })
}

func Test_Inject_Success(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	expected := 15

	s := NewStream(data).Inject(0, func(res int, it int) int { return res + it })

	if s.entity != expected {
		t.Errorf("Expected %v but %v", expected, s.entity)
	}
}

func Test_Inject_InvalidInit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		}
	}()

	NewStream([]int{1, 2, 3, 4}).Inject(float32(0), func(res int, it int) int { return res + it })
}

func Test_Inject_InvalidFnIn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		}
	}()

	NewStream([]int{1, 2, 3, 4}).Inject(0, func(it int) int { return it })
}

func Test_Inject_InvalidFnOut(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		}
	}()

	NewStream([]int{1, 2, 3, 4}).Inject(0, func(res int, it int) float64 { return float64(res + it) })
}
