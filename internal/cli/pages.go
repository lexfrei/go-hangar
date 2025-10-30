package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

var projectPageCmd = &cobra.Command{
	Use:   "page <slug> [path]",
	Short: "Get project page content",
	Long:  "Retrieve the Markdown content of a project page. Defaults to 'home' page if path not specified.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		var pagePath string
		if len(args) > 1 {
			pagePath = args[1]
		}

		client := createClient()
		page, err := client.GetProjectPage(ctx, slug, pagePath)
		if err != nil {
			return errors.Wrap(err, "failed to get project page")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(page); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		default:
			// For table and other formats, print Markdown content
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "# %s (%s)\n\n", page.Name, page.Slug)
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), page.Contents)
		}

		return nil
	},
}

var projectReadmeCmd = &cobra.Command{
	Use:   "readme <slug>",
	Short: "Get project README (main page)",
	Long:  "Retrieve the main README page content of a project.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		client := createClient()
		page, err := client.GetProjectMainPage(ctx, slug)
		if err != nil {
			return errors.Wrap(err, "failed to get project README")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(page); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		default:
			// For table and other formats, print Markdown content
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), page.Contents)
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectPageCmd)
	projectCmd.AddCommand(projectReadmeCmd)
}
