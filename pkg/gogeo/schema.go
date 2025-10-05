package gogeo

import (
	"reflect"
)

// PropertyType represents the inferred type of a GeoJSON property
type PropertyType int

const (
	PropertyTypeUnknown PropertyType = iota
	PropertyTypeString
	PropertyTypeInt
	PropertyTypeFloat
	PropertyTypeBool
	PropertyTypeNull
)

// inferPropertyType infers the Parquet type from a GeoJSON property value
func inferPropertyType(value any) PropertyType {
	if value == nil {
		return PropertyTypeNull
	}

	switch v := value.(type) {
	case bool:
		return PropertyTypeBool
	case int, int8, int16, int32, int64:
		return PropertyTypeInt
	case uint, uint8, uint16, uint32, uint64:
		return PropertyTypeInt
	case float32, float64:
		return PropertyTypeFloat
	case string:
		return PropertyTypeString
	default:
		// Use reflection for JSON unmarshaled types
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Bool:
			return PropertyTypeBool
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return PropertyTypeInt
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return PropertyTypeInt
		case reflect.Float32, reflect.Float64:
			return PropertyTypeFloat
		case reflect.String:
			return PropertyTypeString
		case reflect.Map, reflect.Slice, reflect.Array:
			// Complex types stored as JSON strings
			return PropertyTypeString
		default:
			return PropertyTypeString
		}
	}
}

// String returns the string representation of a PropertyType
func (pt PropertyType) String() string {
	switch pt {
	case PropertyTypeString:
		return "string"
	case PropertyTypeInt:
		return "int64"
	case PropertyTypeFloat:
		return "double"
	case PropertyTypeBool:
		return "boolean"
	case PropertyTypeNull:
		return "null"
	default:
		return "unknown"
	}
}
