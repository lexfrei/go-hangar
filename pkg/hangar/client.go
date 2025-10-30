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

// ListUsers retrieves a paginated list of users matching a query.
func (c *Client) ListUsers(ctx context.Context, query string, opts ListOptions) (*UserList, error) {
	endpoint := fmt.Sprintf("%s/users", c.baseURL)

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	if query != "" {
		params.Set("query", query)
	}

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list UserList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}

	return &list, nil
}

// GetUser retrieves detailed information about a specific user.
func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/users/%s", c.baseURL, url.PathEscape(username))

	var user User
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &user); err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

// GetUserStarred retrieves projects starred by a user.
func (c *Client) GetUserStarred(ctx context.Context, username string, opts ListOptions) (*ProjectsList, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/users/%s/starred", c.baseURL, url.PathEscape(username))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list ProjectsList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get starred projects")
	}

	return &list, nil
}

// GetUserWatching retrieves projects watched by a user.
func (c *Client) GetUserWatching(ctx context.Context, username string, opts ListOptions) (*ProjectsList, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/users/%s/watching", c.baseURL, url.PathEscape(username))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list ProjectsList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get watching projects")
	}

	return &list, nil
}

// GetUserPinned retrieves projects pinned by a user.
func (c *Client) GetUserPinned(ctx context.Context, username string) (*ProjectsList, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/users/%s/pinned", c.baseURL, url.PathEscape(username))

	var list ProjectsList
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get pinned projects")
	}

	return &list, nil
}

// ListAuthors retrieves a paginated list of authors (users with projects).
func (c *Client) ListAuthors(ctx context.Context, opts ListOptions) (*AuthorList, error) {
	endpoint := fmt.Sprintf("%s/authors", c.baseURL)

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list AuthorList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to list authors")
	}

	return &list, nil
}

// ListStaff retrieves the list of Hangar staff members.
func (c *Client) ListStaff(ctx context.Context) ([]StaffMember, error) {
	endpoint := fmt.Sprintf("%s/staff", c.baseURL)

	var staff []StaffMember
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &staff); err != nil {
		return nil, errors.Wrap(err, "failed to list staff")
	}

	return staff, nil
}

// GetVersionByID retrieves a version by its unique ID.
func (c *Client) GetVersionByID(ctx context.Context, versionID int64) (*Version, error) {
	if versionID <= 0 {
		return nil, errors.New("versionID must be positive")
	}

	endpoint := fmt.Sprintf("%s/versions/%d", c.baseURL, versionID)

	var version Version
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &version); err != nil {
		return nil, errors.Wrap(err, "failed to get version")
	}

	return &version, nil
}

// GetVersionByHash retrieves a version by its file hash.
func (c *Client) GetVersionByHash(ctx context.Context, hash string) (*Version, error) {
	if hash == "" {
		return nil, errors.New("hash cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/versions/find/%s", c.baseURL, url.PathEscape(hash))

	var version Version
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &version); err != nil {
		return nil, errors.Wrap(err, "failed to find version by hash")
	}

	return &version, nil
}

// GetProjectMembers retrieves the list of project team members.
func (c *Client) GetProjectMembers(ctx context.Context, slug string, opts ListOptions) (*MemberList, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/members", c.baseURL, url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list MemberList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get project members")
	}

	return &list, nil
}

// GetProjectStargazers retrieves users who starred the project.
func (c *Client) GetProjectStargazers(ctx context.Context, slug string, opts ListOptions) (*UserList, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/stargazers", c.baseURL, url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list UserList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get project stargazers")
	}

	return &list, nil
}

// GetProjectWatchers retrieves users watching the project.
func (c *Client) GetProjectWatchers(ctx context.Context, slug string, opts ListOptions) (*UserList, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/watchers", c.baseURL, url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(opts.Offset))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	var list UserList
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &list); err != nil {
		return nil, errors.Wrap(err, "failed to get project watchers")
	}

	return &list, nil
}

// GetProjectStats retrieves daily statistics for a project within a date range.
// from and to should be in YYYY-MM-DD format. If empty, returns all available data.
func (c *Client) GetProjectStats(ctx context.Context, slug, fromDate, toDate string) (ProjectStats, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/stats", c.baseURL, url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	if fromDate != "" {
		params.Set("fromDate", fromDate)
	}
	if toDate != "" {
		params.Set("toDate", toDate)
	}

	var fullURL string
	if len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	} else {
		fullURL = endpoint
	}

	var stats ProjectStats
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &stats); err != nil {
		return nil, errors.Wrap(err, "failed to get project stats")
	}

	return stats, nil
}

// GetVersionStats retrieves daily statistics for a specific version within a date range.
// from and to should be in YYYY-MM-DD format. If empty, returns all available data.
func (c *Client) GetVersionStats(ctx context.Context, slug, version, fromDate, toDate string) (VersionStatsData, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}
	if version == "" {
		return nil, errors.New("version cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/versions/%s/stats",
		c.baseURL, url.PathEscape(slug), url.PathEscape(version))

	// Build query parameters
	params := url.Values{}
	if fromDate != "" {
		params.Set("fromDate", fromDate)
	}
	if toDate != "" {
		params.Set("toDate", toDate)
	}

	var fullURL string
	if len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	} else {
		fullURL = endpoint
	}

	var stats VersionStatsData
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &stats); err != nil {
		return nil, errors.Wrap(err, "failed to get version stats")
	}

	return stats, nil
}

// GetProjectPage retrieves a specific page content from a project.
func (c *Client) GetProjectPage(ctx context.Context, slug, pagePath string) (*Page, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}
	if pagePath == "" {
		pagePath = "home" // Default to home page
	}

	endpoint := fmt.Sprintf("%s/projects/%s/pages/%s",
		c.baseURL, url.PathEscape(slug), url.PathEscape(pagePath))

	var page Page
	if err := c.doRequest(ctx, http.MethodGet, endpoint, nil, &page); err != nil {
		return nil, errors.Wrap(err, "failed to get project page")
	}

	return &page, nil
}

// GetProjectMainPage retrieves the main (home) page of a project.
func (c *Client) GetProjectMainPage(ctx context.Context, slug string) (*Page, error) {
	return c.GetProjectPage(ctx, slug, "home")
}

// GetLatestVersion retrieves the latest version of a project with optional filters.
// channel, platform, and minecraftVersion are all optional filters.
func (c *Client) GetLatestVersion(ctx context.Context, slug, channel, platform, minecraftVersion string) (*Version, error) {
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/projects/%s/latest", c.baseURL, url.PathEscape(slug))

	// Build query parameters
	params := url.Values{}
	if channel != "" {
		params.Set("channel", channel)
	}
	if platform != "" {
		params.Set("platform", platform)
	}
	if minecraftVersion != "" {
		params.Set("platformVersion", minecraftVersion)
	}

	var fullURL string
	if len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	} else {
		fullURL = endpoint
	}

	var version Version
	if err := c.doRequest(ctx, http.MethodGet, fullURL, nil, &version); err != nil {
		return nil, errors.Wrap(err, "failed to get latest version")
	}

	return &version, nil
}

// GetLatestReleaseVersion retrieves the latest release version of a project.
func (c *Client) GetLatestReleaseVersion(ctx context.Context, slug string) (*Version, error) {
	return c.GetLatestVersion(ctx, slug, "Release", "", "")
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
