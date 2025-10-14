package ankr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// DefaultTag is the struct tag key for default values
const DefaultTag = "default"

// ApplyDefaults applies default values to a struct based on "default" tags
// This allows us to use non-pointer fields in Opts structs and still have optional parameters
func ApplyDefaults(opts any) error {
	if opts == nil {
		return nil
	}

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
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
			return err
		}
	}

	return nil
}

// ApplyDefaultsAndExtractParams applies default values and extracts API parameters in one step
// Returns a map of API parameters, excluding non-API fields like OnLimitExceeded
func ApplyDefaultsAndExtractParams[T any](opts *T) (map[string]string, error) {
	// If opts is nil, return empty params map
	if opts == nil {
		opts = new(T)
	}

	// First apply defaults
	if err := ApplyDefaults(opts); err != nil {
		return nil, err
	}

	// Then extract API parameters
	return ExtractAPIParams(opts)
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
		// For custom types, try to find a string representation
		if field.Type().Implements(reflect.TypeOf((*interface{ String() string })(nil)).Elem()) {
			// This is a custom type that implements String() method
			// We'll need to handle this case by case
			return nil
		}
	}

	return nil
}

// ExtractAPIParams extracts API parameters from opts struct, excluding non-API fields
// Returns a map of parameter names to values, only including fields that should be sent to the API
func ExtractAPIParams(opts any) (map[string]string, error) {
	if opts == nil {
		return make(map[string]string), nil
	}

	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("opts must be a struct or pointer to struct")
	}

	params := make(map[string]string)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get the json tag to determine the API parameter name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Extract the parameter name from json tag
		paramName := strings.Split(jsonTag, ",")[0]
		if paramName == "" {
			continue
		}

		// Skip OnLimitExceeded field as it's not an API parameter
		if paramName == "on_limit_exceeded" {
			continue
		}

		// Convert field value to string
		var value string
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(field.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		default:
			// For custom types, try to convert to string
			if field.Type().Implements(reflect.TypeOf((*interface{ String() string })(nil)).Elem()) {
				value = field.Interface().(interface{ String() string }).String()
			} else {
				value = fmt.Sprintf("%v", field.Interface())
			}
		}

		// Only include non-zero values (except for string fields where empty string is valid)
		if value != "" || field.Kind() == reflect.String {
			params[paramName] = value
		}
	}

	return params, nil
}

// ParseDefaultTag parses a default tag value that might contain multiple options
// Format: "value" or "value|description"
func ParseDefaultTag(tagValue string) (value, description string) {
	parts := strings.SplitN(tagValue, "|", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
