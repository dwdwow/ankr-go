package ankr

import (
	"fmt"
	"reflect"
	"strconv"
)

// DefaultTag is the struct tag key for default values
const DefaultTag = "default"

// ApplyDefaults applies default values to a struct based on "default" tags
// This allows us to use non-pointer fields in Opts structs and still have optional parameters
func ApplyDefaults[T any](opts T) (result T, err error) {
	v := reflect.ValueOf(opts)
	isPointer := v.Kind() == reflect.Pointer
	if isPointer {
		v = v.Elem()
	} else {
		v = reflect.ValueOf(&opts).Elem()
	}

	if v.Kind() != reflect.Struct {
		return result, fmt.Errorf("opts must be a struct or pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get the default tag value
		defaultValue := fieldType.Tag.Get(DefaultTag)
		if defaultValue == "" {
			continue
		}

		// Skip if field already has a value (not zero value)
		if !field.IsZero() {
			continue
		}

		// Apply default value based on field type
		if err := setFieldValue(field, defaultValue); err != nil {
			return result, err
		}
	}

	return opts, nil
}

// setFieldValue sets a field value from a string representation
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)

	default:
		panic(fmt.Sprintf("ankr: setFieldValue: unsupported field type: %s", field.Kind()))
	}

	return nil
}
