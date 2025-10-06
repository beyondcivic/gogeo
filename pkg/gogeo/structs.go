package gogeo

// GeoParquetRecord represents a single record in a GeoParquet file
type GeoParquetRecord struct {
	Geometry []byte  `parquet:"name=geometry, type=BYTE_ARRAY"`
	Name     *string `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, repetition=OPTIONAL"`
}

type GeoParquet struct {
	// GeoParquet version.
	Version string `json:"version" parquet:"name=version, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Name of the primary geometry column.
	PrimaryColumn string `json:"primary_column" parquet:"name=primary_column, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Map of all geometry columns (column name -> column metadata).
	Columns map[string]GeoParquetColumn `json:"columns" parquet:"name=columns"`
}

type GeoParquetColumn struct {
	// Encoding type (e.g., WKB).
	Encoding string `json:"encoding" parquet:"name=encoding, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// List of geometry types (e.g., ["Point"], ["LineString"], etc.).
	GeometryTypes []string `json:"geometry_types" parquet:"name=geometry_types, type=LIST"`
	// Coordinate reference system (e.g., "EPSG:4326").
	CRS *string `json:"crs,omitempty" parquet:"name=crs, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY, repetition=OPTIONAL"`
}

type GeoParquetProperty struct {
	// Name of the property column.
	Name string `json:"name" parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Data type of the property (e.g., INT32, FLOAT, BYTE_ARRAY, etc.).
	Type string `json:"type" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Indicates if the property can have null values.
	Nullable bool `json:"nullable" parquet:"name=nullable, type=BOOLEAN"`
}
