# 2. Unified Caching Strategy

Date: 2026-02-10

## Status

Accepted

## Context

Performance bottlenecks were observed in:
1. Checking if thumbnails exist on disk (repeated I/O).
2. Complex search queries (repeated SQL execution).

## Decision

We will introduce a unified in-memory caching layer.

- **Thumbnail Cache**: Stores `(path, size) -> thumbnail_path`. TTL: 1 hour.
- **Search Cache**: Stores `query -> results`. TTL: 5 minutes.
- **Implementation**: `CacheManager` struct in `internal/services` layer, using Go's `sync.Map` or a thread-safe cache library.

## Consequences

### Positive
- Reduced disk I/O and database load.
- Faster UI response for repeated actions.

### Negative
- Increased memory usage.
- Cache invalidation complexity (need to clear search cache on file changes).
