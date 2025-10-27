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
