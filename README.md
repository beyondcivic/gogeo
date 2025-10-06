# gogeo

[![Version](https://img.shields.io/badge/version-v0.3.0-blue)](https://github.com/beyondcivic/gogeo/releases/tag/v0.3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24.7+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/gogeo.svg)](https://pkg.go.dev/github.com/beyondcivic/gogeo)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go implementation for converting GeoJSON format files to GeoParquet format. This library simplifies working with geospatial data by providing both a command-line interface and a Go library for GeoJSON to GeoParquet conversion.

## Overview

GeoParquet is a standardized way to describe geospatial data in a columnar format that provides efficient storage and query performance. This tool streamlines the process of converting GeoJSON data to GeoParquet format by:

- **Parsing standard GeoJSON files** with full specification compliance
- **Converting to efficient GeoParquet format** with optimized columnar storage
- **Automatic type inference** from GeoJSON data structures
- **Preserving all geospatial metadata** including coordinate reference systems and properties
- **Handling complex geometries** and feature collections
- **Providing both CLI and library interfaces** for different use cases

This project provides both a command-line interface and a Go library for working with geospatial data conversion.

## Key Features

- ✅ **GeoJSON Parsing**: Full GeoJSON specification compliant file parsing
- ✅ **GeoParquet Conversion**: Efficient columnar format output with WKB geometry encoding
- ✅ **Basic Property Support**: Handles the `name` property from GeoJSON features
- ✅ **Geometry Support**: Complete support for all GeoJSON geometry types
- ✅ **Feature Collections**: Handle complex multi-feature datasets
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows
- ✅ **GeoParquet 1.1.0**: Compliant with GeoParquet specification v1.1.0

## Getting Started

### Prerequisites

- Go 1.24.7 or later
- Nix 2.25.4 or later (optional but recommended)
- PowerShell v7.5.1 or later (for building)

### Installation

#### Option 1: Install from Source

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/gogeo.git
cd gogeo
```

2. Build the application:

```bash
go build -o gogeo .
```

#### Option 2: Using Nix (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/gogeo.git
cd gogeo
```

2. Prepare the environment using Nix flakes:

```bash
nix develop
```

3. Build the application:

```bash
./build.ps1
```

#### Option 3: Go Install

```bash
go install github.com/beyondcivic/gogeo@latest
```

## Quick Start

### Command Line Interface

The `gogeo` tool provides commands for converting GeoJSON files to GeoParquet:

```bash
# Convert GeoJSON to GeoParquet
gogeo generate data.geojson -o data.geoparquet

# Show version information
gogeo version
```

### Go Library Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/beyondcivic/gogeo/pkg/gogeo"
)

func main() {
	// Convert GeoJSON to GeoParquet
	featureCollection, err := gogeo.Generate("data.geojson", "data.geoparquet")
	if err != nil {
		log.Fatalf("Error converting data: %v", err)
	}

	fmt.Printf("Converted %d features to GeoParquet\n", len(featureCollection.Features))
}
```

## Detailed Command Reference

### `generate` - Convert GeoJSON to GeoParquet

Convert a GeoJSON file to efficient GeoParquet format with WKB geometry encoding.

```bash
gogeo generate [GEOJSON_FILE] [OPTIONS]
```

**Options:**

- `-o, --output`: Output file path (default: `[filename]_parsed.geoparquet`)

**Examples:**

```bash
# Basic conversion
gogeo generate locations.geojson

# With custom output path
gogeo generate locations.geojson -o my-locations.geoparquet
```

**Environment Variables:**

- `GOGEO_OUTPUT_PATH`: Default output path for generated files

### `version` - Show Version Information

Display version, build information, and system details.

```bash
gogeo version
```

## GeoParquet Output Format

The tool converts GeoJSON data to GeoParquet format, which provides:

- **Columnar Storage**: Efficient storage and query performance
- **WKB Geometry Encoding**: Well-Known Binary format for geometry data
- **Compression**: Built-in Zstd compression for reduced file sizes
- **Interoperability**: Wide support across geospatial tools and libraries
- **GeoParquet Metadata**: Embedded geo metadata following GeoParquet 1.1.0 specification

### Current Implementation Scope

The current implementation focuses on core functionality with:

| Feature               | Status         | Notes                                         |
| --------------------- | -------------- | --------------------------------------------- |
| Geometry Conversion   | ✅ Complete    | All GeoJSON geometry types supported          |
| Name Property         | ✅ Complete    | Extracts and stores the `name` property       |
| Additional Properties | ⚠️ Limited     | Type inference implemented but schema limited |
| Complex Schemas       | 🚧 In Progress | Future enhancement planned                    |

### Supported GeoJSON Elements

| GeoJSON Element      | GeoParquet Representation | Description                     |
| -------------------- | ------------------------- | ------------------------------- |
| `Point`              | WKB geometry column       | Single coordinate point         |
| `LineString`         | WKB geometry column       | Connected line segments         |
| `Polygon`            | WKB geometry column       | Closed area with optional holes |
| `MultiPoint`         | WKB geometry column       | Collection of points            |
| `MultiLineString`    | WKB geometry column       | Collection of line strings      |
| `MultiPolygon`       | WKB geometry column       | Collection of polygons          |
| `GeometryCollection` | WKB geometry column       | Mixed geometry types            |
| `properties.name`    | Optional string column    | Feature name attribute          |

## Examples

### Example 1: Basic GeoJSON Conversion

```bash
# Convert a simple GeoJSON file
$ gogeo generate locations.geojson -o locations.geoparquet

Generating GeoParquet file for 'locations.geojson'...
✓ GeoParquet file generated successfully and saved to: locations.geoparquet
```

### Example 2: Processing Feature Collections

Given a GeoJSON file with multiple features, the tool will create a GeoParquet file with:

- All geometries encoded as WKB (Well-Known Binary) in a single geometry column
- Feature names extracted to an optional `name` column
- GeoParquet metadata embedded following v1.1.0 specification
- Zstd compression applied for efficient storage

**Sample GeoJSON:**

```json
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": { "type": "Point", "coordinates": [1.0, 2.0] },
      "properties": { "name": "Location A" }
    }
  ]
}
```

**Resulting GeoParquet schema:**

- `geometry`: BYTE_ARRAY (WKB-encoded geometry)
- `name`: BYTE_ARRAY OPTIONAL (feature name)

## Current Limitations & Roadmap

### Current Limitations

- **Property Support**: Currently only extracts the `name` property from GeoJSON features
- **Schema Flexibility**: Uses a fixed schema structure (`GeoParquetRecord`)
- **Complex Properties**: Nested objects and arrays are not yet supported

### Planned Enhancements

- 🔄 **Dynamic Schema Generation**: Support for arbitrary GeoJSON properties
- 🔄 **Advanced Type Inference**: Better handling of mixed-type properties
- 🔄 **Complex Property Support**: Nested objects and array properties
- 🔄 **CRS Support**: Coordinate reference system handling beyond EPSG:4326
- 🔄 **Performance Optimizations**: Streaming processing for large files

## Examples

## API Reference

### Core Functions

#### `Generate(geojsonPath, outputPath string) (*geojson.FeatureCollection, error)`

Converts a GeoJSON file to GeoParquet format with WKB geometry encoding.

**Parameters:**

- `geojsonPath`: Path to the input .geojson file
- `outputPath`: Path for the output .geoparquet file

**Returns:**

- `*geojson.FeatureCollection`: Parsed feature collection structure
- `error`: Any error that occurred during processing

#### `ValidateOutputPath(outputPath string) error`

Validates the output path for GeoParquet file generation.

#### `IsGeoJsonFile(filename string) bool`

Checks if a file is a valid GeoJSON file based on file extension.

### Data Structures

#### `GeoParquetRecord`

Represents a single record in the output GeoParquet file:

```go
type GeoParquetRecord struct {
    Geometry []byte  `parquet:"geometry"`        // WKB-encoded geometry
    Name     *string `parquet:"name,optional"`   // Optional name property
}
```

#### `GeoParquet`

GeoParquet metadata structure following v1.1.0 specification:

```go
type GeoParquet struct {
    Version       string                           `json:"version"`
    PrimaryColumn string                           `json:"primary_column"`
    Columns       map[string]GeoParquetColumn      `json:"columns"`
}
```

## GeoParquet Compliance Validation

Files generated by `gogeo` are fully compliant with the [GeoParquet specification v1.1.0](https://geoparquet.org/releases/v1.1.0/). This can be verified using validation tools like `gpq`:

```bash
$ gpq validate ./test_simple.geoparquet

Summary: Passed 20 checks.

 ✓ file must include a "geo" metadata key
 ✓ metadata must be a JSON object
 ✓ metadata must include a "version" string
 ✓ metadata must include a "primary_column" string
 ✓ metadata must include a "columns" object
 ✓ column metadata must include the "primary_column" name
 ✓ column metadata must include a valid "encoding" string
 ✓ column metadata must include a "geometry_types" list
 ✓ optional "crs" must be null or a PROJJSON object
 ✓ optional "orientation" must be a valid string
 ✓ optional "edges" must be a valid string
 ✓ optional "bbox" must be an array of 4 or 6 numbers
 ✓ optional "epoch" must be a number
 ✓ geometry columns must not be grouped
 ✓ geometry columns must be stored using the BYTE_ARRAY parquet type
 ✓ geometry columns must be required or optional, not repeated
 ✓ all geometry values match the "encoding" metadata
 ✓ all geometry types must be included in the "geometry_types" metadata (if not empty)
 ✓ all polygon geometries must follow the "orientation" metadata (if present)
 ✓ all geometries must fall within the "bbox" metadata (if present)
```

This validation confirms that `gogeo` correctly implements:

- Proper GeoParquet metadata structure
- Compliant geometry encoding (WKB)
- Valid column definitions and types
- Specification-adherent file structure

## Architecture

The library is organized into several key components:

### Core Package (`pkg/gogeo`)

- **Parsing**: GeoJSON file parsing using `github.com/paulmach/orb/geojson`
- **Conversion**: GeoJSON to GeoParquet format conversion with WKB encoding
- **Type Inference**: Property type detection (string, int, float, bool, null)
- **Metadata Generation**: GeoParquet 1.1.0 compliant metadata creation
- **Utilities**: File validation and path handling

### Command Line Interface (`cmd/gogeo`)

- **Cobra-based CLI** with subcommands for each major function
- **Comprehensive help system** with detailed usage examples
- **Flexible output options** and error handling
- **Environment variable support** for configuration

### Dependencies

Key external libraries used:

- `github.com/parquet-go/parquet-go`: Parquet file format handling
- `github.com/paulmach/orb`: Geospatial geometry processing and WKB encoding
- `github.com/spf13/cobra`: Command-line interface framework
- `github.com/spf13/viper`: Configuration management

## Technical Implementation

### GeoParquet Specification Compliance

The library implements GeoParquet specification v1.1.0:

- **Metadata Key**: Uses `geo` metadata key as specified
- **Geometry Encoding**: WKB (Well-Known Binary) encoding for all geometries
- **Primary Column**: Default geometry column named `geometry`
- **Schema Validation**: Ensures GeoParquet-compliant file structure

### File Processing Pipeline

1. **GeoJSON Parsing**: Uses `orb/geojson` for standards-compliant parsing
2. **Geometry Conversion**: Converts geometries to WKB using `orb/encoding/wkb`
3. **Property Extraction**: Currently extracts `name` property with optional handling
4. **Metadata Creation**: Generates GeoParquet metadata with geometry type analysis
5. **Parquet Writing**: Uses `parquet-go` with Zstd compression

### Error Handling

The library uses a custom `AppError` type for structured error reporting:

```go
type AppError struct {
    Message string  // User-friendly error message
    Value   any     // Underlying error or additional context
}
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `go test ./...`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature/new-feature`
7. Submit a pull request

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Build Environment

### Using Nix (Recommended)

Use Nix flakes to set up the build environment:

```bash
nix develop
```

### Manual Build

Check the build arguments in `build.ps1`:

```bash
# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Then run:

```bash
./build.ps1
```

Or build manually:

```bash
go build -o gogeo .
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
