package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var projectStatsCmd = &cobra.Command{
	Use:   "stats <slug>",
	Short: "Get project statistics",
	Long:  "Retrieve daily statistics for a project, optionally filtered by date range.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		fromDate, _ := cmd.Flags().GetString("from")
		toDate, _ := cmd.Flags().GetString("to")

		client := createClient()
		stats, err := client.GetProjectStats(ctx, slug, fromDate, toDate)
		if err != nil {
			return errors.Wrap(err, "failed to get project stats")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(stats); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Date", "Downloads", "Views"})
			for date, dailyStats := range stats {
				t.AppendRow(table.Row{
					date,
					dailyStats.Downloads,
					dailyStats.Views,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal days: %d\n", len(stats))
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var versionStatsCmd = &cobra.Command{
	Use:   "stats <slug> <version>",
	Short: "Get version statistics",
	Long:  "Retrieve daily statistics for a specific version, optionally filtered by date range.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]
		version := args[1]

		fromDate, _ := cmd.Flags().GetString("from")
		toDate, _ := cmd.Flags().GetString("to")

		client := createClient()
		stats, err := client.GetVersionStats(ctx, slug, version, fromDate, toDate)
		if err != nil {
			return errors.Wrap(err, "failed to get version stats")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(stats); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Date", "Downloads", "Views"})
			for date, dailyStats := range stats {
				t.AppendRow(table.Row{
					date,
					dailyStats.Downloads,
					dailyStats.Views,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal days: %d\n", len(stats))
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectStatsCmd)
	versionCmd.AddCommand(versionStatsCmd)

	// Project stats flags
	projectStatsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	projectStatsCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")

	// Version stats flags
	versionStatsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD)")
	versionStatsCmd.Flags().String("to", "", "End date (YYYY-MM-DD)")
}
