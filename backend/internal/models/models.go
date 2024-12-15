package models

import (
	"time"

	"github.com/google/uuid"
)

type HDBRecord struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key"`
	Time              time.Time
	Town              string
	FlatType          string
	Block             int
	StreetName        string
	StoreyRange       string
	FloorArea         int
	FlatModel         string
	LeaseCommenceDate int
	ResalePrice       int
	PricePerArea      float64
}
