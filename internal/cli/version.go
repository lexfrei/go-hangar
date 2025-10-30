package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/jedib0t/go-pretty/v6/table"
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
		versionName := args[1]

		platform, _ := cmd.Flags().GetString("platform")

		client := createClient()

		// First get project to find owner
		project, err := client.GetProject(ctx, slug)
		if err != nil {
			return errors.Wrap(err, "failed to get project")
		}

		downloadURL, err := client.GetDownloadURL(ctx, project.Namespace.Owner, slug, versionName, platform)
		if err != nil {
			return errors.Wrap(err, "failed to get download URL")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			result := map[string]string{
				"owner":       project.Namespace.Owner,
				"slug":        slug,
				"version":     versionName,
				"platform":    platform,
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

var versionGetByIDCmd = &cobra.Command{
	Use:   "get-by-id <id>",
	Short: "Get version by ID",
	Long:  "Retrieve version information by its unique identifier.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		versionID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid version ID")
		}

		client := createClient()
		version, err := client.GetVersionByID(ctx, versionID)
		if err != nil {
			return errors.Wrap(err, "failed to get version")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(version); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Field", "Value"})
			t.AppendRows([]table.Row{
				{"ID", version.ID},
				{"Name", version.Name},
				{"Author", version.Author},
				{"Created", version.CreatedAt.Format("2006-01-02 15:04:05")},
				{"Visibility", version.Visibility},
				{"Review State", version.ReviewState},
				{"Downloads", version.Stats.TotalDownloads},
			})
			t.Render()
			if version.Description != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nDescription:\n%s\n", version.Description)
			}
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var versionFindByHashCmd = &cobra.Command{
	Use:   "find-by-hash <hash>",
	Short: "Find version by file hash",
	Long:  "Find version information by file hash (MD5, SHA-256, or SHA-512).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		hash := args[0]

		client := createClient()
		version, err := client.GetVersionByHash(ctx, hash)
		if err != nil {
			return errors.Wrap(err, "failed to find version by hash")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(version); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Field", "Value"})
			t.AppendRows([]table.Row{
				{"ID", version.ID},
				{"Name", version.Name},
				{"Author", version.Author},
				{"Created", version.CreatedAt.Format("2006-01-02 15:04:05")},
				{"Visibility", version.Visibility},
				{"Review State", version.ReviewState},
				{"Downloads", version.Stats.TotalDownloads},
			})
			t.Render()
			if version.Description != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nDescription:\n%s\n", version.Description)
			}
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

var versionLatestCmd = &cobra.Command{
	Use:   "latest <slug>",
	Short: "Get latest version of a project",
	Long:  "Retrieve the latest version of a project, optionally filtered by channel, platform, and Minecraft version.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		slug := args[0]

		channel, _ := cmd.Flags().GetString("channel")
		platform, _ := cmd.Flags().GetString("platform")
		minecraftVersion, _ := cmd.Flags().GetString("minecraft-version")

		client := createClient()
		version, err := client.GetLatestVersion(ctx, slug, channel, platform, minecraftVersion)
		if err != nil {
			return errors.Wrap(err, "failed to get latest version")
		}

		// Output based on format
		outputFormat := cmd.Flag("output").Value.String()
		switch outputFormat {
		case "json":
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(version); err != nil {
				return errors.Wrap(err, "failed to encode JSON")
			}
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Field", "Value"})
			t.AppendRows([]table.Row{
				{"ID", version.ID},
				{"Name", version.Name},
				{"Author", version.Author},
				{"Created", version.CreatedAt.Format("2006-01-02 15:04:05")},
				{"Visibility", version.Visibility},
				{"Review State", version.ReviewState},
				{"Downloads", version.Stats.TotalDownloads},
			})
			t.Render()
			if version.Description != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nDescription:\n%s\n", version.Description)
			}
		default:
			return errors.Newf("unsupported output format: %s", outputFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(versionDownloadURLCmd)
	versionCmd.AddCommand(versionGetByIDCmd)
	versionCmd.AddCommand(versionFindByHashCmd)
	versionCmd.AddCommand(versionLatestCmd)

	// download-url command flags
	versionDownloadURLCmd.Flags().String("platform", "PAPER", "Platform to download for (PAPER, WATERFALL, VELOCITY)")

	// latest command flags
	versionLatestCmd.Flags().String("channel", "", "Release channel (Release, Snapshot, etc.)")
	versionLatestCmd.Flags().String("platform", "", "Platform filter (PAPER, WATERFALL, VELOCITY)")
	versionLatestCmd.Flags().String("minecraft-version", "", "Minecraft version filter (e.g., 1.20.1)")
}
