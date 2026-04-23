# Dew-DB Feature Roadmap

This roadmap is ordered to reduce rework: protocol and correctness first, durability and operations second, then scale-out and advanced features.

## Planning Principles

1. Keep Redis CLI compatibility as a hard requirement for core commands.
2. Ship vertical slices that are production-testable, not isolated subsystems.
3. Add observability before major complexity jumps.
4. Prefer correctness and deterministic behavior over premature optimization.
5. For deployment tracks, map each design choice to Azure equivalents.

## Phase 0: Project Foundation

### Scope

1. Entrypoint and process lifecycle baseline.
2. Configuration loading and validation.
3. Structured logging baseline.
4. CI checks for build, test, lint.

### Deliverables

1. Runtime config file schema and env overrides.
2. Consistent startup/shutdown hooks.
3. Repo automation for fmt, lint, test.

### Exit Criteria

1. App boots and shuts down cleanly.
2. Config validation fails fast on invalid values.
3. CI pipeline passes on clean branch.

## Phase 1: Network + RESP Core

### Scope

1. TCP listener with concurrent connection handling.
2. RESP frame parser and encoder for basic request/response.
3. Basic command dispatcher shell.
4. Connection limits and idle timeout controls.

### Deliverables

1. Stable connection accept loop with graceful stop.
2. RESP decoding for inline and array command forms.
3. End-to-end command loop for a single connection.

### Exit Criteria

1. Redis CLI can connect and issue PING.
2. Invalid RESP payloads return protocol errors without crashing server.
3. Multiple clients can run commands concurrently.

## Phase 2: Core Data Commands (MVP)

### Scope

1. Thread-safe key-value storage engine.
2. Command handlers: PING, GET, SET, DEL, EXISTS.
3. Basic error model and response formatting.

### Deliverables

1. Storage API with clear read/write semantics.
2. Command execution pipeline with validation.
3. Unit tests for storage and command handlers.

### Exit Criteria

1. Correct command behavior under concurrent clients.
2. Basic compatibility with Redis CLI command expectations.
3. Deterministic error responses.

## Phase 3: TTL and Expiration Semantics

### Scope

1. EXPIRE, PEXPIRE, TTL, PTTL commands.
2. Passive expiration during key access.
3. Active expiration background sweeper.

### Deliverables

1. Expiration index structure and scheduler loop.
2. Unified key metadata model for value and expiration.
3. Load tests for expiration behavior under churn.

### Exit Criteria

1. Expired keys are not returned by GET.
2. Active sweeper keeps memory bounded under TTL-heavy workloads.
3. TTL command responses match expected semantics.

## Phase 4: AOF Durability

### Scope

1. Append-only write logging for mutating commands.
2. Configurable fsync policy.
3. Startup replay for crash recovery.

### Deliverables

1. AOF writer abstraction with batching/flush strategy.
2. Recovery bootstrap path from AOF.
3. Corruption handling strategy for truncated tail records.

### Exit Criteria

1. Data survives process restart via AOF replay.
2. Fsync policy behavior is measurable and documented.
3. Replay is idempotent and consistent.

## Phase 5: Performance and Concurrency Upgrades

### Scope

1. Sharded map storage to reduce lock contention.
2. Command pipeline optimization and allocation reduction.
3. Benchmark suite and profiling workflow.

### Deliverables

1. Read/write path benchmarks with baseline and target metrics.
2. pprof instrumentation and profiling runbook.
3. Latency and throughput dashboard definitions.

### Exit Criteria

1. Demonstrated throughput gain versus single-lock baseline.
2. Tail latency targets met under multi-client benchmark.
3. No correctness regressions in integration tests.

## Phase 6: Replication and High Availability (Optional after stable single-node)

### Scope

1. Primary-replica replication protocol.
2. Replica sync and partial resync strategy.
3. Failover coordination design.

### Deliverables

1. Replication backlog model.
2. Snapshot or streamed sync bootstrap path.
3. Documented failure and recovery behavior.

### Exit Criteria

1. Replica catches up after transient disconnect.
2. Read-only replica mode is enforced.
3. Data divergence tests pass in fail/recover scenarios.

## Phase 7: Security and Multi-Tenancy Controls

### Scope

1. AUTH and ACL baseline.
2. TLS support for client connections.
3. Connection quotas and optional namespace isolation.

### Deliverables

1. Credential and permission model.
2. TLS certificate loading and rotation strategy.
3. Security audit checklist.

### Exit Criteria

1. Unauthorized commands are rejected.
2. TLS-enabled clients can connect and run commands.
3. ACL tests cover allow/deny matrix.

## Phase 8: Operations, Packaging, and Azure Deployment

### Scope

1. Container packaging and health endpoints.
2. Metrics, tracing, and alerting.
3. Deployment manifests for Azure environments.

### Deliverables

1. Container image and runtime tuning guide.
2. Metrics export with key SLO signals.
3. Azure deployment assets and runbooks.

### Azure Equivalents

1. Managed Redis alternative for comparison: Azure Cache for Redis.
2. Single-instance container runtime: Azure Container Instances.
3. Cluster orchestration target: Azure Kubernetes Service.
4. Durable persistence target for artifacts/logs: Azure Blob Storage.
5. Monitoring target: Azure Monitor plus Application Insights.

### Exit Criteria

1. Reproducible deployment path for dev and staging.
2. Health and readiness checks integrated into orchestration.
3. Alerting covers crash loops, high latency, and memory pressure.

## Cross-Phase Test Strategy

1. Unit tests: parser, storage, TTL, command handlers.
2. Integration tests: RESP compatibility, multi-client concurrency, restart recovery.
3. Fault tests: malformed protocol, disk-full simulation, forced process kill and replay.
4. Performance tests: throughput, p95/p99 latency, memory usage over time.

## Suggested Initial Implementation Sequence

1. Phase 1 slice: PING end-to-end over RESP.
2. Phase 2 slice: GET and SET with thread-safe storage.
3. Phase 3 slice: EXPIRE plus passive expiration.
4. Phase 4 slice: SET and DEL persistence with replay.
5. Phase 5 slice: shard storage and benchmark improvement.

## Definition of Done for Each Feature

1. Behavior spec written first in docs.
2. Unit and integration tests included.
3. Metrics added for success, error, and latency path.
4. Failure modes documented and tested where practical.
5. No lint or test regressions in CI.
