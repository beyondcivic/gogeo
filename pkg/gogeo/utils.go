package gogeo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsCSVFile checks if a file appears to be a CSV file based on extension
func IsGeoJsonFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".geojson" || ext == ".json"
}

// ValidateOutputPath validates if the given path is a valid file path
func ValidateOutputPath(outputPath string) error {
	if outputPath == "" {
		return AppError{Message: "output path cannot be empty"}
	}

	// Check if the directory exists or can be created
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return AppError{
				Message: fmt.Sprintf("cannot create directory %s", dir),
				Value:   err,
			}
		}
	}

	// Check if we can write to the file (create a temporary file to test)
	tempFile := outputPath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return AppError{
			Message: fmt.Sprintf("cannot write to path %s", outputPath),
			Value:   err,
		}
	}
	file.Close()
	return os.Remove(tempFile) // Clean up the temporary file
}
