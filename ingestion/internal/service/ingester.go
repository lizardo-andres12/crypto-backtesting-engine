package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"cbe.com/internal/models"
)

const (
	ScaleFactor = 100_000_000 // 10^8 (Satoshi precision)
	BinanceEndpoint  = "https://api.binance.com/api/v3/klines"
)

// BinanceIngester handles data fetching and parsing
type BinanceIngester struct {
	client *http.Client
}

// NewBinanceIngester returns a pointer to a new BinanceIngester object
func NewBinanceIngester() *BinanceIngester {
	return &BinanceIngester{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchAndParse downloads candles and converts them to the internal model
func (b *BinanceIngester) FetchAndParse(ctx context.Context, symbol string, limit uint64) ([]models.MarketDataModel, error) {
	url := fmt.Sprintf("%s?symbol=%s&interval=1m&limit=%d", BinanceEndpoint, symbol, limit)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var rawData [][]any
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return parseCandles(symbol, rawData)
}

func parseCandles(symbol string, rawData [][]any) ([]models.MarketDataModel, error) {
	result := make([]models.MarketDataModel, 0, len(rawData))

	for _, candle := range rawData {
		// Binance Format: [OpenTime, Open, High, Low, Close, Volume, ...]
		if len(candle) < 6 {
			continue // Skip malformed rows
		}

		tsMs, ok := candle[0].(float64)
		if !ok { return nil, fmt.Errorf("invalid timestamp format") }
		
		openStr, ok1 := candle[1].(string)
		highStr, ok2 := candle[2].(string)
		lowStr, ok3 := candle[3].(string)
		closeStr, ok4 := candle[4].(string)
		volStr, ok5 := candle[5].(string)

		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
			return nil, fmt.Errorf("invalid price/volume format (expected string)")
		}

		open, err := parsePrice(openStr)
		if err != nil { return nil, err }

		high, err := parsePrice(highStr)
		if err != nil { return nil, err }

		low, err := parsePrice(lowStr)
		if err != nil { return nil, err }

		closePrice, err := parsePrice(closeStr)
		if err != nil { return nil, err }

		volume, err := strconv.ParseFloat(volStr, 64)
		if err != nil { return nil, err }

		result = append(result, models.MarketDataModel{
			Symbol:    symbol,
			Timestamp: time.UnixMilli(int64(tsMs)),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}
	return result, nil
}

func parsePrice(priceStr string) (uint64, error) {
	valFloat, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}
	scaled := math.Round(valFloat * ScaleFactor)
	return uint64(scaled), nil
}

