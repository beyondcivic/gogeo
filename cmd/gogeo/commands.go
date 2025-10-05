// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/gogeo/pkg/gogeo"
	"github.com/beyondcivic/gogeo/pkg/version"
	"github.com/spf13/cobra"
)

// Version Command.
// Displays tool version and build information
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the gogeo tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			stamp := version.RetrieveStamp()
			fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
			fmt.Printf("  Git ref: %s\n", stamp.VCSRevision)
			fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", stamp.InfoGoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		},
	}
}

// Generate command
func generateCmd() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate [geojsonPath]",
		Short: "Generate GeoParquet from a GeoJsonfile",
		Long:  `Generate GeoParquet from a GeoJsonfile, automatically inferring data types.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			geojsonPath := args[0]
			flagOutputPath, _ := cmd.Flags().GetString("output")

			// Validate input file
			if !fileExists(geojsonPath) {
				fmt.Printf("Error: GeoJsonfile '%s' does not exist.\n", geojsonPath)
				os.Exit(1)
			}

			if !isGeoJsonFile(geojsonPath) {
				fmt.Printf("Error: File '%s' does not appear to be a GeoJsonfile.\n", geojsonPath)
				os.Exit(1)
			}

			// Determine output path
			outputPath := determineOutputPath(flagOutputPath, geojsonPath)

			// Validate output path
			if err := gogeo.ValidateOutputPath(outputPath); err != nil {
				fmt.Printf("Error: Invalid output path: %v\n", err)
				os.Exit(1)
			}

			// Generate metadata
			fmt.Printf("Generating GeoParquet file for '%s'...\n", geojsonPath)
			_, err := gogeo.Generate(geojsonPath, outputPath)
			if err != nil {
				fmt.Printf("Error generating metadata: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("✓ GeoParquet file generated successfully")
			if outputPath != "" {
				fmt.Printf(" and saved to: %s\n", outputPath)
			}

		},
	}
	generateCmd.Flags().StringP("output", "o", "", "Output path for the GeoParquet file")

	return generateCmd
}
