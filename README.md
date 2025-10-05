# gogeo

[![Version](https://img.shields.io/badge/version-v0.2.3-blue)](https://github.com/beyondcivic/gogeo/releases/tag/v0.2.3)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
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
- ✅ **GeoParquet Conversion**: Efficient columnar format output
- ✅ **Automatic Type Inference**: Smart detection of data types from GeoJSON properties
- ✅ **Geometry Support**: Complete support for all GeoJSON geometry types
- ✅ **Feature Collections**: Handle complex multi-feature datasets
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Getting Started

### Prerequisites

- Go 1.24 or later
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

Convert a GeoJSON file to efficient GeoParquet format with automatic type inference.

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
- **Type Safety**: Strongly typed columns with automatic inference
- **Compression**: Built-in compression for reduced file sizes
- **Interoperability**: Wide support across geospatial tools and libraries

### Supported GeoJSON Properties

The parser supports all standard GeoJSON properties and geometry types:

| GeoJSON Element      | GeoParquet Representation       | Description                            |
| -------------------- | ------------------------------- | -------------------------------------- |
| `Point`              | Point geometry column           | Single coordinate point                |
| `LineString`         | LineString geometry column      | Connected line segments                |
| `Polygon`            | Polygon geometry column         | Closed area with optional holes        |
| `MultiPoint`         | MultiPoint geometry column      | Collection of points                   |
| `MultiLineString`    | MultiLineString geometry column | Collection of line strings             |
| `MultiPolygon`       | MultiPolygon geometry column    | Collection of polygons                 |
| `GeometryCollection` | GeometryCollection column       | Mixed geometry types                   |
| `properties`         | Typed property columns          | Feature attributes with inferred types |

## Examples

### Example 1: Basic GeoJSON Conversion

```bash
# Convert a simple GeoJSON file
$ gogeo generate locations.geojson -o locations.geoparquet

Generating GeoParquet file for 'locations.geojson'...
✓ GeoParquet file generated successfully and saved to: locations.geoparquet
```

### Example 2: Processing Feature Collections

Given a GeoJSON file with multiple features, the tool will create a GeoParquet file with all features properly typed and stored in columnar format for efficient querying.

## API Reference

### Core Functions

#### `Generate(geojsonPath, outputPath string) (*geojson.FeatureCollection, error)`

Converts a GeoJSON file to GeoParquet format.

**Parameters:**

- `geojsonPath`: Path to the input .geojson file
- `outputPath`: Path for the output .geoparquet file

**Returns:**

- `*geojson.FeatureCollection`: Parsed feature collection structure
- `error`: Any error that occurred during processing

#### `ValidateOutputPath(outputPath string) error`

Validates the output path for GeoParquet file generation.

#### `IsGeoJsonFile(filename string) bool`

Checks if a file is a valid GeoJSON file.

### Data Structures

## Architecture

The library is organized into several key components:

### Core Package (`pkg/gogeo`)

- **Parsing**: GeoJSON file parsing and validation
- **Conversion**: GeoJSON to GeoParquet format conversion
- **Type Inference**: Automatic detection of data types from GeoJSON properties
- **Utilities**: Helper functions for file handling and validation

### Command Line Interface (`cmd/gogeo`)

- **Cobra-based CLI** with subcommands for each major function
- **Comprehensive help system** with detailed usage examples
- **Flexible output options** and error handling
- **Environment variable support** for configuration

### Contributing

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
