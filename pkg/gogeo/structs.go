package gogeo

// GeoParquet represents the GeoParquet metadata structure
type GeoParquet struct {
	// GeoParquet version.
	Version string `json:"version"`
	// Name of the primary geometry column.
	PrimaryColumn string `json:"primary_column"`
	// Map of all geometry columns (column name -> column metadata).
	Columns map[string]GeoParquetColumn `json:"columns"`
}

// GeoParquetColumn represents metadata for a geometry column
type GeoParquetColumn struct {
	// Encoding type (e.g., WKB).
	Encoding string `json:"encoding"`
	// List of geometry types (e.g., ["Point"], ["LineString"], etc.).
	GeometryTypes []string `json:"geometry_types"`
	// Coordinate reference system (can be null for WGS84/EPSG:4326).
	CRS *string `json:"crs,omitempty"`
}

// GeoParquetProperty represents metadata for a property column (not used in actual schema)
type GeoParquetProperty struct {
	// Name of the property column.
	Name string `json:"name"`
	// Data type of the property (e.g., INT32, FLOAT, BYTE_ARRAY, etc.).
	Type string `json:"type"`
	// Indicates if the property can have null values.
	Nullable bool `json:"nullable"`
}
