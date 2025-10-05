// gogeo is a command-line tool and Go library for working with GeoParquet and GeoJSON formats.
//
// GeoParquet is a standardized way to describe geospatial data in a columnar format.
// This tool simplifies the creation of GeoParquet-compatible metadata from GeoJSON data sources.
//
// # Installation
//
// Install the latest version:
//
//	go install github.com/beyondcivic/gogeo@latest
//
// # Usage
//
// Generate GeoParquet from a GeoJSON file:
//
//	gogeo generate data.geojson -o data.geoparquet
//
// For detailed usage information, run:
//
//	gogeo --help
package main

import (
	cmd "github.com/beyondcivic/gogeo/cmd/gogeo"
)

func main() {
	cmd.Init()
	cmd.Execute()
}
