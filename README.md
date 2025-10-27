# go-hangar

CLI tool and Go library for interacting with the PaperMC Hangar API.

## Features

- Get detailed information about plugins
- List plugins with filtering and pagination
- Retrieve download URLs for specific versions
- Support for both CLI usage and library import
- Structured logging with slog
- Context-based cancellation
- Multiple output formats (table, JSON)

## Installation

```bash
go install github.com/lexfrei/go-hangar/cmd/hangar@latest
```

## Usage

### CLI

#### Get Plugin Information

```bash
hangar project get <slug>
```

Example:
```bash
hangar project list --limit 5
```

#### List Plugins

```bash
hangar project list [flags]
```

Flags:
- `--limit` - Maximum number of results (default: 25)
- `--offset` - Offset for pagination (default: 0)
- `--category` - Filter by category

Example:
```bash
hangar project list --category gameplay --limit 10
```

#### Get Download URL

```bash
hangar version download-url <slug> <version>
```

Example:
```bash
hangar version download-url essentialsx 2.20.1
```

### Global Flags

- `--base-url` - Hangar API base URL (default: https://hangar.papermc.io/api/v1)
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

	// Get project information
	ctx := context.Background()
	project, err := client.GetProject(ctx, "fancyglow")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plugin: %s\n", project.Name)
	fmt.Printf("Downloads: %d\n", project.Stats.Downloads)
	fmt.Printf("Category: %s\n", project.Category)

	// List projects
	list, err := client.ListProjects(ctx, hangar.ListOptions{
		Limit:    10,
		Offset:   0,
		Category: "gameplay",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total projects: %d\n", list.Pagination.Count)
	for _, p := range list.Result {
		fmt.Printf("- %s (%s)\n", p.Name, p.Namespace.Slug)
	}

	// Get download URL
	downloadURL, err := client.GetDownloadURL(ctx, "fancyglow", "1.0.0")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Download URL: %s\n", downloadURL)
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

```
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
