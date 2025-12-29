# Distributed Crypto Backtesting Engine

A high-performance, distributed quantitative finance platform designed to backtest trading strategies against massive historical datasets. This system leverages a split-plane architecture, separating orchestration (Go) from heavy computation (C++), and utilizes columnar storage (ClickHouse) to minimize read-latency during multi-year simulations.

## 🏗 Architecture

The system follows a **Cloud-Native Split-Plane Architecture**:

* **Control Plane (Go):** A stateless REST/gRPC API gateway that manages job scheduling, metadata indexing, and user requests. It acts as the "brain," dispatching logical units of work to the cluster.

* **Compute Plane (C++):** Ephemeral, containerized worker pods orchestrated by Kubernetes. These "muscle" units stream compressed columnar data directly into memory, execute vectorized strategy logic, and push results back to storage.

* **Data Lake (ClickHouse):** An OLAP database optimized for time-series financial data. It replaces traditional row-stores (like Cassandra) to eliminate read amplification and provide 10-100x faster data retrieval for analytical workloads.

## 🚀 Tech Stack

| **Component** | **Technology** | **Rationale** |
|---------------|----------------|---------------|
| **Orchestration** | **Kubernetes** | Manages the lifecycle of ephemeral C++ worker pods and scales execution horizontally. |
| **Control Plane** | **Go (Golang)** | Handles high-concurrency API requests, K8s API interaction, and IO-bound orchestration tasks. |
| **Execution Engine** | **C++ (C++20)** | Provides low-level memory control, SIMD vectorization, and zero-copy data streaming for maximum backtest speed. |
| **Storage** | **ClickHouse** | Columnar storage with high compression ratios, optimized for massive range scans (OLAP). |
| **Data Source** | **Binance API** | Industry standard for crypto spot market data (OHLCV). |

## 🧩 System Components

### 1. The Feed Handler (Go)

A persistent background service that connects to exchange Websockets (e.g., Binance).

* **Function:** Normalizes raw JSON/Strings into strict types (`Int64`, `DateTime64`).

* **Optimization:** Buffers ticks in memory to perform batch inserts, preventing "Too Many Parts" errors in ClickHouse.

### 2. The Strategy Engine (C++)

The core executable running inside K8s Jobs.

* **Input:** Receives execution parameters (Symbol, Time Range, Strategy ID) via Environment Variables or gRPC.

* **Data Access:** Uses `clickhouse-cpp` to stream compressed blocks.

* **Memory Model:** Structure of Arrays (SoA) layout to maximize CPU cache hits and auto-vectorization.

### 3. The API Gateway (Go)

The user-facing entry point.

* **Metadata Index:** Checks availability of data before spinning up costly K8s resources.

* **Job Scheduler:** Splits massive time-ranges (e.g., 10 years) into parallel K8s Jobs (e.g., 10 x 1-year jobs).

## 💾 Data Schema

The database design optimizes for **Sequential Replay** rather than Random Access.

### Market Data Table

Designed for the Binance 1-minute candle format.

```sql
CREATE TABLE market_data_1m
(
    symbol LowCardinality(String),  -- Optimized storage for repeated strings
    timestamp DateTime64(3),        -- Millisecond precision
    open UInt64,
    high UInt64,
    low UInt64,
    close UInt64,
    volume UInt64
)
ENGINE = ReplacingMergeTree()       -- Automatically handles data corrections/duplicates
PARTITION BY toYYYYMM(timestamp)    -- efficient drops of old data
ORDER BY (symbol, timestamp);       -- Fast physical seek for time ranges
```

### Backtest Results Table

Stores the output of the C++ engine for asynchronous retrieval.

```sql
CREATE TABLE backtest_results
(
    job_id String,
    strategy_id String,
    total_pnl Int64,
    sharpe_ratio Int64,
    max_drawdown Int64,
    execution_time_ms UInt64,
    created_at DateTime DEFAULT now()
)
ENGINE = MergeTree()
ORDER BY (job_id);
```

## ⚡ Getting Started

### Prerequisites

* Docker & Kubernetes Cluster (Minikube or Kind)
* ClickHouse Instance (Docker)
* Go 1.21+
* C++ Compiler (GCC 12+ / Clang)

### Local Development Setup

1. **Start Database:**

```bash
docker run -d --name clickhouse-server --ulimit nofile=262144:262144 -p 8123:8123 -p 9000:9000 clickhouse/clickhouse-server
```

2. **Initialize Schema:** Connect to ClickHouse and run the SQL commands in the Data Schema section.

3. **Run Feed Handler (Ingest Data):**

```bash
# (Coming soon: Implementation of the Go Feed Handler)
go run cmd/feed_handler/main.go --symbol BTCUSDT --days 30
```

4. **Build Execution Engine:**

```bash
docker build -t strategy-engine:local -f build/Dockerfile.cpp .
```

5. **Trigger Backtest:**

```bash
curl -X POST http://localhost:8080/api/v1/backtest \
  -d '{"symbol": "BTCUSDT", "strategy": "momentum_v1", "range": "2023"}'
```

## 🗺 Roadmap

- [ ] Phase 1: Implement Go Feed Handler & ClickHouse Ingestion.
- [ ] Phase 2: Create Base C++ Strategy Engine with clickhouse-cpp integration.
- [ ] Phase 3: Build K8s Job Dispatcher in Go.
- [ ] Phase 4: Advanced Analytics Dashboard (Frontend).
- [ ] Phase 5: Support for Dynamic Strategy Loading (Plugin Architecture).
