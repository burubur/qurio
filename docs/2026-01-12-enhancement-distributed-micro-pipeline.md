This technical requirements document outlines the implementation of the Distributed Embedding Micro-Pipeline for the Qurio ingestion engine. This architecture decouples high-speed coordination (Link Discovery) from high-latency I/O operations (Gemini Vectorization and Weaviate Storage) by utilizing a durable message queue.

1. Architectural Overview
The system transitions from a sequential processing model to a Durable Coordinator-Worker model.
• The Coordinator (Go Backend): Focused on the "Fast Path," it processes crawling results to update the frontier and immediately enqueues content for vectorization without waiting for API responses.
• The Shock Absorber (NSQ): A new topic, ingest.embed, persists text chunks, protecting against system crashes and managing backpressure from Gemini API rate limits.
• The Embedding Worker (Go Backend): A specialized, horizontally scalable consumer that handles vector generation and storage.
2. Updated Architecture Diagram
This diagram incorporates the recently implemented Topic Splitting for web and file tasks, alongside the new embedding pipeline.
graph TD
    A[NSQ: ingest.result] --> B[ResultConsumer / Coordinator]
    
    subgraph "Fast Path: Coordination & Discovery"
        B --> C{Discover Links}
        C --> D[DB: Create Pending Pages]
        D --> E{Task Type?}
        E -->|Web| E1[NSQ: ingest.task.web]
        E -->|File| E2[NSQ: ingest.task.file]
        E1 & E2 -->|Crawl Continues| F[Python Workers]
    end

    subgraph "Heavy Path: Durable Embedding Pipeline"
        B --> G[NSQ: ingest.embed]
        G --> H[Embedding Worker 1]
        G --> I[Embedding Worker 2]
        G --> J[Embedding Worker N]
        H & I & J --> K[Gemini API]
        K --> L[Weaviate Storage]
    end

    style G fill:#f9f,stroke:#333,stroke-width:2px
    style B fill:#3B82F6,color:#fff
3. Functional Requirements
3.1 Topic and Channel Management
• Topic Seeding: The bootstrap logic in apps/backend/internal/app/bootstrap.go MUST be updated to initialize the ingest.embed topic on startup.
• Channel Design: The embedding workers will use a stable channel named backend-embedder to ensure exactly-once delivery within the worker group.
3.2 The Coordinator (Fast Path)
• Non-Blocking Handoff: Upon receiving a result from ingest.result, the ResultConsumer MUST immediately execute DiscoverLinks and publish new tasks to ingest.task.web.
• Chunk Persistence: The Coordinator MUST split content using text.ChunkMarkdown and publish these chunks to ingest.embed immediately.
• Transactional Durability: The Coordinator MUST NOT acknowledge (FIN) the ingest.result message until chunks are successfully persisted in the ingest.embed queue.
3.3 The Embedding Worker (Heavy Path)
• Role Specialization: The system MUST utilize environment toggles (ENABLE_API and ENABLE_EMBEDDER_WORKER) to allow a single Go binary to act as either the API/Coordinator or the Scaling Worker.
• Concurrent Handling: The worker MUST utilize AddConcurrentHandlers to process multiple chunks in parallel, governed by the INGESTION_CONCURRENCY variable.
• Contextual Integrity: The worker MUST reconstruct the composite embedding string (Title, Source, Path, Type) before vectorization to ensure semantic precision.
4. Scaling and Configuration
Variable
Location
Purpose
ENABLE_API
docker-compose.yml
Role Toggle: Enables HTTP/MCP and Link Discovery.
ENABLE_EMBEDDER_WORKER
docker-compose.yml
Role Toggle: Enables the ingest.embed consumer.
INGESTION_CONCURRENCY
.env
Vertical Scaling: Controls internal goroutines per container.
deploy.replicas
docker-compose.yml
Horizontal Scaling: Controls the number of worker containers.
5. Reliability and Error Handling
5.1 Data Safety
• Late Acknowledgment: The Embedding Worker MUST NOT acknowledge (FIN) the NSQ message until the vector has been successfully stored in Weaviate.
• Poison Pill Protection: The worker MUST acknowledge and discard malformed JSON payloads to prevent infinite retry loops in the queue.
5.2 I/O Isolation
• Hard Timeouts: All calls to the Gemini API and Weaviate MUST be wrapped in a 60-second context.WithTimeout to prevent "Zombie" workers from clogging the pipeline.
• Circuit Breaking: The worker SHOULD implement a backoff strategy for HTTP 429 (Too Many Requests) errors from the Gemini API.


Analogy: This implementation is like upgrading a busy kitchen. Previously, the head chef (the Backend) had to stop everything to hand-deliver every plate to a table (the API) before they could start cooking the next order. With the Micro-Pipeline, the head chef quickly writes orders on sturdy tickets (NSQ). Now, ten line cooks (Embedding Workers) can grab those tickets and cook in parallel, while the head chef stays focused on taking more orders, ensuring no customer is left waiting at the door.
