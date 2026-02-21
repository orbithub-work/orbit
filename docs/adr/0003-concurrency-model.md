# 3. Refined Concurrency Model

Date: 2026-02-10

## Status

Accepted

## Context

The initial implementation needed to ensure that long-running operations (like workspace initialization) don't block the UI thread or other read operations.

## Decision

We will use Go's native concurrency primitives (Goroutines and Channels) to handle background tasks and protect shared state.

- **Services**: `FileService`, `AssetService`, etc., will be stateless or handle their own internal concurrency using `sync.Mutex` or `sync.RWMutex`.
- **Background Tasks**: Long-running tasks like scanning will run in separate Goroutines.
- **Communication**: Use Channels to report progress or events from background tasks to the UI or other services.
- **AppState**: Global application state will be managed using a thread-safe structure, with granular locking where necessary.

## Consequences

### Positive
- Improved concurrency: Long-running service operations don't block the UI.
- Better separation of concerns.

### Negative
- Breaking change for existing API handlers (need update).
