CREATE DATABASE IF NOT EXISTS crypto;

USE crypto;

-- ==========================================
-- 1. MARKET DATA TABLE (ReplacingMergeTree)
-- ==========================================
-- Stores raw OHLCV data. 
-- 'ReplacingMergeTree' handles duplicate inserts/corrections automatically.
-- Partitioned by Month for efficient data lifecycle management.
CREATE TABLE IF NOT EXISTS market_data_1m
(
    symbol LowCardinality(String),
    timestamp DateTime64(3),
    open UInt64,
    high UInt64,
    low UInt64,
    close UInt64,
    volume Float64
)
ENGINE = ReplacingMergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (symbol, timestamp);

-- ==========================================
-- 2. METADATA INDEX (AggregatingMergeTree)
-- ==========================================
-- Stores summary statistics. 
-- Used by the Go API to quickly check date ranges without scanning the main table.
CREATE TABLE IF NOT EXISTS market_metadata
(
    symbol LowCardinality(String),
    min_time SimpleAggregateFunction(min, DateTime64(3)),
    max_time SimpleAggregateFunction(max, DateTime64(3)),
    row_count SimpleAggregateFunction(sum, UInt64)
)
ENGINE = AggregatingMergeTree()
ORDER BY symbol;

-- ==========================================
-- 3. METADATA TRIGGER (Materialized View)
-- ==========================================
-- Automatically updates 'market_metadata' whenever data is inserted into 'market_data_1m'.
CREATE MATERIALIZED VIEW IF NOT EXISTS market_metadata_mv TO market_metadata
AS SELECT
    symbol,
    min(timestamp) as min_time,
    max(timestamp) as max_time,
    count() as row_count
FROM market_data_1m
GROUP BY symbol;

-- ==========================================
-- 4. RESULTS TABLE (MergeTree)
-- ==========================================
-- Stores the output of C++ backtest jobs.
CREATE TABLE IF NOT EXISTS backtest_results
(
    job_id String,
    strategy_id String,
    symbol LowCardinality(String),
    start_time DateTime64(3),
    end_time DateTime64(3),
    total_pnl Int64,
    sharpe_ratio Float64,
    max_drawdown Float64,
    execution_time_ms UInt64,
    created_at DateTime DEFAULT now()
)
ENGINE = MergeTree()
ORDER BY (job_id, symbol);
