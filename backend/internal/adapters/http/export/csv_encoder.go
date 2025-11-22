package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// EncodeCSV writes a slice of structs as CSV to the provided writer.
// It uses reflection to read csv struct tags for headers and field order.
// Nil/optional pointer fields are written as empty strings.
// String slice fields are joined with ", " separator.
func EncodeCSV(w io.Writer, data interface{}) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Get the slice value
	sliceValue := reflect.ValueOf(data)
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice, got %s", sliceValue.Kind())
	}

	// Handle empty slice
	if sliceValue.Len() == 0 {
		// For empty data, we still need to write headers if we can determine the type
		sliceType := reflect.TypeOf(data)
		if sliceType.Kind() == reflect.Slice {
			elemType := sliceType.Elem()
			headers := extractHeaders(elemType)
			if len(headers) > 0 {
				if err := writer.Write(headers); err != nil {
					return fmt.Errorf("failed to write CSV headers: %w", err)
				}
			}
		}
		return nil
	}

	// Get the element type
	elemType := sliceValue.Index(0).Type()

	// Extract headers from struct tags
	headers := extractHeaders(elemType)
	if len(headers) == 0 {
		return fmt.Errorf("no csv tags found in struct")
	}

	// Write headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data rows
	for i := 0; i < sliceValue.Len(); i++ {
		row, err := structToRow(sliceValue.Index(i), headers, elemType)
		if err != nil {
			return fmt.Errorf("failed to convert struct to row at index %d: %w", i, err)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row at index %d: %w", i, err)
		}
	}

	return nil
}

// extractHeaders extracts CSV headers from struct tags in field order
func extractHeaders(t reflect.Type) []string {
	var headers []string

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return headers
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csvTag := field.Tag.Get("csv")
		if csvTag != "" && csvTag != "-" {
			headers = append(headers, csvTag)
		}
	}

	return headers
}

// structToRow converts a struct value to a CSV row based on headers
func structToRow(v reflect.Value, headers []string, t reflect.Type) ([]string, error) {
	row := make([]string, len(headers))

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// Return empty strings for all fields
			return row, nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}

	// Create a map of csv tag to field value
	fieldMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		csvTag := field.Tag.Get("csv")
		if csvTag == "" || csvTag == "-" {
			continue
		}

		fieldValue := v.Field(i)
		fieldMap[csvTag] = fieldValueToString(fieldValue)
	}

	// Fill row in header order
	for i, header := range headers {
		row[i] = fieldMap[header]
	}

	return row, nil
}

// fieldValueToString converts a field value to its string representation
func fieldValueToString(v reflect.Value) string {
	// Handle nil pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	// Handle string slices - join with ", "
	if v.Kind() == reflect.Slice {
		if v.Type().Elem().Kind() == reflect.String {
			var strs []string
			for i := 0; i < v.Len(); i++ {
				strs = append(strs, v.Index(i).String())
			}
			return strings.Join(strs, ", ")
		}
	}

	// Handle basic types
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
