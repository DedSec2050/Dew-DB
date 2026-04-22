# Dew-DB Architecture (Scaffold)

## Core Components

1. Network Manager: owns TCP listener, client sessions, and connection goroutines.
2. Protocol Parser: parses RESP frames into command structures and serializes responses.
3. Execution Engine: validates and routes commands to internal services.
4. Storage Engine: thread-safe in-memory key-value and metadata access.
5. TTL Manager: tracks expirations (passive checks + active sweeps).
6. Persistence (AOF): appends write operations and supports replay on boot.

## High-Level Data Flow

1. Client connects over TCP.
2. Request bytes are parsed as RESP.
3. Engine resolves command and executes handler.
4. Handler reads/writes storage and optional TTL metadata.
5. Mutating commands are appended to AOF.
6. Response is encoded to RESP and sent to client.

## Azure Equivalents (for deployment and platform discussions)

- Managed Redis-like service: Azure Cache for Redis
- Container runtime: Azure Container Instances
- Kubernetes orchestration: Azure Kubernetes Service (AKS)
- Durable file/object storage for logs/snapshots: Azure Blob Storage
- Monitoring and telemetry: Azure Monitor + Application Insights
