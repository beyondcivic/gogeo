package gogeo

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/parquet-go/parquet-go"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
)

const (
	GeoParquetVersion       = "1.1.0"
	DefaultGeometryColumn   = "geometry"
	DefaultGeometryEncoding = "WKB"
	GeoParquetMetadataKey   = "geo"
	DefaultCRS              = "EPSG:4326"
)

// Generate generates Geo Parquet file from a geojson file with automatic type inference.
func Generate(geojsonPath string, outputPath string) (*geojson.FeatureCollection, error) {
	// Read and parse GeoJSON file
	fc, err := readGeoJSON(geojsonPath)
	if err != nil {
		return nil, AppError{Message: "failed to read GeoJSON file", Value: err}
	}

	if len(fc.Features) == 0 {
		return nil, AppError{Message: "no features found in GeoJSON file"}
	}

	// Write GeoParquet file
	if err := writeGeoParquet(outputPath, fc); err != nil {
		return nil, AppError{Message: "failed to write GeoParquet file", Value: err}
	}

	return fc, nil
}

// readGeoJSON reads and parses a GeoJSON file
func readGeoJSON(path string) (*geojson.FeatureCollection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}

	return fc, nil
}

// PropertyInfo holds information about a property column
type PropertyInfo struct {
	Name     string
	Type     PropertyType
	Nullable bool
}

// writeGeoParquet writes features to a GeoParquet file using dynamic schema
func writeGeoParquet(path string, fc *geojson.FeatureCollection) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Analyze properties to build schema
	propertyInfos := analyzeProperties(fc)

	// Create GeoParquet metadata
	geoMeta := createGeoParquetMetadata(fc, propertyInfos)
	geoMetaJSON, err := json.Marshal(geoMeta)
	if err != nil {
		return fmt.Errorf("failed to marshal geo metadata: %w", err)
	}

	// Build dynamic struct type and create schema
	structType := buildDynamicType(propertyInfos)
	schemaInstance := reflect.New(structType).Elem().Interface()
	schema := parquet.SchemaOf(schemaInstance)

	// Create writer with schema and options
	writerOpts := []parquet.WriterOption{
		parquet.KeyValueMetadata(GeoParquetMetadataKey, string(geoMetaJSON)),
		parquet.Compression(&parquet.Zstd),
	}

	writer := parquet.NewWriter(file, writerOpts...)
	defer writer.Close()

	// Convert features to rows
	for _, feature := range fc.Features {
		record := buildRecord(feature, propertyInfos, structType)
		row := schema.Deconstruct(nil, record)
		if _, err := writer.WriteRows([]parquet.Row{row}); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// buildDynamicType creates a dynamic struct type based on property information
func buildDynamicType(propertyInfos []PropertyInfo) reflect.Type {
	fields := make([]reflect.StructField, 0, len(propertyInfos)+1)

	// Add geometry field
	fields = append(fields, reflect.StructField{
		Name: "Geometry",
		Type: reflect.TypeOf([]byte{}),
		Tag:  `parquet:"name=geometry, type=BYTE_ARRAY, repetition=REQUIRED"`,
	})

	// Add property fields
	for _, propInfo := range propertyInfos {
		var fieldType reflect.Type
		var tag string

		switch propInfo.Type {
		case PropertyTypeInt:
			fieldType = reflect.TypeOf((*int64)(nil))
			tag = fmt.Sprintf(`parquet:"name=%s, type=INT64, repetition=OPTIONAL"`, propInfo.Name)
		case PropertyTypeFloat:
			fieldType = reflect.TypeOf((*float64)(nil))
			tag = fmt.Sprintf(`parquet:"name=%s, type=DOUBLE, repetition=OPTIONAL"`, propInfo.Name)
		case PropertyTypeBool:
			fieldType = reflect.TypeOf((*bool)(nil))
			tag = fmt.Sprintf(`parquet:"name=%s, type=BOOLEAN, repetition=OPTIONAL"`, propInfo.Name)
		default: // PropertyTypeString
			fieldType = reflect.TypeOf((*string)(nil))
			tag = fmt.Sprintf(`parquet:"name=%s, type=BYTE_ARRAY, convertedtype=UTF8, repetition=OPTIONAL"`, propInfo.Name)
		}

		// Convert property name to exported field name
		fieldName := exportFieldName(propInfo.Name)
		fields = append(fields, reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(tag),
		})
	}

	return reflect.StructOf(fields)
}

// exportFieldName converts a property name to an exported field name
func exportFieldName(name string) string {
	if len(name) == 0 {
		return "Field"
	}
	// Capitalize first letter
	runes := []rune(name)
	runes[0] = []rune(string(runes[0]))[0] & ^('a' - 'A')
	result := string(runes)

	// Replace invalid characters with underscore
	validRunes := make([]rune, 0, len(result))
	for i, r := range result {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9' && i > 0) || r == '_' {
			validRunes = append(validRunes, r)
		} else {
			validRunes = append(validRunes, '_')
		}
	}
	return string(validRunes)
}

// buildRecord creates a record (struct value) from a GeoJSON feature
func buildRecord(feature *geojson.Feature, propertyInfos []PropertyInfo, structType reflect.Type) any {
	record := reflect.New(structType).Elem()

	// Set geometry field
	var wkbBytes []byte
	if feature.Geometry != nil {
		var err error
		wkbBytes, err = wkb.Marshal(feature.Geometry)
		if err != nil {
			wkbBytes = []byte{}
		}
	} else {
		wkbBytes = []byte{}
	}
	record.Field(0).SetBytes(wkbBytes)

	// Set property fields
	for i, propInfo := range propertyInfos {
		fieldIndex := i + 1 // +1 because geometry is field 0
		field := record.Field(fieldIndex)

		var value interface{}
		if feature.Properties != nil {
			value = feature.Properties[propInfo.Name]
		}

		if value == nil {
			// Field is already nil (zero value for pointer types)
			continue
		}

		// Set value based on type
		switch propInfo.Type {
		case PropertyTypeInt:
			if v, ok := toInt64(value); ok {
				ptr := reflect.New(field.Type().Elem())
				ptr.Elem().SetInt(v)
				field.Set(ptr)
			}
		case PropertyTypeFloat:
			if v, ok := toFloat64(value); ok {
				ptr := reflect.New(field.Type().Elem())
				ptr.Elem().SetFloat(v)
				field.Set(ptr)
			}
		case PropertyTypeBool:
			if v, ok := value.(bool); ok {
				ptr := reflect.New(field.Type().Elem())
				ptr.Elem().SetBool(v)
				field.Set(ptr)
			}
		default: // PropertyTypeString
			var strValue string
			if v, ok := value.(string); ok {
				strValue = v
			} else {
				// Convert other types to JSON string
				if jsonBytes, err := json.Marshal(value); err == nil {
					strValue = string(jsonBytes)
				}
			}
			if strValue != "" {
				ptr := reflect.New(field.Type().Elem())
				ptr.Elem().SetString(strValue)
				field.Set(ptr)
			}
		}
	}

	return record.Interface()
}

// Helper functions to convert values
func toInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	case float32:
		return int64(val), true
	default:
		return 0, false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

// analyzeProperties collects and analyzes all properties from features
func analyzeProperties(fc *geojson.FeatureCollection) []PropertyInfo {
	propertyTypes := make(map[string]PropertyType)
	propertyNames := make(map[string]bool)

	for _, feature := range fc.Features {
		if feature.Properties == nil {
			continue
		}

		for key, value := range feature.Properties {
			propertyNames[key] = true
			inferredType := inferPropertyType(value)

			if existingType, exists := propertyTypes[key]; exists {
				// Handle type conflicts by promoting to string
				if existingType != inferredType && inferredType != PropertyTypeNull {
					if existingType != PropertyTypeNull {
						propertyTypes[key] = PropertyTypeString
					} else {
						propertyTypes[key] = inferredType
					}
				}
			} else {
				propertyTypes[key] = inferredType
			}
		}
	}

	// Convert to sorted slice for consistent ordering
	names := make([]string, 0, len(propertyNames))
	for name := range propertyNames {
		names = append(names, name)
	}
	sort.Strings(names)

	infos := make([]PropertyInfo, len(names))
	for i, name := range names {
		propType := propertyTypes[name]
		if propType == PropertyTypeNull {
			propType = PropertyTypeString
		}
		infos[i] = PropertyInfo{
			Name:     name,
			Type:     propType,
			Nullable: true,
		}
	}

	return infos
}

// createGeoParquetMetadata creates GeoParquet metadata from a feature collection
func createGeoParquetMetadata(fc *geojson.FeatureCollection, propertyInfos []PropertyInfo) *GeoParquet {
	// Collect geometry types and bounds
	geomTypes := make(map[string]bool)
	var bounds *orb.Bound

	for _, feature := range fc.Features {
		if feature.Geometry != nil {
			geomType := feature.Geometry.GeoJSONType()
			geomTypes[geomType] = true

			featureBounds := feature.Geometry.Bound()
			if bounds == nil {
				b := featureBounds
				bounds = &b
			} else {
				*bounds = bounds.Union(featureBounds)
			}
		}
	}

	// Convert geometry types to slice
	var typesList []string
	for gt := range geomTypes {
		typesList = append(typesList, gt)
	}
	sort.Strings(typesList)

	// Create geometry column metadata
	geomColumn := GeoParquetColumn{
		Encoding:      DefaultGeometryEncoding,
		GeometryTypes: typesList,
		CRS:           nil,
	}

	// Create columns map
	columns := make(map[string]GeoParquetColumn)
	columns[DefaultGeometryColumn] = geomColumn

	// Create GeoParquet metadata
	metadata := &GeoParquet{
		Version:       GeoParquetVersion,
		PrimaryColumn: DefaultGeometryColumn,
		Columns:       columns,
	}

	return metadata
}
