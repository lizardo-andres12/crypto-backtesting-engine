package repository

import (
	"context"
	"fmt"

	"cbe.com/internal/models"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const (
	DefaultDBName   = "crypto"
	DefaultUsername = "default"
	DefaultPassword = "default"
)

// ClickhouseRepository wraps a clichouse connection with I/O methods
type ClickhouseRepository struct {
	conn driver.Conn
}

// NewClickhouseRepository attempts to create a new ClickhouseRepository object and errors if any connection errors arise
func NewClickhouseRepository(ctx context.Context, addr string) (*ClickhouseRepository, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: DefaultDBName,
			Username: DefaultUsername,
			Password: DefaultPassword,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &ClickhouseRepository{conn: conn}, nil
}

// Close closes the clickhouse connection and should always be deferred after successful construcction of a ClickhouseRepository object
func (r *ClickhouseRepository) Close() error {
	return r.conn.Close()
}

// BatchInsert takes a slice of MarketDataModel and batch inserts them to the clickhouse instance
func (r *ClickhouseRepository) BatchInsert(ctx context.Context, data []models.MarketDataModel) error {
	if len(data) == 0 {
		return nil
	}

	batch, err := r.conn.PrepareBatch(ctx, "INSERT INTO market_data_1m")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, record := range data {
		if err := batch.AppendStruct(&record); err != nil {
			return err
		}
	}

	return batch.Send()
}
