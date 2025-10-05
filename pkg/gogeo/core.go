package gogeo

import (
	"os"

	"github.com/paulmach/orb/geojson"
)

// Generate generates Geo Parquet file from a geojson file with automatic type inference.
func Generate(geojsonPath string, outputPath string) (*geojson.FeatureCollection, error) {
	// Get file information
	_, err := os.Stat(geojsonPath)
	if err != nil {
		return nil, AppError{Message: "failed to get file info", Value: err}
	}

	// Read and parse GeoJsonfile
	file, err := os.Open(geojsonPath)
	if err != nil {
		return nil, AppError{Message: "failed to open GeoJsonfile", Value: err}
	}
	defer file.Close()

	return nil, nil
}
