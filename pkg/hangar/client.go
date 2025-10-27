package hangar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	// DefaultBaseURL is the default Hangar API base URL.
	DefaultBaseURL = "https://hangar.papermc.io/api/v1"
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
	// DefaultLimit is the default pagination limit.
	DefaultLimit = 25
)

// Client is the Hangar API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// Config contains configuration for the Hangar client.
type Config struct {
	// BaseURL is the API base URL (defaults to DefaultBaseURL).
	BaseURL string
	// Token is the optional API authentication token.
	Token string
	// Timeout is the HTTP client timeout (defaults to DefaultTimeout).
	Timeout time.Duration
	// HTTPClient is an optional custom HTTP client.
	HTTPClient *http.Client
}

// NewClient creates a new Hangar API client.
func NewClient(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = DefaultTimeout
		}
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		baseURL:    cfg.BaseURL,
		token:      cfg.Token,
		httpClient: httpClient,
	}
}

// ListOptions contains options for listing resources.
type ListOptions struct {
	// Limit is the maximum number of items to return (default: 25).
	Limit int
	// Offset is the starting position (default: 0).
	Offset int
	// Category filters projects by category (optional).
	Category string
}

// GetProject retrieves information about a specific project.
func (c *Client) GetProject(ctx context.Context, slug string) (*Project, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s", c.baseURL, url.PathEscape(slug))

	var project Project
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &project); err != nil {
		return nil, errors.Wrap(err, "failed to get project")
	}

	return &project, nil
}

// ListProjects retrieves a paginated list of projects.
func (c *Client) ListProjects(ctx context.Context, opts ListOptions) (*ProjectsList, error) {
	endpoint := fmt.Sprintf("%s/projects", c.baseURL)

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	if opts.Category != "" {
		params.Set("category", opts.Category)
	}

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list ProjectsList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to list projects")
	}

	return &list, nil
}

// ListVersions retrieves a paginated list of versions for a project.
// owner is the project owner username, slug is the project identifier.
func (c *Client) ListVersions(ctx context.Context, owner, slug string, opts ListOptions) (*VersionsList, error) {
	if owner == "" {
		return nil, errors.New("owner cannot be empty")
	}
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/%s/versions",
		c.baseURL, url.PathEscape(owner), url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list VersionsList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to list versions")
	}

	return &list, nil
}

// GetDownloadURL retrieves the download URL for a specific version.
// owner is the project owner username, slug is the project identifier, version is the version name.
// platform specifies which platform to get the download for (e.g., "PAPER", "WATERFALL").
func (c *Client) GetDownloadURL(ctx context.Context, owner, slug, version, platform string) (string, error) {
	if owner == "" {
		return "", errors.New("owner cannot be empty")
	}
	if slug == "" {
		return "", errors.New("slug cannot be empty")
	}
	if version == "" {
		return "", errors.New("version cannot be empty")
	}
	if platform == "" {
		platform = "PAPER" // Default to PAPER platform
	}

	// List versions to find the specific one
	versions, err := c.ListVersions(ctx, owner, slug, ListOptions{Limit: 100})
	if err != nil {
		return "", errors.Wrap(err, "failed to list versions")
	}

	// Find matching version
	for _, v := range versions.Result {
		if v.Name == version {
			if downloadInfo, ok := v.Downloads[platform]; ok {
				// Prefer downloadUrl, fallback to externalUrl
				if downloadInfo.DownloadURL != "" {
					return downloadInfo.DownloadURL, nil
				}
				if downloadInfo.ExternalURL != "" {
					return downloadInfo.ExternalURL, nil
				}
			}
			return "", errors.Newf("no download URL found for platform %s", platform)
		}
	}

	return "", errors.Newf("version %s not found", version)
}

// doRequest performs an HTTP request with proper error handling.
func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "go-hangar/1.0")

	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	slog.DebugContext(ctx, "making API request",
		"method", method,
		"url", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "HTTP request failed")
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.WarnContext(ctx, "failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return errors.Newf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return errors.Wrap(err, "failed to decode response")
		}
	}

	return nil
}
