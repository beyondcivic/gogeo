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
	// For EPSG:4326 (WGS84), we can set CRS to null as it's the default
	geomColumn := GeoParquetColumn{
		Encoding:      DefaultGeometryEncoding,
		GeometryTypes: typesList,
		CRS:           nil, // null for WGS84/EPSG:4326
	}

	// Create columns map
	columns := make(map[string]GeoParquetColumn)
	columns[DefaultGeometryColumn] = geomColumn

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
		Columns:       columns,
	}

	return metadata
}
