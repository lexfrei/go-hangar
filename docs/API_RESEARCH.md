# PaperMC Hangar API - Complete Research Report

## Executive Summary

The Hangar API is a comprehensive REST API for the PaperMC Hangar plugin repository. Most endpoints work WITHOUT authentication, contrary to the OpenAPI specification which marks them as requiring auth. The API is well-structured with consistent patterns and good pagination support.

**Base URL:** `https://hangar.papermc.io/api/v1`

**Total Endpoints:** 40 (categorized into 8 groups)

## Current Implementation Status

### ✅ Already Implemented (4/40 endpoints)

1. `GET /api/v1/projects` - List projects with filters
2. `GET /api/v1/projects/{slugOrId}` - Get project details  
3. `GET /api/v1/projects/{author}/{slugOrId}/versions` - List versions
4. `GET /api/v1/projects/{author}/{slugOrId}/versions/{nameOrId}/{platform}/download` - Get download URL (indirect)

### ❌ Not Implemented (36/40 endpoints)

## Complete Endpoint Catalog

### 1. PROJECTS (20 endpoints)

#### Public Read Endpoints (No Auth Required - TESTED ✓)

- `GET /api/v1/projects` - Search/list projects
  - Query params: limit, offset, sort, category, platform, owner, query, license, version, tag, member
  - Sort options: views, downloads, newest, stars, updated, recent_downloads, recent_views, slug
  
- `GET /api/v1/projects/{slugOrId}` - Get specific project
  
- `GET /api/v1/projects/{slugOrId}/versions` - List project versions
  - Query params: limit, offset, includeHiddenChannels, channel, platform, platformVersion
  
- `GET /api/v1/projects/{slugOrId}/versions/{nameOrId}` - Get specific version
  - nameOrId can be version name (e.g., "1.0.0") or internal ID
  
- `GET /api/v1/projects/{slugOrId}/versions/{nameOrId}/{platform}/download` - Download version
  - Returns 301 redirect to actual file (CDN or external)
  - Platform: PAPER, WATERFALL, VELOCITY
  
- `GET /api/v1/projects/{slugOrId}/latest` - Get latest version (any channel)
  - Query param: channel (optional, defaults to all)
  
- `GET /api/v1/projects/{slugOrId}/latestrelease` - Get latest release version
  
- `GET /api/v1/projects/{slugOrId}/members` - List project members
  - Query params: limit, offset
  
- `GET /api/v1/projects/{slugOrId}/stargazers` - List users who starred
  - Query params: limit, offset
  
- `GET /api/v1/projects/{slugOrId}/watchers` - List users watching
  - Query params: limit, offset
  
- `GET /api/v1/projects/{slugOrId}/stats` - Get project statistics
  - Query params: fromDate*, toDate* (ISO 8601 format: 2024-01-01T00:00:00Z)
  - Returns daily download/view counts
  
- `GET /api/v1/projects/{author}/{slugOrId}/versions/{nameOrId}/stats` - Get version statistics
  - Query params: fromDate*, toDate*

#### Alternative Paths (author/slug format)

All project endpoints also support `{author}/{slugOrId}` format:

- `/api/v1/projects/{author}/{slugOrId}/latest`
- `/api/v1/projects/{author}/{slugOrId}/latestrelease`
- `/api/v1/projects/{author}/{slugOrId}/versions`
- `/api/v1/projects/{author}/{slugOrId}/versions/{nameOrId}`
- `/api/v1/projects/{author}/{slugOrId}/versions/{nameOrId}/stats`
- `/api/v1/projects/{author}/{slugOrId}/versions/{nameOrId}/{platform}/download`

#### Write Endpoints (Auth Required)

- `POST /api/v1/projects/{slugOrId}/upload` - Upload new version
- `POST /api/v1/projects/{author}/{slugOrId}/upload` - Upload new version (alt)

### 2. VERSIONS (4 endpoints)

- `GET /api/v1/versions/{id}` - Get version by internal ID (works without auth ✓)
- `GET /api/v1/versions/{id}/stats` - Get version stats by ID
  - Query params: fromDate*, toDate*
- `GET /api/v1/versions/{id}/{platform}/download` - Download by version ID
- `GET /api/v1/versions/hash/{hash}` - Find version by SHA-256 file hash (works without auth ✓)
  - Returns full project info with the matching version

### 3. USERS (5 endpoints)

All tested working WITHOUT authentication:

- `GET /api/v1/users` - Search/list users
  - Query params: query*, limit, offset, sort
  - Sort options: username, projectCount, joinDate
  
- `GET /api/v1/users/{user}` - Get specific user details
  
- `GET /api/v1/users/{user}/pinned` - Get user's pinned projects
  
- `GET /api/v1/users/{user}/starred` - Get projects user starred
  - Query params: limit, offset, sort
  - Sort options: -stars, -downloads, -updated, -newest
  
- `GET /api/v1/users/{user}/watching` - Get projects user is watching
  - Query params: limit, offset, sort

### 4. AUTHORS (1 endpoint)

- `GET /api/v1/authors` - List users with at least one public project (works without auth ✓)
  - Query params: query*, limit, offset, sort
  - Sort options: username, projectCount, joinDate

### 5. STAFF (1 endpoint)

- `GET /api/v1/staff` - List Hangar staff members (works without auth ✓)
  - Query params: query*, limit, offset, sort

### 6. PAGES (4 endpoints)

- `GET /api/v1/pages/main/{project}` - Get project main page content (Markdown, works without auth ✓)
- `GET /api/v1/pages/page/{project}` - Get specific project page
  - Query param: path* (page path)
- `PATCH /api/v1/pages/editmain/{project}` - Edit main page (auth required)
- `PATCH /api/v1/pages/edit/{project}` - Edit page (auth required)

### 7. AUTHENTICATION & KEYS (4 endpoints)

- `POST /api/v1/authenticate` - Create JWT from API key
  - Query param: apiKey*
- `GET /api/v1/keys` - List API keys (auth required)
- `POST /api/v1/keys` - Create new API key (auth required)
- `DELETE /api/v1/keys` - Delete API key (auth required)
  - Query param: name*

### 8. PERMISSIONS (3 endpoints)

All require authentication:

- `GET /api/v1/permissions` - Get your permissions
  - Query params: slug, organization, project
- `GET /api/v1/permissions/hasAll` - Check if you have all permissions
  - Query params: permissions*, slug, organization, project
- `GET /api/v1/permissions/hasAny` - Check if you have any permission
  - Query params: permissions*, slug, organization, project

## API Features & Characteristics

### Pagination

- Standard params: `limit` (default: 25), `offset` (default: 0)
- Response includes: `pagination: {count, limit, offset}` and `result: []`

### Filtering & Searching

**Projects:**

- `category`: admin_tools, chat, dev_tools, economy, gameplay, misc, protection, world_management
- `platform`: PAPER, WATERFALL, VELOCITY
- `version`: Minecraft version (e.g., "1.20", "1.21")
- `license`: License type filter (e.g., "MIT", "GPL")
- `owner`: Filter by project owner username
- `query`: Text search
- `tag`: Filter by tag
- `member`: Filter by member username

**Versions:**

- `channel`: Filter by channel name (e.g., "Release", "Beta", "Alpha")
- `platform`: PAPER, WATERFALL, VELOCITY
- `platformVersion`: Minecraft version
- `includeHiddenChannels`: Include hidden channels (default: true)

### Sorting

Projects support: views, downloads, newest, stars, updated, recent_downloads, recent_views, slug

- Prefix with `-` for descending order (e.g., `-downloads`)

### Statistics

- Date range queries require ISO 8601 format: `2024-01-01T00:00:00Z`
- Returns daily breakdown of downloads/views
- Available for both projects and versions

### Error Handling

- Standard HTTP status codes
- JSON error responses with structure:

  ```json
  {
    "message": "Error description",
    "messageArgs": [],
    "isHangarApiException": true,
    "httpError": {"statusCode": 404}
  }
  ```

### Rate Limiting

- No rate limiting detected in testing (5 rapid requests all succeeded)
- No rate limit headers observed

### Authentication

- **Important Finding:** Despite OpenAPI spec marking endpoints as auth-required, most read endpoints work WITHOUT authentication
- Only actual auth-required endpoints: uploads, key management, permissions, page editing
- Auth header format: `Authorization: Bearer {token}`

## Data Structures

### Core Types (Already Implemented ✅)

- Project
- ProjectsList
- Version
- VersionsList
- Pagination
- Namespace
- Stats
- VersionStats
- DownloadInfo
- FileInfo
- Channel
- PluginDependency

### Missing Types (Need Implementation ❌)

#### User

```go
type User struct {
    ID             int64     `json:"id"`
    Name           string    `json:"name"`
    Tagline        string    `json:"tagline"`
    Roles          []int     `json:"roles"`  // Role IDs
    ProjectCount   int       `json:"projectCount"`
    Locked         bool      `json:"locked"`
    NameHistory    []UserNameChange `json:"nameHistory"`
    AvatarURL      string    `json:"avatarUrl"`
    CreatedAt      time.Time `json:"createdAt"`
    IsOrganization bool      `json:"isOrganization"`
    Socials        map[string]string `json:"socials"`  // github, website, etc.
}
```

#### ProjectMember

```go
type ProjectMember struct {
    User   string `json:"user"`
    UserID int64  `json:"userId"`
    Roles  []Role `json:"roles"`
}

type Role struct {
    Title    string `json:"title"`
    Color    string `json:"color"`
    Rank     int    `json:"rank"`
    Category string `json:"category"`
}
```

#### ProjectStats (for stats endpoint)

```go
type ProjectStats map[string]DailyStats

type DailyStats struct {
    TotalDownloads    int64            `json:"totalDownloads"`
    PlatformDownloads map[string]int64 `json:"platformDownloads"`
}
```

#### UserList

```go
type UserList struct {
    Pagination Pagination `json:"pagination"`
    Result     []User     `json:"result"`
}
```

#### MemberList

```go
type MemberList struct {
    Pagination Pagination      `json:"pagination"`
    Result     []ProjectMember `json:"result"`
}
```

#### ApiKey

```go
type ApiKey struct {
    Name            string      `json:"name"`
    TokenIdentifier string      `json:"tokenIdentifier"`
    Permissions     []string    `json:"permissions"`
    CreatedAt       time.Time   `json:"createdAt"`
    LastUsed        *time.Time  `json:"lastUsed"`
}
```

## Priority Implementation Plan

### Priority 1: Core Read Operations (High Value, Low Effort)

1. **Users & Authors** (5 endpoints)
   - List/search users
   - Get user details
   - User's starred/watching/pinned projects
   - List authors
   - **Reason:** Complete user discovery, user-centric workflows

2. **Version Lookup** (2 endpoints)
   - Get version by ID
   - Find version by file hash
   - **Reason:** Version validation, file integrity checks

3. **Project Social** (3 endpoints)
   - Stargazers
   - Watchers
   - Members
   - **Reason:** Social features, project team information

### Priority 2: Statistics & Analytics (Medium Value, Medium Effort)

4. **Statistics** (3 endpoints)
   - Project stats with date ranges
   - Version stats with date ranges
   - **Reason:** Analytics, download tracking, trend analysis

5. **Staff List** (1 endpoint)
   - List Hangar staff
   - **Reason:** Low priority but trivial to implement

### Priority 3: Content & Pages (Low Value, Low Effort)

6. **Pages** (2 read endpoints)
   - Get main page
   - Get specific page
   - **Reason:** Full project information, documentation access

### Priority 4: Latest Version Helpers (High Value, Low Effort)

7. **Version Shortcuts** (4 endpoints)
   - Get latest version
   - Get latest release
   - Both with author/slug variants
   - **Reason:** Common use case, simplifies version discovery

### Priority 5: Write Operations (Lower Priority)

8. **Upload** (2 endpoints - auth required)
   - Upload new version
   - **Reason:** Publishing workflow (less common for CLI)

9. **API Keys** (4 endpoints - auth required)
   - Manage API keys
   - **Reason:** Authentication management

10. **Permissions** (3 endpoints - auth required)
    - Check permissions
    - **Reason:** Authorization checks

11. **Page Editing** (2 endpoints - auth required)
    - Edit pages
    - **Reason:** Content management (less common for CLI)

## Implementation Recommendations

### Query Parameter Handling

Add new option types:

```go
type ListOptions struct {
    Limit    int
    Offset   int
    Sort     string
    Query    string
}

type ProjectListOptions struct {
    ListOptions
    Category        string
    Platform        string
    Owner           string
    License         string
    Version         string  // Minecraft version
    Tag             string
    Member          string
}

type VersionListOptions struct {
    ListOptions
    Channel              string
    Platform             string
    PlatformVersion      string
    IncludeHiddenChannels bool
}

type StatsOptions struct {
    FromDate time.Time
    ToDate   time.Time
}
```

### API Client Methods

**Priority 1 Methods:**

```go
// Users
func (c *Client) ListUsers(ctx context.Context, opts ListOptions) (*UserList, error)
func (c *Client) GetUser(ctx context.Context, username string) (*User, error)
func (c *Client) GetUserStarredProjects(ctx context.Context, username string, opts ListOptions) (*ProjectsList, error)
func (c *Client) GetUserWatchingProjects(ctx context.Context, username string, opts ListOptions) (*ProjectsList, error)
func (c *Client) GetUserPinnedProjects(ctx context.Context, username string) ([]Project, error)

// Authors & Staff
func (c *Client) ListAuthors(ctx context.Context, opts ListOptions) (*UserList, error)
func (c *Client) ListStaff(ctx context.Context, opts ListOptions) (*UserList, error)

// Versions
func (c *Client) GetVersionByID(ctx context.Context, id int64) (*Version, error)
func (c *Client) GetVersionByHash(ctx context.Context, hash string) (*Project, error)

// Project Social
func (c *Client) GetProjectMembers(ctx context.Context, slug string, opts ListOptions) (*MemberList, error)
func (c *Client) GetProjectStargazers(ctx context.Context, slug string, opts ListOptions) (*UserList, error)
func (c *Client) GetProjectWatchers(ctx context.Context, slug string, opts ListOptions) (*UserList, error)
```

**Priority 2 Methods:**

```go
// Statistics
func (c *Client) GetProjectStats(ctx context.Context, slug string, opts StatsOptions) (ProjectStats, error)
func (c *Client) GetVersionStats(ctx context.Context, slug, version string, opts StatsOptions) (ProjectStats, error)
func (c *Client) GetVersionStatsByID(ctx context.Context, id int64, opts StatsOptions) (ProjectStats, error)
```

**Priority 3 Methods:**

```go
// Pages
func (c *Client) GetProjectMainPage(ctx context.Context, slug string) (string, error)
func (c *Client) GetProjectPage(ctx context.Context, slug, path string) (string, error)
```

**Priority 4 Methods:**

```go
// Version Helpers
func (c *Client) GetLatestVersion(ctx context.Context, slug string, channel string) (*Version, error)
func (c *Client) GetLatestRelease(ctx context.Context, slug string) (*Version, error)
```

### CLI Commands Structure

```bash
# Users
hangar user get <username>
hangar user list [--query=<text>] [--sort=<field>]
hangar user starred <username>
hangar user watching <username>
hangar user pinned <username>

# Authors & Staff
hangar author list [--query=<text>]
hangar staff list

# Versions
hangar version get <id>
hangar version find-by-hash <sha256>
hangar version latest <slug> [--channel=<name>]
hangar version latest-release <slug>

# Project Social
hangar project members <slug>
hangar project stargazers <slug>
hangar project watchers <slug>

# Statistics
hangar project stats <slug> --from=<date> --to=<date>
hangar version stats <slug> <version> --from=<date> --to=<date>

# Pages
hangar project page <slug> [--path=<page>]
```

## Edge Cases & Quirks

1. **Version Identification:**
   - Versions can be referenced by name OR internal ID
   - Some endpoints require ID, others accept both
   - Latest version endpoints return empty body (not null, just empty)

2. **Pagination:**
   - Empty results return valid pagination with count=0
   - No maximum limit enforced (tested up to 100)

3. **Date Formats:**
   - Stats require ISO 8601 with timezone: `2024-01-01T00:00:00Z`
   - Other dates in responses are ISO 8601 strings

4. **Downloads:**
   - Download endpoints return 301 redirects
   - External URLs (Modrinth, GitHub) vs. Hangar CDN
   - `externalUrl` XOR `downloadUrl` will be set in DownloadInfo

5. **Project Stats:**
   - Endpoint path is `/projects/{slug}/stats` but some projects return 404
   - Likely requires the project to have stats enabled

6. **Author vs Slug:**
   - Most endpoints accept just `{slugOrId}`
   - Some have explicit `{author}/{slugOrId}` variants
   - Both patterns work for versions endpoints

7. **Sort Direction:**
   - Prefix with `-` for descending (e.g., `-downloads`)
   - No prefix for ascending

8. **Empty Responses:**
   - Latest version endpoints may return empty (not 404) if no versions
   - Pinned projects returns empty array if none pinned

## Testing Recommendations

1. **Unit Tests:**
   - Mock HTTP responses for all endpoints
   - Test pagination edge cases (empty, single page, multi-page)
   - Test error handling (404, 400, 500)
   - Test date parsing/formatting

2. **Integration Tests:**
   - Use real API with known stable projects (e.g., "Maintenance", "ViaVersion")
   - Test without authentication for public endpoints
   - Test rate limiting behavior
   - Verify redirect following for downloads

3. **CLI Tests:**
   - Test all output formats (table, JSON)
   - Test pagination with --limit and --offset
   - Test filtering combinations
   - Test error messages for user-friendly output

## Open Questions

1. **Authentication:**
   - What's the actual scope of endpoints requiring auth?
   - OpenAPI spec vs. actual behavior discrepancy

2. **Rate Limiting:**
   - No limits observed, but are there production limits?
   - Should we implement client-side rate limiting?

3. **Caching:**
   - Should we implement response caching?
   - Cache headers not observed in responses

4. **Webhooks:**
   - No webhook endpoints found - is there webhook support?

5. **Organizations:**
   - Organizations vs. Users - permission system?
   - Organization management endpoints?

## Conclusion

The Hangar API is comprehensive and well-designed with 40 endpoints covering all aspects of plugin repository management. The current go-hangar implementation covers only 10% (4/40) of available functionality.

**Recommended Next Steps:**

1. Implement Priority 1 endpoints (users, authors, version lookup, project social) - 13 endpoints
2. Add comprehensive types for User, ProjectMember, ProjectStats
3. Enhance ListOptions to support all filter parameters
4. Add CLI commands for new functionality
5. Write integration tests against real API

This would bring coverage to 42.5% (17/40 endpoints) while delivering the highest-value features for end users.
