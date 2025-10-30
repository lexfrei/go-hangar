package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/spf13/cobra"
)

var projectMembersCmd = &cobra.Command{
	Use:   "members <slug>",
	Short: "Get project team members",
	Long:  "Retrieve a list of team members for a project.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.GetProjectMembers(ctx, slug, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get project members")
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
			t.AppendHeader(table.Row{"Username", "Roles", "Accepted"})
			for _, member := range list.Result {
				roles := ""
				for i, role := range member.Roles {
					if i > 0 {
						roles += ", "
					}
					roles += role.Name
				}
				accepted := "Yes"
				if !member.Accepted {
					accepted = "No (Pending)"
				}
				t.AppendRow(table.Row{
					member.User,
					roles,
					accepted,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d members\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var projectStargazersCmd = &cobra.Command{
	Use:   "stargazers <slug>",
	Short: "Get users who starred the project",
	Long:  "Retrieve a list of users who have starred a project.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.GetProjectStargazers(ctx, slug, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get project stargazers")
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
			t.AppendHeader(table.Row{"Username", "Projects", "Joined"})
			for _, user := range list.Result {
				t.AppendRow(table.Row{
					user.Name,
					user.ProjectCount,
					user.JoinDate.Format("2006-01-02"),
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d stargazers\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var projectWatchersCmd = &cobra.Command{
	Use:   "watchers <slug>",
	Short: "Get users watching the project",
	Long:  "Retrieve a list of users who are watching a project.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.GetProjectWatchers(ctx, slug, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get project watchers")
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
			t.AppendHeader(table.Row{"Username", "Projects", "Joined"})
			for _, user := range list.Result {
				t.AppendRow(table.Row{
					user.Name,
					user.ProjectCount,
					user.JoinDate.Format("2006-01-02"),
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d watchers\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectMembersCmd)
	projectCmd.AddCommand(projectStargazersCmd)
	projectCmd.AddCommand(projectWatchersCmd)

	// Members command flags
	projectMembersCmd.Flags().Int("limit", 25, "Maximum number of results")
	projectMembersCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Stargazers command flags
	projectStargazersCmd.Flags().Int("limit", 25, "Maximum number of results")
	projectStargazersCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Watchers command flags
	projectWatchersCmd.Flags().Int("limit", 25, "Maximum number of results")
	projectWatchersCmd.Flags().Int("offset", 0, "Offset for pagination")
}
