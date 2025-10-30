package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Commands for working with users",
	Long:  "Commands for retrieving information about Hangar users.",
}

var userGetCmd = &cobra.Command{
	Use:   "get <username>",
	Short: "Get information about a specific user",
	Long:  "Retrieve detailed information about a user by username.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		username := args[0]

		client := createClient()
		user, err := client.GetUser(ctx, username)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(user); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Field", "Value"})
			t.AppendRows([]table.Row{
				{"Username", user.Name},
				{"Tagline", user.TagLine},
				{"Joined", user.JoinDate.Format("2006-01-02")},
				{"Projects", user.ProjectCount},
				{"Locked", user.Locked},
			})
			if len(user.Roles) > 0 {
				roles := ""
				for i, role := range user.Roles {
					if i > 0 {
						roles += ", "
					}
					roles += role.Name
				}
				t.AppendRow(table.Row{"Roles", roles})
			}
			t.Render()
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list [query]",
	Short: "List or search users",
	Long:  "Retrieve a paginated list of users, optionally filtered by search query.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var query string
		if len(args) > 0 {
			query = args[0]
		}

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.ListUsers(ctx, query, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to list users")
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
			for _, u := range list.Result {
				roles := ""
				for i, role := range u.Roles {
					if i > 0 {
						roles += ", "
					}
					roles += role.Name
				}
				t.AppendRow(table.Row{
					u.Name,
					u.ProjectCount,
					u.JoinDate.Format("2006-01-02"),
					roles,
				})
			}
			t.Render()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nTotal: %d users\n", list.Pagination.Count)
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var userStarredCmd = &cobra.Command{
	Use:   "starred <username>",
	Short: "Get projects starred by a user",
	Long:  "Retrieve a list of projects that a user has starred.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		username := args[0]

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.GetUserStarred(ctx, username, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get starred projects")
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
			t.AppendHeader(table.Row{"Name", "Slug", "Category", "Downloads", "Stars"})
			for _, proj := range list.Result {
				t.AppendRow(table.Row{
					proj.Name,
					proj.Namespace.Slug,
					proj.Category,
					proj.Stats.Downloads,
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

var userWatchingCmd = &cobra.Command{
	Use:   "watching <username>",
	Short: "Get projects watched by a user",
	Long:  "Retrieve a list of projects that a user is watching.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		username := args[0]

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		client := createClient()
		list, err := client.GetUserWatching(ctx, username, hangar.ListOptions{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "failed to get watching projects")
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
			t.AppendHeader(table.Row{"Name", "Slug", "Category", "Downloads", "Watchers"})
			for _, proj := range list.Result {
				t.AppendRow(table.Row{
					proj.Name,
					proj.Namespace.Slug,
					proj.Category,
					proj.Stats.Downloads,
					proj.Stats.Watchers,
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

var userPinnedCmd = &cobra.Command{
	Use:   "pinned <username>",
	Short: "Get projects pinned by a user",
	Long:  "Retrieve a list of projects that a user has pinned to their profile.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		username := args[0]

		client := createClient()
		list, err := client.GetUserPinned(ctx, username)
		if err != nil {
			return errors.Wrap(err, "failed to get pinned projects")
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
			t.AppendHeader(table.Row{"Name", "Slug", "Category", "Downloads", "Stars"})
			for _, proj := range list.Result {
				t.AppendRow(table.Row{
					proj.Name,
					proj.Namespace.Slug,
					proj.Category,
					proj.Stats.Downloads,
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
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userStarredCmd)
	userCmd.AddCommand(userWatchingCmd)
	userCmd.AddCommand(userPinnedCmd)

	// List command flags
	userListCmd.Flags().Int("limit", 25, "Maximum number of results")
	userListCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Starred command flags
	userStarredCmd.Flags().Int("limit", 25, "Maximum number of results")
	userStarredCmd.Flags().Int("offset", 0, "Offset for pagination")

	// Watching command flags
	userWatchingCmd.Flags().Int("limit", 25, "Maximum number of results")
	userWatchingCmd.Flags().Int("offset", 0, "Offset for pagination")
}
