# go-hangar

CLI tool and Go library for interacting with the PaperMC Hangar API.

## Features

- **Projects**: Get information, list projects, retrieve versions and download URLs
- **Users**: Search users, view profiles, starred/watching/pinned projects
- **Authors & Staff**: List authors and Hangar staff members
- **Project Social**: View project members, stargazers, and watchers
- **Statistics**: Retrieve daily project and version statistics with date filtering
- **Pages**: Access project documentation and README pages
- **Version Utilities**: Find versions by ID or file hash, get latest releases
- **Full CLI & Library**: Support for both CLI usage and library import
- **Structured Logging**: Built-in slog integration
- **Context Cancellation**: Graceful shutdown support
- **Multiple Formats**: Table and JSON output modes
- **API Coverage**: 27/40 endpoints (67.5% - all read operations)

> 📋 See [ROADMAP.md](docs/ROADMAP.md) for detailed implementation status and future plans

## Installation

```bash
go install github.com/lexfrei/go-hangar/cmd/hangar@latest
```

## Usage

### CLI

#### Projects

Get project information:

```bash
hangar project get <slug>
hangar project get fancyglow
```

List projects with filtering:

```bash
hangar project list --category gameplay --limit 10
```

View project members, stargazers, and watchers:

```bash
hangar project members <slug>
hangar project stargazers <slug> --limit 50
hangar project watchers <slug>
```

Get project pages:

```bash
hangar project page <slug> [path]      # Get specific page
hangar project readme <slug>           # Get main README
```

Project statistics:

```bash
hangar project stats <slug> --from 2024-01-01 --to 2024-01-31
```

#### Versions

Get download URL:

```bash
hangar version download-url <slug> <version>
hangar version download-url essentialsx 2.20.1
```

Version utilities:

```bash
hangar version get-by-id 12345
hangar version find-by-hash abc123def456
hangar version latest <slug> --channel Release --platform PAPER
```

Version statistics:

```bash
hangar version stats <slug> <version> --from 2024-01-01 --to 2024-01-31
```

#### Users

Get user information:

```bash
hangar user get <username>
hangar user list [query] --limit 25
```

View user's projects:

```bash
hangar user starred <username>
hangar user watching <username>
hangar user pinned <username>
```

#### Authors & Staff

List authors and staff:

```bash
hangar authors list --limit 50
hangar staff list
```

### Global Flags

- `--base-url` - Hangar API base URL (default: <https://hangar.papermc.io/api/v1>)
- `--token` - Hangar API token for authenticated requests
- `--timeout` - HTTP client timeout (default: 30s)
- `--output` / `-o` - Output format: table, json (default: table)
- `--config` - Config file path (default: $HOME/.config/hangar/config.yaml)

### Library Usage

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/lexfrei/go-hangar/pkg/hangar"
)

func main() {
 // Create client
 client := hangar.NewClient(hangar.Config{
  BaseURL: hangar.DefaultBaseURL,
  Timeout: 30 * time.Second,
 })
 ctx := context.Background()

 // === Projects ===

 // Get project information
 project, err := client.GetProject(ctx, "fancyglow")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Plugin: %s (Downloads: %d)\n", project.Name, project.Stats.Downloads)

 // List projects with filtering
 list, err := client.ListProjects(ctx, hangar.ListOptions{
  Limit:    10,
  Category: "gameplay",
 })
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Total projects: %d\n", list.Pagination.Count)

 // Get project page
 page, err := client.GetProjectMainPage(ctx, "fancyglow")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("README: %s\n", page.Contents)

 // === Versions ===

 // Get download URL
 downloadURL, err := client.GetDownloadURL(ctx, "fancyglow", "2.0.0", "PAPER")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Download: %s\n", downloadURL)

 // Get latest version
 latest, err := client.GetLatestReleaseVersion(ctx, "fancyglow")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Latest: %s\n", latest.Name)

 // === Users ===

 // Get user information
 user, err := client.GetUser(ctx, "username")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("User: %s (Projects: %d)\n", user.Name, user.ProjectCount)

 // Get user's starred projects
 starred, err := client.GetUserStarred(ctx, "username", hangar.ListOptions{Limit: 10})
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Starred: %d projects\n", starred.Pagination.Count)

 // === Statistics ===

 // Get project statistics
 stats, err := client.GetProjectStats(ctx, "fancyglow", "2024-01-01", "2024-01-31")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Stats for %d days\n", len(stats))

 // === Project Social ===

 // Get project members
 members, err := client.GetProjectMembers(ctx, "fancyglow", hangar.ListOptions{Limit: 25})
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Team: %d members\n", members.Pagination.Count)

 // Get stargazers
 stargazers, err := client.GetProjectStargazers(ctx, "fancyglow", hangar.ListOptions{Limit: 25})
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("Stars: %d users\n", stargazers.Pagination.Count)
}
```

## Configuration

Configuration can be provided via:

1. Command-line flags (highest priority)
2. Environment variables (prefix: `HANGAR_`)
3. Config file (`~/.config/hangar/config.yaml`)
4. Defaults (lowest priority)

### Environment Variables

- `HANGAR_API_TOKEN` - API authentication token
- `HANGAR_API_BASE_URL` - Base URL for Hangar API
- `HANGAR_CONFIG` - Path to config file
- `HANGAR_OUTPUT_FORMAT` - Output format (table, json)
- `HANGAR_TIMEOUT` - API request timeout in seconds
- `HANGAR_LOG_LEVEL` - Logging level (debug, info, warn, error)

### Config File Example

```yaml
base_url: https://hangar.papermc.io/api/v1
api_token: your_token_here
timeout: 30s
output: table
```

## Development

### Requirements

- Go 1.25 or later
- golangci-lint for linting

### Build

```bash
go build ./cmd/hangar
```

### Test

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...
```

### Lint

```bash
golangci-lint run
```

## Architecture

Project follows standard Go project layout:

```text
go-hangar/
├── cmd/hangar/          # CLI application entry point
├── internal/            # Private application code
│   ├── cli/            # CLI command implementations
│   ├── client/         # Internal HTTP client
│   └── config/         # Configuration management
├── pkg/hangar/         # Public library API (importable)
│   ├── client.go       # Main client interface
│   └── types.go        # Public types and models
└── testdata/           # Test fixtures
```

## License

See LICENSE file for details.

## Maintainer

Aleksei Sviridkin <f@lex.la>

GPG: F57F 85FC 7975 F22B BC3F 2504 9C17 3EB1 B531 AA1F
