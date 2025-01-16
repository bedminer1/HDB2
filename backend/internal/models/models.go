package models

import (
	"time"

	"github.com/google/uuid"
)

type HDBRecord struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Time              time.Time `gorm:"index" json:"time"`
	Town              string    `gorm:"index" json:"town"`
	FlatType          string    `gorm:"index" json:"flat_type"`
	Block             int       `json:"block"`
	StreetName        string    `json:"street_name"`
	StoreyRange       string    `json:"storey_range"`
	FloorArea         int       `json:"floor_area"`
	FlatModel         string    `json:"flat_model"`
	LeaseCommenceDate int       `json:"lease_commence_date"`
	ResalePrice       int       `json:"resale_price"`
	PricePerArea      float64   `json:"price_per_area"`
}

type TimeBasedRecord struct {
	Time                time.Time `json:"time"`
	Towns               []string  `json:"towns,omitempty"`
	FlatTypes           []string  `json:"flatTypes,omitempty"`
	NumberOfUnits       int       `json:"numberOfUnits"`
	AverageResalePrice  float64   `json:"averageResalePrice"`
	AveragePricePerArea float64   `json:"averagePricePerArea"`
}

type SimplifiedTimeBasedRecord struct {
	Date               time.Time `json:"date"`
	AverageResalePrice float64   `json:"average_resale_price"`
}

type TownBasedRecord struct {
	Town             string            `json:"town"`
	TimeBasedRecords []TimeBasedRecord `json:"records"`
}

type TownPrediction struct {
	Town                  string
	HistoricalData        []SimplifiedTimeBasedRecord `json:"-"`
	PredictedData         []SimplifiedTimeBasedRecord `json:"-"`
	MostRecentPrice       float64
	FinalPredictedPrice   float64
	ExpectedROI           float64
	PredictionModel       string
}
