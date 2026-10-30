package domain

import "time"

type Crypto struct {
	Symbol       string
	Name         string
	CurrentPrice float64
	LastUpdated  time.Time
}

type PricePoint struct {
	Price     float64
	Timestamp time.Time
}

type Stats struct {
	MinPrice           float64
	MaxPrice           float64
	AvgPrice           float64
	PriceChange        float64
	PriceChangePercent float64
	RecordsCount       int
}
