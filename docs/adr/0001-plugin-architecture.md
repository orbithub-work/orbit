# 1. Internal Plugin Architecture

Date: 2026-02-10

## Status

Accepted

## Context

We need to allow external tools (like CapCut assistant) to interact with Media Assistant OS. The architecture should support bi-directional communication and be language-agnostic for plugins.

## Decision

We will use an internal HTTP/WebSocket server hosted by the Go application backend.

- **Framework**: `net/http` (Standard Library)
- **Protocol**: HTTP/1.1 for commands, WebSocket for events.
- **Port**: Default 32000, auto-increment if occupied.
- **Discovery**: The app will write its port to a local file, and plugins will read it.

## Consequences

### Positive
- Loose coupling between Core and Plugins.
- Plugins can be written in any language (Python, Node.js, etc.).
- Browser-based plugins are supported.

### Negative
- Security implications (need to restrict to localhost).
- Need to handle port management.
