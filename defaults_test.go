package ankr

import "testing"

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
