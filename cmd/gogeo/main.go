// Package cmd provides the command-line interface for gogeo.
//
// gogeo is a Go implementation for converting GeoJSON to GeoParquet format.
//
// The command-line tool provides functionality to:
//   - Generate GeoParquet from GeoJSON files with WKB geometry encoding
//   - Display version and build information
//
// # Command Reference
//
// Generate parquet with default output path:
//
//	gogeo generate data.geojson
//
// Show version information:
//
//	gogeo version
//
// # Features
//
// Metadata Generation:
//   - WKB geometry encoding for all supported geometry types
//   - Configurable output paths and validation options
//   - Support for environment variable configuration
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beyondcivic/gogeo/pkg/gogeo"
	"github.com/beyondcivic/gogeo/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root cobra command.
// Call Init() once to initialize child commands.
// Global so it can be picked up by docs/doc-gen.go.
// nolint:gochecknoglobals
var RootCmd = &cobra.Command{
	Use:   "gogeo",
	Short: "GeoParquet tools",
	Long: `A Go implementation for working with the GeoParquet and GeoJson format.
GeoParquet is a standardized way to describe geospatial data in a columnar format.`,
	Version: version.Version,
}

// Call Once.
func Init() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("GOGEO")
	viper.AutomaticEnv()

	// Add child commands
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(generateCmd())
}

func Execute() {
	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper functions

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isGeoJsonFile(filename string) bool {
	return gogeo.IsGeoJsonFile(filename)
}

func determineOutputPath(providedPath, csvPath string) string {
	if providedPath != "" {
		return providedPath
	}

	// Check environment variable
	envOutputPath := os.Getenv("GOGEO_OUTPUT_PATH")
	if envOutputPath != "" {
		return envOutputPath
	}

	// Generate default path based on CSV filename
	baseName := strings.TrimSuffix(filepath.Base(csvPath), filepath.Ext(csvPath))
	return baseName + ".parquet"
}
