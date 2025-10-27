package hangar_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lexfrei/go-hangar/pkg/hangar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProject_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		jsonFile     string
		wantID       int64
		wantName     string
		wantSlug     string
		wantCategory string
		wantViews    int64
		wantErr      bool
	}{
		{
			name:         "valid project response",
			jsonFile:     "project_response.json",
			wantID:       1950,
			wantName:     "FancyGlow",
			wantSlug:     "fancyglow",
			wantCategory: "gameplay",
			wantViews:    3618,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(filepath.Join("../../testdata", tt.jsonFile))
			require.NoError(t, err, "failed to read test data file")

			var project hangar.Project
			err = json.Unmarshal(data, &project)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantID, project.ID)
			assert.Equal(t, tt.wantName, project.Name)
			assert.Equal(t, tt.wantSlug, project.Namespace.Slug)
			assert.Equal(t, tt.wantCategory, project.Category)
			assert.Equal(t, tt.wantViews, project.Stats.Views)
			assert.NotZero(t, project.CreatedAt)
		})
	}
}

func TestProjectsList_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		jsonFile      string
		wantCount     int64
		wantLimit     int
		wantOffset    int
		wantResults   int
		wantFirstName string
		wantErr       bool
	}{
		{
			name:          "valid projects list response",
			jsonFile:      "projects_list_response.json",
			wantCount:     2426,
			wantLimit:     25,
			wantOffset:    0,
			wantResults:   1,
			wantFirstName: "FancyGlow",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(filepath.Join("../../testdata", tt.jsonFile))
			require.NoError(t, err, "failed to read test data file")

			var response hangar.ProjectsList
			err = json.Unmarshal(data, &response)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, response.Pagination.Count)
			assert.Equal(t, tt.wantLimit, response.Pagination.Limit)
			assert.Equal(t, tt.wantOffset, response.Pagination.Offset)
			assert.Len(t, response.Result, tt.wantResults)
			if tt.wantResults > 0 {
				assert.Equal(t, tt.wantFirstName, response.Result[0].Name)
			}
		})
	}
}

func TestStats_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		stats   hangar.Stats
		wantErr bool
	}{
		{
			name: "valid stats",
			stats: hangar.Stats{
				Views:           100,
				Downloads:       50,
				RecentViews:     10,
				RecentDownloads: 5,
				Stars:           3,
				Watchers:        2,
			},
			wantErr: false,
		},
		{
			name: "negative values should not cause error in unmarshaling",
			stats: hangar.Stats{
				Views:           -1,
				Downloads:       -1,
				RecentViews:     0,
				RecentDownloads: 0,
				Stars:           0,
				Watchers:        0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(tt.stats)
			require.NoError(t, err)

			var stats hangar.Stats
			err = json.Unmarshal(data, &stats)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.stats.Views, stats.Views)
			assert.Equal(t, tt.stats.Downloads, stats.Downloads)
		})
	}
}

func TestNamespace_Fields(t *testing.T) {
	t.Parallel()

	ns := hangar.Namespace{
		Owner: "testowner",
		Slug:  "testplugin",
	}

	assert.Equal(t, "testowner", ns.Owner)
	assert.Equal(t, "testplugin", ns.Slug)
}

func TestProject_TimeFields(t *testing.T) {
	t.Parallel()

	now := time.Now()
	project := hangar.Project{
		CreatedAt:   now,
		LastUpdated: now,
	}

	assert.Equal(t, now, project.CreatedAt)
	assert.Equal(t, now, project.LastUpdated)
}
