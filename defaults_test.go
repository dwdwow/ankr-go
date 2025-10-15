package ankr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestApplyDefaults(t *testing.T) {
	type TestStruct struct {
		Field1 string  `default:"test"`
		Field2 int     `default:"10"`
		Field3 bool    `default:"true"`
		Field4 float64 `default:"3.14"`
	}

	testStruct := &TestStruct{}
	testStruct, err := ApplyDefaults(testStruct)
	if err != nil {
		t.Fatalf("Failed to apply defaults: %v", err)
	}

	t.Logf("TestStruct: %+v", testStruct)
}

func TestX(t *testing.T) {
	s := &[]bool{false}[0]
	v := reflect.ValueOf(s)
	fmt.Println(v.Kind())
	v = v.Elem()
	fmt.Println(v.Kind())
	v.SetBool(true)
	fmt.Println(*s)
}
