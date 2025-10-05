package gogeo

type GeoParquet struct {
	// GeoParquet version.
	Version string `json:"version" parquet:"name=version, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Name of the primary geometry column.
	PrimaryColumn string `json:"primary_column" parquet:"name=primary_column, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// List of all geometry columns.
	Columns []GeoParquetColumn `json:"columns" parquet:"name=columns, type=LIST"`
	// List of all property columns.
	Properties []GeoParquetProperty `json:"properties" parquet:"name=properties, type=LIST"`
}

type GeoParquetColumn struct {
	// Name of the geometry column.
	Name string `json:"name" parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Geometry type (e.g., Point, LineString, Polygon, etc.).
	GeometryType string `json:"geometry_type" parquet:"name=geometry_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Coordinate reference system in WKT format.
	CRS string `json:"crs" parquet:"name=crs, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

type GeoParquetProperty struct {
	// Name of the property column.
	Name string `json:"name" parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Data type of the property (e.g., INT32, FLOAT, BYTE_ARRAY, etc.).
	Type string `json:"type" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// Indicates if the property can have null values.
	Nullable bool `json:"nullable" parquet:"name=nullable, type=BOOLEAN"`
}
