package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/spf13/cobra"
)

var authorsCmd = &cobra.Command{
	Use:   "authors",
	Short: "Commands for working with authors",
	Long:  "Commands for retrieving information about Hangar authors (users with projects).",
}

var authorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List authors",
	Long:  "Retrieve a paginated list of authors (users who have published projects).",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.ListAuthors(ctx, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to list authors")
		}

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
			t.AppendHeader(table.Row{"Username", "Projects", "Joined", "Roles"})
			for _, author := range list.Result {
				roles := ""
				for i, role := range author.Roles {
					if i > 0 {
						roles += ", "
					}
					roles += role.Name
				}
				t.AppendRow(table.Row{
					author.Name,
					author.ProjectCount,
					author.JoinDate.Format("2006-01-02"),
					roles,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d authors\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var staffCmd = &cobra.Command{
	Use:   "staff",
	Short: "Commands for working with Hangar staff",
	Long:  "Commands for retrieving information about Hangar staff members.",
}

var staffListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Hangar staff members",
	Long:  "Retrieve a list of all Hangar staff members.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		client := createClient()
		staff, err := client.ListStaff(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to list staff")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(staff); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Username", "Roles", "Joined"})
			for _, member := range staff {
				roles := ""
				for i, role := range member.Roles {
					if i > 0 {
						roles += ", "
					}
					roles += role.Name
				}
				t.AppendRow(table.Row{
					member.Name,
					roles,
					member.JoinDate.Format("2006-01-02"),
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d staff members\n", len(staff))
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authorsCmd)
	authorsCmd.AddCommand(authorsListCmd)

	rootCmd.AddCommand(staffCmd)
	staffCmd.AddCommand(staffListCmd)

	// Authors list command flags
	authorsListCmd.Flags().Int("limit", 25, "Maximum number of results")
	authorsListCmd.Flags().Int("offset", 0, "Offset for pagination")
}
