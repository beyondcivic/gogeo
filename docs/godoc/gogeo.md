# gogeo

```go
import "github.com/beyondcivic/gogeo/pkg/gogeo"
```

## Index

- [func Generate\(geojsonPath string, outputPath string\) \(\*geojson.FeatureCollection, error\)](<#Generate>)
- [func IsGeoJsonFile\(filePath string\) bool](<#IsGeoJsonFile>)
- [func ValidateOutputPath\(outputPath string\) error](<#ValidateOutputPath>)
- [type AppError](<#AppError>)
  - [func \(e AppError\) Error\(\) string](<#AppError.Error>)
- [type GeoParquet](<#GeoParquet>)
- [type GeoParquetColumn](<#GeoParquetColumn>)
- [type GeoParquetProperty](<#GeoParquetProperty>)


<a name="Generate"></a>
## func [Generate](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/core.go#L10>)

```go
func Generate(geojsonPath string, outputPath string) (*geojson.FeatureCollection, error)
```

Generate generates Geo Parquet file from a geojson file with automatic type inference.

<a name="IsGeoJsonFile"></a>
## func [IsGeoJsonFile](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/utils.go#L11>)

```go
func IsGeoJsonFile(filePath string) bool
```

IsCSVFile checks if a file appears to be a CSV file based on extension

<a name="ValidateOutputPath"></a>
## func [ValidateOutputPath](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/utils.go#L17>)

```go
func ValidateOutputPath(outputPath string) error
```

ValidateOutputPath validates if the given path is a valid file path

<a name="AppError"></a>
## type [AppError](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/error.go#L5-L10>)



```go
type AppError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="AppError.Error"></a>
### func \(AppError\) [Error](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/error.go#L12>)

```go
func (e AppError) Error() string
```



<a name="GeoParquet"></a>
## type [GeoParquet](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L3-L12>)



```go
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
```

<a name="GeoParquetColumn"></a>
## type [GeoParquetColumn](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L14-L21>)



```go
type GeoParquetColumn struct {
    // Name of the geometry column.
    Name string `json:"name" parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
    // Geometry type (e.g., Point, LineString, Polygon, etc.).
    GeometryType string `json:"geometry_type" parquet:"name=geometry_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
    // Coordinate reference system in WKT format.
    CRS string `json:"crs" parquet:"name=crs, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}
```

<a name="GeoParquetProperty"></a>
## type [GeoParquetProperty](<https://github.com:beyondcivic/gogeo/blob/main/pkg/gogeo/structs.go#L23-L30>)



```go
type GeoParquetProperty struct {
    // Name of the property column.
    Name string `json:"name" parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
    // Data type of the property (e.g., INT32, FLOAT, BYTE_ARRAY, etc.).
    Type string `json:"type" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
    // Indicates if the property can have null values.
    Nullable bool `json:"nullable" parquet:"name=nullable, type=BOOLEAN"`
}
```