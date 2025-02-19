package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PocketbaseHDBRecord struct {
	Block             int       `gorm:"column:block"`
	Created           string    `gorm:"column:created"` // You might want to use time.Time
	FlatModel         string    `gorm:"column:flat_model"`
	FlatType          string    `gorm:"column:flat_type"`
	FloorArea         int       `gorm:"column:floor_area"`
	LeastCommenceDate int       `gorm:"column:least_commence_date"`
	PricePerArea      float32   `gorm:"column:price_per_area"`
	ResalePrice       int       `gorm:"column:resale_price"`
	StoreyRange       string    `gorm:"column:storey_range"`
	StreetName        string    `gorm:"column:street_name"`
	Time              time.Time `gorm:"column:time"` // Time as string
	Town              string    `gorm:"column:town"`
	Updated           string    `gorm:"column:updated"` // You might want to use time.Time
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current Working Directory:", cwd)

	sourcePath := "./hdb.db"
	sourceDB, err := gorm.Open(sqlite.Open(sourcePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to source SQLite database: %v", err)
	}

	destPath := "./pb/pb_data/data.db"
	destDB, err := gorm.Open(sqlite.Open(destPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), 
	})
	if err != nil {
		log.Fatalf("failed to connect to destination SQLite database: %v", err)
	}

	var hdbRecords []models.HDBRecord
	err = sourceDB.Find(&hdbRecords).Error
	if err != nil {
		log.Fatalf("failed to fetch records from source SQLite: %v", err)
	}

	for _, hdbRecord := range hdbRecords {
		pbRecord := PocketbaseHDBRecord{
			Block:             hdbRecord.Block,
			FlatModel:         hdbRecord.FlatModel,
			FlatType:          hdbRecord.FlatType,
			FloorArea:         hdbRecord.FloorArea,
			LeastCommenceDate: hdbRecord.LeaseCommenceDate, // Corrected field name
			PricePerArea:      float32(hdbRecord.PricePerArea),
			ResalePrice:       hdbRecord.ResalePrice,
			StoreyRange:       hdbRecord.StoreyRange,
			StreetName:        hdbRecord.StreetName,
			Time:              hdbRecord.Time,
			Town:              hdbRecord.Town,
			// Created and Updated are handled by Pocketbase, so you can leave them empty
		}

		err = destDB.Table("hdb_records").Create(&pbRecord).Error // Use the new struct
		if err != nil {
			log.Printf("error creating record in destination DB: %v (Record ID: %s)", err, hdbRecord.ID)
		}
	}

	fmt.Println("Migration complete!")
}
