# Dew-DB

Dew-DB is a high-performance, in-memory key-value database written in Go, inspired by Redis/Valkey.

This repository is currently scaffolded for the following core systems:

1. Network Manager (TCP + goroutines)
2. Protocol Parser (RESP-compatible)
3. Execution Engine (command routing)
4. Storage Engine (thread-safe in-memory data layer)
5. TTL Manager (active and passive expiration)
6. Persistence (AOF)

## Initial Directory Layout

```text
cmd/dewdb/                  # Entrypoint package (no code yet)
internal/network/           # TCP listener and connection lifecycle
internal/protocol/resp/     # RESP parsing/encoding
internal/engine/            # Command dispatch and execution pipeline
internal/storage/           # Thread-safe data structures and abstractions
internal/ttl/               # Expiration tracking and eviction loops
internal/persistence/aof/   # Append-only log handling
internal/config/            # Runtime configuration loading/validation
configs/                    # Example config files
docs/                       # Architecture and design notes
deploy/azure/               # Azure deployment notes/manifests
tests/integration/          # Integration test suite layout
scripts/                    # Helper scripts
build/                      # Build artifacts/config templates
```

## Notes

- Module path in `go.mod` is intentionally a placeholder and can be updated later.
- No runtime code has been added yet by design.
