package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beyondcivic/gogeo/pkg/gogeo"
)

func main() {
	// Example 1: Simple conversion
	fmt.Println("Example 1: Converting GeoJSON to GeoParquet")
	fmt.Println("===========================================")

	fc, err := gogeo.Generate("input.geojson", "output.geoparquet")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("✓ Successfully converted %d features to GeoParquet\n", len(fc.Features))
	fmt.Println()

	// Example 2: Check if file is GeoJSON before processing
	fmt.Println("Example 2: Validating input file")
	fmt.Println("=================================")

	inputFile := "data.geojson"
	if !gogeo.IsGeoJsonFile(inputFile) {
		fmt.Printf("Warning: %s does not appear to be a GeoJSON file\n", inputFile)
	} else {
		fmt.Printf("✓ %s is a valid GeoJSON file\n", inputFile)
	}
	fmt.Println()

	// Example 3: Validate output path before generation
	fmt.Println("Example 3: Validating output path")
	fmt.Println("==================================")

	outputPath := "./output/data.geoparquet"
	if err := gogeo.ValidateOutputPath(outputPath); err != nil {
		fmt.Printf("Error: Invalid output path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Output path %s is valid\n", outputPath)
	fmt.Println()

	// Example 4: Complete workflow with error handling
	fmt.Println("Example 4: Complete workflow")
	fmt.Println("============================")

	inputPath := "cities.geojson"
	outputPath = "cities.geoparquet"

	// Validate input
	if !gogeo.IsGeoJsonFile(inputPath) {
		log.Fatalf("Input file %s is not a GeoJSON file", inputPath)
	}

	// Validate output path
	if err := gogeo.ValidateOutputPath(outputPath); err != nil {
		log.Fatalf("Invalid output path: %v", err)
	}

	// Generate GeoParquet
	featureCollection, err := gogeo.Generate(inputPath, outputPath)
	if err != nil {
		log.Fatalf("Failed to generate GeoParquet: %v", err)
	}

	// Print statistics
	fmt.Printf("✓ Conversion complete!\n")
	fmt.Printf("  Features: %d\n", len(featureCollection.Features))
	fmt.Printf("  Output: %s\n", outputPath)

	if len(featureCollection.Features) > 0 {
		firstFeature := featureCollection.Features[0]
		fmt.Printf("  First feature geometry type: %s\n", firstFeature.Geometry.GeoJSONType())
		if firstFeature.Properties != nil {
			fmt.Printf("  Property keys: %v\n", getPropertyKeys(firstFeature.Properties))
		}
	}
}

// Helper function to get property keys
func getPropertyKeys(props map[string]interface{}) []string {
	keys := make([]string, 0, len(props))
	for key := range props {
		keys = append(keys, key)
	}
	return keys
}
