package models

import "time"

// MarketDataModel is the universal data model with relevant fields for all time intervals
type MarketDataModel struct {
	Symbol    string `ch:"symbol"`
	Timestamp time.Time `ch:"timestamp"`
	Open      uint64 `ch:"open"`
	High      uint64 `ch:"high"`
	Low       uint64 `ch:"low"`
	Close     uint64 `ch:"close"`
	Volume    float64 `ch:"volume"`
}

