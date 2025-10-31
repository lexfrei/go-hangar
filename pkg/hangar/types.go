// Package hangar provides a client for interacting with the PaperMC Hangar API.
package hangar

import "time"

// Project represents a plugin/mod project on Hangar.
type Project struct {
	// ID is the unique identifier for the project.
	ID int64 `json:"id"`
	// Name is the display name of the project.
	Name string `json:"name"`
	// Namespace contains owner and slug information.
	Namespace Namespace `json:"namespace"`
	// Category is the project category (e.g., "gameplay", "admin_tools").
	Category string `json:"category"`
	// Description is a short description of the project.
	Description string `json:"description"`
	// CreatedAt is when the project was created.
	CreatedAt time.Time `json:"createdAt"`
	// LastUpdated is when the project was last modified.
	LastUpdated time.Time `json:"lastUpdated"`
	// Stats contains download and view statistics.
	Stats Stats `json:"stats"`
	// Visibility indicates if project is "public" or "private".
	Visibility string `json:"visibility"`
	// AvatarURL is the URL to the project avatar image.
	AvatarURL string `json:"avatarUrl"`
	// Settings contains additional project configuration.
	Settings Settings `json:"settings"`
}

// Namespace identifies the owner and unique slug of a project.
type Namespace struct {
	// Owner is the username of the project owner.
	Owner string `json:"owner"`
	// Slug is the URL-friendly identifier for the project.
	Slug string `json:"slug"`
}

// Stats contains engagement metrics for a project.
type Stats struct {
	// Views is the total number of project page views.
	Views int64 `json:"views"`
	// Downloads is the total number of downloads across all versions.
	Downloads int64 `json:"downloads"`
	// RecentViews is the number of views in the recent period.
	RecentViews int64 `json:"recentViews"`
	// RecentDownloads is the number of downloads in the recent period.
	RecentDownloads int64 `json:"recentDownloads"`
	// Stars is the number of stars/favorites.
	Stars int64 `json:"stars"`
	// Watchers is the number of users watching the project.
	Watchers int64 `json:"watchers"`
}

// Settings contains project configuration and metadata.
type Settings struct {
	// Links contains external links (homepage, source, issues, etc.).
	Links []Link `json:"links"`
	// Tags are project tags for categorization.
	Tags []string `json:"tags"`
	// License contains licensing information.
	License License `json:"license"`
	// Keywords are search keywords for the project.
	Keywords []string `json:"keywords"`
	// Donation contains donation configuration.
	Donation Donation `json:"donation"`
}

// Link represents an external link associated with the project.
type Link struct {
	// ID is the link identifier.
	ID int64 `json:"id"`
	// Name is the link name/type.
	Name string `json:"name"`
	// URL is the link destination.
	URL string `json:"url"`
}

// License contains project license information.
type License struct {
	// Name is the license name (e.g., "MIT", "GPL").
	Name *string `json:"name"`
	// URL is a link to the license text.
	URL *string `json:"url"`
	// Type is the license type identifier.
	Type string `json:"type"`
}

// Donation contains donation/sponsorship configuration.
type Donation struct {
	// Enable indicates if donations are enabled.
	Enable bool `json:"enable"`
	// Subject is the donation subject/message.
	Subject string `json:"subject,omitempty"`
}

// ProjectsList represents a paginated list of projects.
type ProjectsList struct {
	// Pagination contains pagination metadata.
	Pagination Pagination `json:"pagination"`
	// Result is the list of projects in this page.
	Result []Project `json:"result"`
}

// Pagination contains metadata for paginated responses.
type Pagination struct {
	// Count is the total number of items available.
	Count int64 `json:"count"`
	// Limit is the maximum number of items per page.
	Limit int `json:"limit"`
	// Offset is the starting position for this page.
	Offset int `json:"offset"`
}

// Version represents a specific version of a project.
type Version struct {
	// ID is the unique identifier for this version.
	ID int64 `json:"id"`
	// ProjectID is the ID of the parent project.
	ProjectID int64 `json:"projectId"`
	// Name is the version name (e.g., "1.0.0", "2.1-SNAPSHOT").
	Name string `json:"name"`
	// Description is a changelog or description of changes.
	Description string `json:"description"`
	// CreatedAt is when the version was created.
	CreatedAt time.Time `json:"createdAt"`
	// Author is the username of the version author.
	Author string `json:"author"`
	// Visibility is the version visibility ("public", "unlisted", etc.).
	Visibility string `json:"visibility"`
	// ReviewState is the review status ("reviewed", "under_review", etc.).
	ReviewState string `json:"reviewState"`
	// Stats contains download statistics for this version.
	Stats VersionStats `json:"stats"`
	// Downloads contains platform-specific download information.
	Downloads map[string]DownloadInfo `json:"downloads"`
	// PluginDependencies lists required plugin dependencies per platform.
	PluginDependencies map[string][]PluginDependency `json:"pluginDependencies"`
	// Channel contains channel information.
	Channel Channel `json:"channel"`
	// PinnedStatus indicates if version is pinned ("CHANNEL", "VERSION", "NONE").
	PinnedStatus string `json:"pinnedStatus"`
}

// VersionStats contains download statistics for a version.
type VersionStats struct {
	// TotalDownloads is the total download count.
	TotalDownloads int64 `json:"totalDownloads"`
	// PlatformDownloads is downloads per platform.
	PlatformDownloads map[string]int64 `json:"platformDownloads"`
}

// DownloadInfo contains download information for a specific platform.
type DownloadInfo struct {
	// FileInfo contains file metadata (may be null for external URLs).
	FileInfo *FileInfo `json:"fileInfo"`
	// ExternalURL is a direct download URL (e.g., from Modrinth, GitHub).
	ExternalURL string `json:"externalUrl"`
	// DownloadURL is the Hangar-hosted download URL.
	DownloadURL string `json:"downloadUrl"`
}

// PluginDependency represents a required plugin dependency.
type PluginDependency struct {
	// Name is the plugin name.
	Name string `json:"name"`
	// ProjectID is the Hangar project ID if available.
	ProjectID *int64 `json:"projectId"`
	// Required indicates if this dependency is mandatory.
	Required bool `json:"required"`
	// ExternalURL is a link to the plugin if not on Hangar.
	ExternalURL string `json:"externalUrl"`
	// Platform is the platform this dependency applies to.
	Platform string `json:"platform"`
}

// Channel represents a version release channel.
type Channel struct {
	// Name is the channel name (e.g., "Release", "Beta", "Alpha").
	Name string `json:"name"`
	// Description explains the channel purpose.
	Description string `json:"description"`
	// Color is the hex color code for the channel.
	Color string `json:"color"`
	// Flags are channel configuration flags.
	Flags []string `json:"flags"`
	// CreatedAt is when the channel was created.
	CreatedAt time.Time `json:"createdAt"`
}

// FileInfo contains metadata about a version's downloadable file.
type FileInfo struct {
	// Name is the filename.
	Name string `json:"name"`
	// SizeBytes is the file size in bytes.
	SizeBytes int64 `json:"sizeBytes"`
	// SHA256Hash is the SHA-256 checksum of the file.
	SHA256Hash string `json:"sha256Hash"`
}

// VersionsList represents a paginated list of versions.
type VersionsList struct {
	// Pagination contains pagination metadata.
	Pagination Pagination `json:"pagination"`
	// Result is the list of versions in this page.
	Result []Version `json:"result"`
}

// User represents a Hangar user.
type User struct {
	// Name is the username.
	Name string `json:"name"`
	// TagLine is the user's tagline/bio.
	TagLine string `json:"tagline"`
	// JoinDate is when the user joined Hangar.
	JoinDate time.Time `json:"joinDate"`
	// Roles are the user's global roles.
	Roles []Role `json:"roles"`
	// ProjectCount is the number of projects owned.
	ProjectCount int `json:"projectCount"`
	// Locked indicates if the account is locked.
	Locked bool `json:"locked"`
	// AvatarURL is the URL to the user's avatar.
	AvatarURL string `json:"avatarUrl,omitempty"`
	// Socials contains social media links (platform name -> username/handle).
	Socials map[string]string `json:"socials,omitempty"`
}

// Role represents a user role.
type Role struct {
	// Name is the role name.
	Name string `json:"name"`
	// Category is the role category.
	Category string `json:"category"`
	// Color is the hex color code for the role.
	Color string `json:"color"`
}

// UserList represents a paginated list of users.
type UserList struct {
	// Pagination contains pagination metadata.
	Pagination Pagination `json:"pagination"`
	// Result is the list of users in this page.
	Result []User `json:"result"`
}

// ProjectMember represents a project team member.
type ProjectMember struct {
	// User is the member's username.
	User string `json:"user"`
	// Roles are the member's project roles.
	Roles []Role `json:"roles"`
	// Accepted indicates if the member accepted the invite.
	Accepted bool `json:"accepted"`
}

// MemberList represents a paginated list of project members.
type MemberList struct {
	// Pagination contains pagination metadata.
	Pagination Pagination `json:"pagination"`
	// Result is the list of members in this page.
	Result []ProjectMember `json:"result"`
}

// ProjectStats represents daily statistics for a project.
// Map key is date in YYYY-MM-DD format.
type ProjectStats map[string]DailyStats

// VersionStatsData represents daily statistics for a version.
// Map key is date in YYYY-MM-DD format.
type VersionStatsData map[string]DailyStats

// DailyStats contains statistics for a single day.
type DailyStats struct {
	// Downloads is the number of downloads on this day.
	Downloads int64 `json:"downloads"`
	// Views is the number of views on this day.
	Views int64 `json:"views"`
}

// Page represents a project page with Markdown content.
type Page struct {
	// ID is the page identifier.
	ID int64 `json:"id"`
	// Name is the page name.
	Name string `json:"name"`
	// Slug is the URL-friendly page identifier.
	Slug string `json:"slug"`
	// Contents is the Markdown content.
	Contents string `json:"contents"`
}

// Author represents an author with project information.
type Author struct {
	// Name is the author username.
	Name string `json:"name"`
	// JoinDate is when the author joined.
	JoinDate time.Time `json:"joinDate"`
	// ProjectCount is the number of projects.
	ProjectCount int `json:"projectCount"`
	// Roles are the author's roles.
	Roles []Role `json:"roles"`
}

// AuthorList represents a paginated list of authors.
type AuthorList struct {
	// Pagination contains pagination metadata.
	Pagination Pagination `json:"pagination"`
	// Result is the list of authors in this page.
	Result []Author `json:"result"`
}

// StaffMember represents a Hangar staff member.
type StaffMember struct {
	// Name is the staff member username.
	Name string `json:"name"`
	// Roles are the staff roles.
	Roles []Role `json:"roles"`
	// JoinDate is when they joined.
	JoinDate time.Time `json:"joinDate"`
}

// CompactProject is a simplified project representation for user lists.
type CompactProject struct {
	// Namespace contains owner and slug.
	Namespace Namespace `json:"namespace"`
	// Stats contains project statistics.
	Stats Stats `json:"stats"`
	// Category is the project category.
	Category string `json:"category"`
}
