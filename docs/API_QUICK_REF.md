# Hangar API Quick Reference

## Base URL

`https://hangar.papermc.io/api/v1`

## Authentication

Most read endpoints work WITHOUT auth. Only uploads, key management, permissions, and page editing require auth.

## Common Parameters

### Pagination

- `limit` - Max items (1-25, default: 25)
- `offset` - Starting position (default: 0)

### Sorting (prefix with `-` for descending)

- Projects: `views`, `downloads`, `newest`, `stars`, `updated`, `recent_downloads`, `recent_views`, `slug`
- Users: `username`, `projectCount`, `joinDate`

## Categories

`admin_tools`, `chat`, `dev_tools`, `economy`, `gameplay`, `misc`, `protection`, `world_management`

## Platforms

`PAPER`, `WATERFALL`, `VELOCITY`

## Quick Endpoints

### Projects

```bash
GET /projects?category=gameplay&platform=PAPER&sort=-downloads
GET /projects/{slug}
GET /projects/{slug}/versions
GET /projects/{slug}/versions/{version}/PAPER/download
GET /projects/{slug}/latest?channel=Release
GET /projects/{slug}/members
GET /projects/{slug}/stargazers
GET /projects/{slug}/watchers
GET /projects/{slug}/stats?fromDate=2024-01-01T00:00:00Z&toDate=2024-12-31T23:59:59Z
```

### Users

```bash
GET /users?query=kenny&limit=10
GET /users/{username}
GET /users/{username}/starred
GET /users/{username}/watching
GET /users/{username}/pinned
```

### Versions

```bash
GET /versions/{id}
GET /versions/hash/{sha256}
GET /versions/{id}/PAPER/download
```

### Other

```bash
GET /authors?query=mini
GET /staff
GET /pages/main/{project}
```

## Response Format

All paginated responses:

```json
{
  "pagination": {
    "count": 2434,
    "limit": 25,
    "offset": 0
  },
  "result": [...]
}
```

## Testing Examples

```bash
# Most popular projects
curl 'https://hangar.papermc.io/api/v1/projects?sort=-downloads&limit=5'

# Search for user
curl 'https://hangar.papermc.io/api/v1/users?query=kenny'

# Get project versions
curl 'https://hangar.papermc.io/api/v1/projects/fancyglow/versions'

# Find version by file hash
curl 'https://hangar.papermc.io/api/v1/versions/hash/4a569f01ac5251fcf17a34937a192fa68f2a4854be90baede3f2f6e132be9d4a'

# Get download URL (returns 301 redirect)
curl -I 'https://hangar.papermc.io/api/v1/projects/fancyglow/versions/2.0.1/PAPER/download'
```
