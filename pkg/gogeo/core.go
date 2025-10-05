package gogeo

import (
	"encoding/json"
	"fmt"
	"os"
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

// writeGeoParquet writes features to a GeoParquet file
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

	// Create writer with options
	writerOpts := []parquet.WriterOption{
		parquet.KeyValueMetadata(GeoParquetMetadataKey, string(geoMetaJSON)),
		parquet.Compression(&parquet.Zstd),
	}

	// Convert features to records
	records := make([]GeoParquetRecord, 0, len(fc.Features))

	for _, feature := range fc.Features {
		record := GeoParquetRecord{}

		// Add geometry as WKB
		if feature.Geometry != nil {
			wkbBytes, err := wkb.Marshal(feature.Geometry)
			if err != nil {
				return fmt.Errorf("failed to encode geometry as WKB: %w", err)
			}
			record.Geometry = wkbBytes
		}

		// Add the name property if it exists (for this simple example)
		if feature.Properties != nil {
			if name, exists := feature.Properties["name"]; exists && name != nil {
				if nameStr, ok := name.(string); ok {
					record.Name = &nameStr
				}
			}
		}

		records = append(records, record)
	}

	// Create writer and write records
	writer := parquet.NewGenericWriter[GeoParquetRecord](file, writerOpts...)
	defer writer.Close()

	if _, err := writer.Write(records); err != nil {
		return fmt.Errorf("failed to write records: %w", err)
	}

	return nil
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

// featureToParquetRow converts a GeoJSON feature to a Parquet row
func featureToParquetRow(feature *geojson.Feature, propertyInfos []PropertyInfo) (parquet.Row, error) {
	// Create value slice: geometry + properties
	values := make([]parquet.Value, len(propertyInfos)+1)

	// Convert geometry to WKB
	if feature.Geometry != nil {
		wkbBytes, err := wkb.Marshal(feature.Geometry)
		if err != nil {
			return nil, fmt.Errorf("failed to encode geometry as WKB: %w", err)
		}
		values[0] = parquet.ByteArrayValue(wkbBytes)
	} else {
		values[0] = parquet.NullValue()
	}

	// Add property values
	for i, info := range propertyInfos {
		var value parquet.Value

		if feature.Properties == nil {
			value = parquet.NullValue()
		} else if propValue, exists := feature.Properties[info.Name]; exists && propValue != nil {
			value = convertToParquetValue(propValue, info.Type)
		} else {
			value = parquet.NullValue()
		}

		values[i+1] = value
	}

	return values, nil
}

// convertToParquetValue converts a property value to a Parquet value
func convertToParquetValue(value any, expectedType PropertyType) parquet.Value {
	if value == nil {
		return parquet.NullValue()
	}

	switch expectedType {
	case PropertyTypeInt:
		switch v := value.(type) {
		case int:
			return parquet.Int64Value(int64(v))
		case int32:
			return parquet.Int64Value(int64(v))
		case int64:
			return parquet.Int64Value(v)
		case float64:
			return parquet.Int64Value(int64(v))
		default:
			return parquet.NullValue()
		}

	case PropertyTypeFloat:
		switch v := value.(type) {
		case float32:
			return parquet.DoubleValue(float64(v))
		case float64:
			return parquet.DoubleValue(v)
		case int:
			return parquet.DoubleValue(float64(v))
		case int64:
			return parquet.DoubleValue(float64(v))
		default:
			return parquet.NullValue()
		}

	case PropertyTypeBool:
		if v, ok := value.(bool); ok {
			return parquet.BooleanValue(v)
		}
		return parquet.NullValue()

	case PropertyTypeString:
		switch v := value.(type) {
		case string:
			return parquet.ByteArrayValue([]byte(v))
		case []any, map[string]any:
			// Convert complex types to JSON strings
			data, _ := json.Marshal(v)
			return parquet.ByteArrayValue(data)
		default:
			return parquet.ByteArrayValue([]byte(fmt.Sprintf("%v", v)))
		}

	default:
		// Fallback to string representation
		return parquet.ByteArrayValue([]byte(fmt.Sprintf("%v", value)))
	}
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

	// Determine geometry type string
	geomTypeStr := "Unknown"
	if len(typesList) == 1 {
		geomTypeStr = typesList[0]
	} else if len(typesList) > 1 {
		geomTypeStr = "Mixed"
	}

	// Create geometry column metadata
	geomColumn := GeoParquetColumn{
		Name:         DefaultGeometryColumn,
		GeometryType: geomTypeStr,
		CRS:          DefaultCRS,
	}

	// Create property metadata
	properties := make([]GeoParquetProperty, len(propertyInfos))
	for i, info := range propertyInfos {
		properties[i] = GeoParquetProperty{
			Name:     info.Name,
			Type:     info.Type.String(),
			Nullable: info.Nullable,
		}
	}

	// Create GeoParquet metadata
	metadata := &GeoParquet{
		Version:       GeoParquetVersion,
		PrimaryColumn: DefaultGeometryColumn,
		Columns:       []GeoParquetColumn{geomColumn},
		Properties:    properties,
	}

	return metadata
}
