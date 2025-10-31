package hangar_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_DefaultConfig(t *testing.T) {
	t.Parallel()

	client := hangar.NewClient(hangar.Config{})

	assert.NotNil(t, client)
}

func TestNewClient_CustomConfig(t *testing.T) {
	t.Parallel()

	cfg := hangar.Config{
		BaseURL: "https://custom.api.test/v1",
		Token:   "test-token",
		Timeout: 5 * time.Second,
	}

	client := hangar.NewClient(cfg)

	assert.NotNil(t, client)
}

func TestClient_GetProject_Success(t *testing.T) {
	t.Parallel()

	// Load test data
	testData, err := os.ReadFile(filepath.Join("../../testdata", "project_response.json"))
	require.NoError(t, err)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/fancyglow", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testData)
	}))
	defer server.Close()

	// Create client with test server
	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	project, err := client.GetProject(ctx, "fancyglow")

	require.NoError(t, err)
	assert.Equal(t, int64(1950), project.ID)
	assert.Equal(t, "FancyGlow", project.Name)
	assert.Equal(t, "fancyglow", project.Namespace.Slug)
	assert.Equal(t, "gameplay", project.Category)
}

func TestClient_GetProject_NotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "project not found"}`))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	_, err := client.GetProject(ctx, "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestClient_GetProject_ContextCanceled(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.GetProject(ctx, "test")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestClient_ListProjects_Success(t *testing.T) {
	t.Parallel()

	testData, err := os.ReadFile(filepath.Join("../../testdata", "projects_list_response.json"))
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "25", query.Get("limit"))
		assert.Equal(t, "0", query.Get("offset"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testData)
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	list, err := client.ListProjects(ctx, hangar.ListOptions{
		Limit:  25,
		Offset: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, int64(2426), list.Pagination.Count)
	assert.Equal(t, 25, list.Pagination.Limit)
	assert.Len(t, list.Result, 1)
	assert.Equal(t, "FancyGlow", list.Result[0].Name)
}

func TestClient_ListProjects_WithCategory(t *testing.T) {
	t.Parallel()

	testData, err := os.ReadFile(filepath.Join("../../testdata", "projects_list_response.json"))
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, "gameplay", query.Get("category"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testData)
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	list, err := client.ListProjects(ctx, hangar.ListOptions{
		Limit:    25,
		Offset:   0,
		Category: "gameplay",
	})

	require.NoError(t, err)
	assert.Len(t, list.Result, 1)
}

func TestClient_ListVersions_Success(t *testing.T) {
	t.Parallel()

	versionsData := `{
		"pagination": {"count": 1, "limit": 25, "offset": 0},
		"result": [{
			"id": 7728,
			"projectId": 1950,
			"name": "2.0.1",
			"description": "Bug fixes",
			"createdAt": "2024-06-30T19:29:53.843453Z",
			"author": "testowner",
			"visibility": "public",
			"reviewState": "reviewed",
			"stats": {"totalDownloads": 69, "platformDownloads": {"PAPER": 69}},
			"downloads": {
				"PAPER": {
					"fileInfo": null,
					"externalUrl": "https://cdn.test.com/file.jar",
					"downloadUrl": ""
				}
			},
			"pluginDependencies": {},
			"channel": {"name": "Release", "description": "Release", "color": "#14b8a6", "flags": [], "createdAt": "2024-04-21T22:07:16.479186Z"},
			"pinnedStatus": "CHANNEL"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/testowner/testplugin/versions", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(versionsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	versions, err := client.ListVersions(ctx, "testowner", "testplugin", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Len(t, versions.Result, 1)
	assert.Equal(t, "2.0.1", versions.Result[0].Name)
	assert.Equal(t, "testowner", versions.Result[0].Author)
}

func TestClient_GetDownloadURL_Success(t *testing.T) {
	t.Parallel()

	versionsData := `{
		"pagination": {"count": 1, "limit": 100, "offset": 0},
		"result": [{
			"id": 7728,
			"projectId": 1950,
			"name": "2.0.1",
			"description": "Bug fixes",
			"createdAt": "2024-06-30T19:29:53.843453Z",
			"author": "testowner",
			"visibility": "public",
			"reviewState": "reviewed",
			"stats": {"totalDownloads": 69, "platformDownloads": {"PAPER": 69}},
			"downloads": {
				"PAPER": {
					"fileInfo": null,
					"externalUrl": "https://cdn.test.com/testplugin-2.0.1.jar",
					"downloadUrl": ""
				}
			},
			"pluginDependencies": {},
			"channel": {"name": "Release", "description": "Release", "color": "#14b8a6", "flags": [], "createdAt": "2024-04-21T22:07:16.479186Z"},
			"pinnedStatus": "CHANNEL"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/testowner/testplugin/versions", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(versionsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
	})

	ctx := context.Background()
	downloadURL, err := client.GetDownloadURL(ctx, "testowner", "testplugin", "2.0.1", "PAPER")

	require.NoError(t, err)
	assert.Contains(t, downloadURL, "testplugin-2.0.1.jar")
}

func TestClient_WithAuthentication(t *testing.T) {
	t.Parallel()

	testData, err := os.ReadFile(filepath.Join("../../testdata", "project_response.json"))
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer test-token-12345", authHeader)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(testData)
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{
		BaseURL: server.URL,
		Token:   "test-token-12345",
	})

	ctx := context.Background()
	_, err = client.GetProject(ctx, "fancyglow")

	require.NoError(t, err)
}

func TestListOptions_Defaults(t *testing.T) {
	t.Parallel()

	opts := hangar.ListOptions{}

	// Test that we can set defaults
	if opts.Limit == 0 {
		opts.Limit = 25
	}

	assert.Equal(t, 25, opts.Limit)
	assert.Equal(t, 0, opts.Offset)
}

// Test Users & Authors methods

func TestClient_ListUsers_Success(t *testing.T) {
	t.Parallel()

	usersData := `{
		"pagination": {"count": 100, "limit": 25, "offset": 0},
		"result": [{
			"name": "testuser",
			"tagline": "Test user tagline",
			"joinDate": "2024-01-01T00:00:00Z",
			"roles": [{"name": "Developer", "color": "#00ff00"}],
			"projectCount": 5,
			"locked": false
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(usersData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.ListUsers(ctx, "", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(100), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
	assert.Equal(t, "testuser", list.Result[0].Name)
}

func TestClient_GetUser_Success(t *testing.T) {
	t.Parallel()

	userData := `{
		"name": "testuser",
		"tagline": "Test user",
		"joinDate": "2024-01-01T00:00:00Z",
		"roles": [{"name": "Developer", "color": "#00ff00"}],
		"projectCount": 5,
		"locked": false
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/testuser", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(userData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	user, err := client.GetUser(ctx, "testuser")

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, 5, user.ProjectCount)
}

func TestClient_GetUserStarred_Success(t *testing.T) {
	t.Parallel()

	projectsData := `{
		"pagination": {"count": 3, "limit": 25, "offset": 0},
		"result": [{
			"id": 1,
			"name": "TestProject",
			"namespace": {"owner": "testowner", "slug": "testproject"},
			"category": "gameplay",
			"stats": {"downloads": 1000, "stars": 50}
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/testuser/starred", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(projectsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetUserStarred(ctx, "testuser", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(3), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
}

func TestClient_GetUserWatching_Success(t *testing.T) {
	t.Parallel()

	projectsData := `{
		"pagination": {"count": 2, "limit": 25, "offset": 0},
		"result": [{
			"id": 1,
			"name": "TestProject",
			"namespace": {"owner": "testowner", "slug": "testproject"},
			"category": "gameplay",
			"stats": {"downloads": 1000, "watchers": 20}
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/testuser/watching", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(projectsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetUserWatching(ctx, "testuser", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Len(t, list.Result, 1)
}

func TestClient_GetUserPinned_Success(t *testing.T) {
	t.Parallel()

	projectsData := `{
		"pagination": {"count": 1, "limit": 25, "offset": 0},
		"result": [{
			"id": 1,
			"name": "PinnedProject",
			"namespace": {"owner": "testuser", "slug": "pinned"},
			"category": "admin_tools"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/testuser/pinned", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(projectsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetUserPinned(ctx, "testuser")

	require.NoError(t, err)
	assert.Len(t, list.Result, 1)
	assert.Equal(t, "PinnedProject", list.Result[0].Name)
}

func TestClient_ListAuthors_Success(t *testing.T) {
	t.Parallel()

	authorsData := `{
		"pagination": {"count": 50, "limit": 25, "offset": 0},
		"result": [{
			"name": "author1",
			"tagline": "Plugin developer",
			"joinDate": "2023-01-01T00:00:00Z",
			"roles": [{"name": "Developer", "color": "#00ff00"}],
			"projectCount": 10,
			"locked": false
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/authors", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(authorsData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.ListAuthors(ctx, hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(50), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
}

func TestClient_ListStaff_Success(t *testing.T) {
	t.Parallel()

	staffData := `[{
		"name": "staffmember",
		"roles": [{"name": "Hangar_Admin", "color": "#ff0000"}],
		"joinDate": "2022-01-01T00:00:00Z"
	}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/staff", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(staffData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	staff, err := client.ListStaff(ctx)

	require.NoError(t, err)
	assert.Len(t, staff, 1)
	assert.Equal(t, "staffmember", staff[0].Name)
}

// Test Version utilities

func TestClient_GetVersionByID_Success(t *testing.T) {
	t.Parallel()

	versionData := `{
		"id": 12345,
		"projectId": 1950,
		"name": "1.0.0",
		"description": "Initial release",
		"createdAt": "2024-01-01T00:00:00Z",
		"author": "testauthor",
		"visibility": "public",
		"reviewState": "reviewed",
		"stats": {"totalDownloads": 100},
		"downloads": {}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/versions/12345", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(versionData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	version, err := client.GetVersionByID(ctx, 12345)

	require.NoError(t, err)
	assert.Equal(t, int64(12345), version.ID)
	assert.Equal(t, "1.0.0", version.Name)
}

func TestClient_GetVersionByHash_Success(t *testing.T) {
	t.Parallel()

	versionData := `{
		"id": 67890,
		"projectId": 1950,
		"name": "2.0.0",
		"description": "Major update",
		"createdAt": "2024-06-01T00:00:00Z",
		"author": "testauthor",
		"visibility": "public",
		"reviewState": "reviewed",
		"stats": {"totalDownloads": 500},
		"downloads": {}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/versions/find/abc123def456", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(versionData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	version, err := client.GetVersionByHash(ctx, "abc123def456")

	require.NoError(t, err)
	assert.Equal(t, int64(67890), version.ID)
	assert.Equal(t, "2.0.0", version.Name)
}

// Test Project social methods

func TestClient_GetProjectMembers_Success(t *testing.T) {
	t.Parallel()

	membersData := `{
		"pagination": {"count": 3, "limit": 25, "offset": 0},
		"result": [{
			"user": "member1",
			"roles": [{"name": "Owner", "color": "#ff0000"}],
			"accepted": true
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/testproject/members", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(membersData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetProjectMembers(ctx, "testproject", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(3), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
	assert.Equal(t, "member1", list.Result[0].User)
}

func TestClient_GetProjectStargazers_Success(t *testing.T) {
	t.Parallel()

	stargazersData := `{
		"pagination": {"count": 50, "limit": 25, "offset": 0},
		"result": [{
			"name": "stargazer1",
			"projectCount": 2,
			"joinDate": "2024-01-01T00:00:00Z"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/testproject/stargazers", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(stargazersData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetProjectStargazers(ctx, "testproject", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(50), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
}

func TestClient_GetProjectWatchers_Success(t *testing.T) {
	t.Parallel()

	watchersData := `{
		"pagination": {"count": 20, "limit": 25, "offset": 0},
		"result": [{
			"name": "watcher1",
			"projectCount": 5,
			"joinDate": "2024-02-01T00:00:00Z"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/testproject/watchers", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(watchersData))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	list, err := client.GetProjectWatchers(ctx, "testproject", hangar.ListOptions{Limit: 25})

	require.NoError(t, err)
	assert.Equal(t, int64(20), list.Pagination.Count)
	assert.Len(t, list.Result, 1)
}

// Test Statistics methods

func TestClient_GetProjectStats_Success(t *testing.T) {
	t.Parallel()

	projectData := `{
		"id": 1,
		"name": "TestProject",
		"namespace": {"owner": "testowner", "slug": "testproject"},
		"category": "test",
		"description": "Test project",
		"createdAt": "2024-01-01T00:00:00Z",
		"lastUpdated": "2024-01-01T00:00:00Z",
		"stats": {"views": 0, "downloads": 0},
		"visibility": "public",
		"avatarUrl": "",
		"settings": {"links": [], "tags": [], "license": {"type": "MIT"}, "keywords": [], "donation": {"enable": false}}
	}`

	statsData := `{
		"2024-01-01": {"downloads": 100, "views": 500},
		"2024-01-02": {"downloads": 150, "views": 600}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/testproject":
			// First request to get project info
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(projectData))
		case "/projects/testowner/testproject/stats":
			// Second request to get stats with author/slug format
			query := r.URL.Query()
			assert.Equal(t, "2024-01-01T00:00:00Z", query.Get("fromDate"))
			assert.Equal(t, "2024-01-02T23:59:59Z", query.Get("toDate"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(statsData))
		default:
			t.Errorf("Unexpected request to %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	stats, err := client.GetProjectStats(ctx, "testproject", "2024-01-01", "2024-01-02")

	require.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, int64(100), stats["2024-01-01"].Downloads)
	assert.Equal(t, int64(500), stats["2024-01-01"].Views)
}

func TestClient_GetVersionStats_Success(t *testing.T) {
	t.Parallel()

	projectData := `{
		"id": 1,
		"name": "TestProject",
		"namespace": {"owner": "testowner", "slug": "testproject"},
		"category": "test",
		"description": "Test project",
		"createdAt": "2024-01-01T00:00:00Z",
		"lastUpdated": "2024-01-01T00:00:00Z",
		"stats": {"views": 0, "downloads": 0},
		"visibility": "public",
		"avatarUrl": "",
		"settings": {"links": [], "tags": [], "license": {"type": "MIT"}, "keywords": [], "donation": {"enable": false}}
	}`

	statsData := `{
		"2024-06-01": {"downloads": 50, "views": 200},
		"2024-06-02": {"downloads": 75, "views": 250}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/testproject":
			// First request to get project info
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(projectData))
		case "/projects/testowner/testproject/versions/1.0.0/stats":
			// Second request to get version stats with author/slug format
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(statsData))
		default:
			t.Errorf("Unexpected request to %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	stats, err := client.GetVersionStats(ctx, "testproject", "1.0.0", "", "")

	require.NoError(t, err)
	assert.Len(t, stats, 2)
	assert.Equal(t, int64(50), stats["2024-06-01"].Downloads)
}

// Test Pages methods

func TestClient_GetProjectPage_Success(t *testing.T) {
	t.Parallel()

	pageContent := "# Welcome\nThis is the home page"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pages/page/testproject", r.URL.Path)
		assert.Equal(t, "home", r.URL.Query().Get("path"))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(pageContent))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	page, err := client.GetProjectPage(ctx, "testproject", "home")

	require.NoError(t, err)
	assert.Equal(t, "home", page.Slug)
	assert.Contains(t, page.Contents, "Welcome")
}

func TestClient_GetProjectMainPage_Success(t *testing.T) {
	t.Parallel()

	pageContent := "# TestProject\nMain documentation"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pages/main/testproject", r.URL.Path)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(pageContent))
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	page, err := client.GetProjectMainPage(ctx, "testproject")

	require.NoError(t, err)
	assert.Equal(t, "README", page.Name)
	assert.Equal(t, "home", page.Slug)
	assert.Contains(t, page.Contents, "TestProject")
}

// Test Latest version shortcuts

func TestClient_GetLatestVersion_Success(t *testing.T) {
	t.Parallel()

	versionData := `{
		"id": 99999,
		"projectId": 1950,
		"name": "3.0.0",
		"description": "Latest release",
		"createdAt": "2024-12-01T00:00:00Z",
		"author": "testauthor",
		"visibility": "public",
		"reviewState": "reviewed",
		"stats": {"totalDownloads": 1000},
		"downloads": {},
		"pluginDependencies": {},
		"channel": {"name": "Release", "description": "", "color": "#00FF00", "flags": [], "createdAt": "2024-01-01T00:00:00Z"},
		"pinnedStatus": "NONE"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/testproject/latest":
			// First request returns plain text version string
			query := r.URL.Query()
			assert.Equal(t, "Release", query.Get("channel"))

			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("3.0.0"))
		case "/projects/testproject/versions/3.0.0":
			// Second request returns full Version object
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(versionData))
		default:
			t.Errorf("Unexpected request to %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	version, err := client.GetLatestVersion(ctx, "testproject", "Release", "", "")

	require.NoError(t, err)
	assert.Equal(t, int64(99999), version.ID)
	assert.Equal(t, "3.0.0", version.Name)
}

func TestClient_GetLatestReleaseVersion_Success(t *testing.T) {
	t.Parallel()

	versionData := `{
		"id": 88888,
		"projectId": 1950,
		"name": "2.5.0",
		"description": "Latest stable",
		"createdAt": "2024-11-01T00:00:00Z",
		"author": "testauthor",
		"visibility": "public",
		"reviewState": "reviewed",
		"stats": {"totalDownloads": 2000},
		"downloads": {},
		"pluginDependencies": {},
		"channel": {"name": "Release", "description": "", "color": "#00FF00", "flags": [], "createdAt": "2024-01-01T00:00:00Z"},
		"pinnedStatus": "NONE"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/projects/testproject/latestrelease":
			// First request returns plain text version string
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("2.5.0"))
		case "/projects/testproject/versions/2.5.0":
			// Second request returns full Version object
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(versionData))
		default:
			t.Errorf("Unexpected request to %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := hangar.NewClient(hangar.Config{BaseURL: server.URL})
	ctx := context.Background()

	version, err := client.GetLatestReleaseVersion(ctx, "testproject")

	require.NoError(t, err)
	assert.Equal(t, int64(88888), version.ID)
	assert.Equal(t, "2.5.0", version.Name)
}
