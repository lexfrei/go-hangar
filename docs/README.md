# go-hangar API Documentation

This directory contains comprehensive documentation about the PaperMC Hangar API and implementation planning.

## Documents

### [API_RESEARCH.md](./API_RESEARCH.md) - Complete API Analysis (557 lines)

**The definitive reference for the Hangar API**

Comprehensive research document covering:

- All 40 API endpoints with detailed descriptions
- Query parameters, filters, and options
- Response structures and data types
- Authentication requirements (TESTED - most endpoints work without auth!)
- Pagination, sorting, and searching
- Edge cases and API quirks
- Complete code examples in Go
- Priority-based implementation recommendations

**Use this for:** Understanding what the API can do, planning implementations, reference during development

### [API_QUICK_REF.md](./API_QUICK_REF.md) - Quick Reference (94 lines)

**Fast lookup for common operations**

Quick reference guide with:

- Common endpoints and usage patterns
- Available categories, platforms, and filters
- Request/response format examples
- curl command examples for testing
- Common query parameters

**Use this for:** Quick lookups during development, testing endpoints with curl, remembering parameter names

### [ROADMAP.md](./ROADMAP.md) - Implementation Roadmap (276 lines)

**Phased implementation plan**

Detailed roadmap covering:

- Current implementation status (10% - 4/40 endpoints)
- 4 implementation phases with priorities
- Required type additions per phase
- CLI command structure
- Effort estimates and timelines
- Success metrics

**Use this for:** Planning sprints, tracking progress, understanding dependencies

## Key Findings

### 🎯 API Coverage

- **Total Endpoints:** 40
- **Currently Implemented:** 4 (10%)
- **Public (No Auth):** ~30 endpoints work without authentication
- **Priority 1 Target:** 17 endpoints (42.5%)

### 🚀 Quick Stats

- **Categories:** 8 (Projects, Versions, Users, Authors, Staff, Pages, Keys, Permissions)
- **Platforms:** 3 (PAPER, WATERFALL, VELOCITY)
- **Project Categories:** 8 (admin_tools, chat, dev_tools, economy, gameplay, misc, protection, world_management)
- **Pagination:** Standard limit/offset (max 25 per page)
- **Rate Limiting:** None observed

### 💡 Important Discoveries

1. **Most endpoints work WITHOUT authentication** - contrary to OpenAPI spec
2. **Version lookup by file hash** available (SHA-256)
3. **Daily statistics** available with date ranges (ISO 8601)
4. **Dual path patterns:** Both `/projects/{slug}` and `/projects/{author}/{slug}` work
5. **Download redirects:** Download endpoints return 301 to actual file (CDN or external)

### 📊 Recommended Implementation Order

1. **Phase 1:** Users, authors, version utilities, project social (13 endpoints) - **2-3 days**
2. **Phase 2:** Statistics with date ranges (4 endpoints) - **1-2 days**
3. **Phase 3:** Pages and version helpers (6 endpoints) - **1 day**
4. **Phase 4:** Write operations and auth (13 endpoints) - **3-4 days**

**Total estimated time:** 7-10 days for 100% coverage

## Using This Documentation

### For Development

1. Start with **ROADMAP.md** to understand phases and priorities
2. Reference **API_RESEARCH.md** for detailed endpoint specifications
3. Use **API_QUICK_REF.md** for quick parameter lookups
4. Test endpoints with curl examples before implementing

### For Testing

1. Use curl examples from **API_QUICK_REF.md**
2. Test without authentication first (most endpoints work)
3. Use stable projects for testing: "Maintenance", "ViaVersion", "Essentials"
4. Reference edge cases in **API_RESEARCH.md**

### For Planning

1. Review current status in **ROADMAP.md**
2. Choose implementation phase based on priorities
3. Check effort estimates and dependencies
4. Track progress against success metrics

## Example Workflow

```bash
# 1. Test an endpoint
curl 'https://hangar.papermc.io/api/v1/projects?sort=-downloads&limit=3'

# 2. Check API_RESEARCH.md for detailed specs
# - Query parameters
# - Response structure
# - Error cases

# 3. Implement in client.go
func (c *Client) ListProjects(ctx context.Context, opts ListOptions) (*ProjectsList, error)

# 4. Add CLI command
hangar project list --sort=-downloads --limit=3

# 5. Write tests
func TestListProjects(t *testing.T) { ... }

# 6. Update ROADMAP.md to mark as complete
```

## Testing Examples

### List Most Downloaded Projects

```bash
curl 'https://hangar.papermc.io/api/v1/projects?sort=-downloads&limit=5' | jq '.result[] | {name, downloads: .stats.downloads}'
```

### Find User's Projects

```bash
curl 'https://hangar.papermc.io/api/v1/users/kennytv/starred' | jq '.result[] | .name'
```

### Lookup Version by Hash

```bash
curl 'https://hangar.papermc.io/api/v1/versions/hash/4a569f01ac5251fcf17a34937a192fa68f2a4854be90baede3f2f6e132be9d4a' | jq '.name'
```

### Get Project Stats

```bash
curl 'https://hangar.papermc.io/api/v1/projects/maintenance/stats?fromDate=2024-01-01T00:00:00Z&toDate=2024-12-31T23:59:59Z' | jq 'to_entries | .[:3]'
```

## Contributing

When implementing new endpoints:

1. Add types to `pkg/hangar/types.go`
2. Add client methods to `pkg/hangar/client.go`
3. Add CLI commands to `internal/cli/`
4. Add tests to `pkg/hangar/client_test.go`
5. Update README.md examples
6. Mark as complete in ROADMAP.md

## Links

- [Hangar Website](https://hangar.papermc.io)
- [Hangar API Docs](https://hangar.papermc.io/api-docs)
- [OpenAPI Spec](https://hangar.papermc.io/v3/api-docs)
- [PaperMC Discord](https://discord.gg/papermc)
- [go-hangar Repository](https://github.com/lexfrei/go-hangar)
