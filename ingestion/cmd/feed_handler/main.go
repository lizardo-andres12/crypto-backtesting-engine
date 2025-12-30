package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"cbe.com/internal/repository"
	"cbe.com/internal/service"
	"cbe.com/pkg/logger"
	"github.com/google/uuid"
)

const (
	LocalClickhouseInstance = "127.0.0.1:9000" // HTTP port
)

func main() {
	// TODO: create a hook registry for custom config
	// Run setup hooks
	logger.Setup()

	// TODO: upgrade to open telemetry
	// Create context
	runID := uuid.New().String()
	ctx := context.WithValue(context.Background(), logger.RunIDKey, runID)

	// TODO: study heap allocation implications for flag vars (String vs StringVar)
	// Create and parse flags
	symbol := flag.String("symbol", "", "Trading pair symbol (e.g. BTCUSDT)")
	limit := flag.Uint64("limit", 500, "Number of candles to fetch")
	flag.Parse()

	if *symbol == "" {
		slog.ErrorContext(ctx, "Symbol is required")
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Starting ingestion job", slog.String("symbol", *symbol), slog.Uint64("limit", *limit))

	// Initialize Components
	// TODO: display address of connection, which is determined via config file (slog.String(addr, config.addr)
	slog.InfoContext(ctx, "Attempting to connect to clickhouse instace", slog.String("addr", LocalClickhouseInstance))
	repo, err := repository.NewClickhouseRepository(ctx, LocalClickhouseInstance)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to ClickHouse", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()
	slog.InfoContext(ctx, "Connected to clickhouse instance")

	ingester := service.NewBinanceIngester()

	// Run ingestion
	slog.InfoContext(ctx, "Fetching candlestick data...",
		slog.String("symbol", *symbol),
		slog.Uint64("limit", *limit),
	)
	candlestickData, err := ingester.FetchAndParse(ctx, *symbol, *limit)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch data", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.InfoContext(ctx, "Fetched candlestick data", slog.Int("count", len(candlestickData)))

	// Batch insert
	slog.InfoContext(ctx, "Inserting results into clickhouse instace...")
	if err := repo.BatchInsert(ctx, candlestickData); err != nil {
		slog.ErrorContext(ctx, "Failed to insert data", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Ingestion completed successfully!")
}
