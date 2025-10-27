package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Commands for working with versions",
	Long:  "Commands for retrieving version information and download URLs.",
}

var versionDownloadURLCmd = &cobra.Command{
	Use:   "download-url <slug> <version>",
	Short: "Get download URL for a specific version",
	Long:  "Retrieve the download URL for a specific version of a project.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]
		version := args[1]

		client := createClient()
		downloadURL, err := client.GetDownloadURL(ctx, slug, version)
		if err != nil {
			return errors.Wrap(err, "failed to get download URL")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			result := map[string]string{
				"slug":        slug,
				"version":     version,
				"downloadUrl": downloadURL,
			}
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(result); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		default:
			// For table and other formats, just print the URL
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), downloadURL)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(versionDownloadURLCmd)
}
