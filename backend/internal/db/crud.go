package db

import (
	"fmt"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
	"gorm.io/gorm"
)

func Fetch(start, end string, towns []string, flatType string, db *gorm.DB) ([]models.HDBRecord, error) {
	records := []models.HDBRecord{}

	startTime, err := time.Parse("2006-01", start)
	if err != nil {
		return nil, fmt.Errorf("invalid 'start' date format, use YYYY-MM")
	}
	endTime, err := time.Parse("2006-01", end)
	if err != nil {
		return nil, fmt.Errorf("invalid 'end' date format, use YYYY-MM")
	}

	query := db.Where("time BETWEEN ? AND ?", startTime, endTime)

	if len(towns) > 0 {
		query = query.Where("town IN ?", towns)
	}
	if flatType != "" {
		query = query.Where("flat_type = ?", flatType)
	}

	// Fetch data with limit and offset
	result := query.Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}
