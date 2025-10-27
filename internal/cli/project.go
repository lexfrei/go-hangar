// Package cli implements the command-line interface for the hangar tool.
package cli

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Commands for working with projects",
	Long:  "Commands for retrieving information about projects/plugins on Hangar.",
}

var projectGetCmd = &cobra.Command{
	Use:   "get <slug>",
	Short: "Get information about a specific project",
	Long:  "Retrieve detailed information about a project by its slug.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		client := createClient()
		project, err := client.GetProject(ctx, slug)
		if err != nil {
			return errors.Wrap(err, "failed to get project")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(project); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Field", "Value"})
			t.AppendRows([]table.Row{
				{"ID", project.ID},
				{"Name", project.Name},
				{"Slug", project.Namespace.Slug},
				{"Owner", project.Namespace.Owner},
				{"Category", project.Category},
				{"Description", project.Description},
				{"Views", project.Stats.Views},
				{"Downloads", project.Stats.Downloads},
				{"Stars", project.Stats.Stars},
				{"Created", project.CreatedAt.Format("2006-01-02")},
				{"Last Updated", project.LastUpdated.Format("2006-01-02")},
			})
			t.Render()
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long:  "Retrieve a paginated list of projects from Hangar.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		category, _ := cmd.Flags().GetString("category")

		client := createClient()
		list, err := client.ListProjects(ctx, hangar.ListOptions{
			Limit:    limit,
			Offset:   offset,
			Category: category,
		})
		if err != nil {
			return errors.Wrap(err, "failed to list projects")
		}

		slog.Info("retrieved projects",
			"count", list.Pagination.Count,
			"limit", list.Pagination.Limit,
			"offset", list.Pagination.Offset)

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(list); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Name", "Slug", "Category", "Downloads", "Views", "Stars"})
			for _, proj := range list.Result {
				t.AppendRow(table.Row{
					proj.Name,
					proj.Namespace.Slug,
					proj.Category,
					proj.Stats.Downloads,
					proj.Stats.Views,
					proj.Stats.Stars,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d projects\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectListCmd)

	// List command flags
	projectListCmd.Flags().Int("limit", 25, "Maximum number of results")
	projectListCmd.Flags().Int("offset", 0, "Offset for pagination")
	projectListCmd.Flags().String("category", "", "Filter by category")
}
