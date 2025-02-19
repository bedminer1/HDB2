package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("../../hdb.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.HDBRecord{})

	err = db.Exec("DELETE FROM hdb_records").Error
	if err != nil {
		fmt.Printf("Error clearing table: %v\n", err)
		return
	}
	fmt.Println("Cleared existing records from the table.")

	csvFiles := []string{
		"../../csvData/file1.csv",
		"../../csvData/file2.csv",
		"../../csvData/file3.csv",
		"../../csvData/file4.csv",
		"../../csvData/file5.csv",
	}

	batchSize := 200

	for _, file := range csvFiles {
		priceIndex := 9
		if file == "../../csvData/file3.csv" || file == "../../csvData/file4.csv" {
			priceIndex = 10
		}
		records, err := readCSV(file, priceIndex)
		if err != nil {
			continue
		}

		// Insert records into the database
		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]

			result := db.Create(&batch)
			if result.Error != nil {
				fmt.Printf("Error inserting batch starting at %d: %v\n", i, result.Error)
				break
			}
		}

		fmt.Printf("Inserted %d records from %s\n", len(records), file)
	}
}

func readCSV(filePath string, priceIndex int) ([]models.HDBRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []models.HDBRecord
	for i, row := range rows {
		if i == 0 {
			continue
		}

		monthStr := row[0]
		parsedTime, _ := time.Parse("2006-01", monthStr)

		block, _ := strconv.Atoi((row[3]))
		floorArea, _ := strconv.Atoi(row[6])
		leaseCommenceDate, _ := strconv.Atoi(row[8])
		resalePrice, _ := strconv.Atoi(row[priceIndex])
		pricePerArea := float64(resalePrice) / float64(floorArea)

		record := models.HDBRecord{
			ID:                uuid.New(),
			Time:              parsedTime,
			Town:              row[1],
			FlatType:          row[2],
			Block:             block,
			StreetName:        row[4],
			StoreyRange:       row[5],
			FloorArea:         floorArea,
			FlatModel:         row[7],
			LeaseCommenceDate: leaseCommenceDate,
			ResalePrice:       resalePrice,
			PricePerArea:      pricePerArea,
		}

		records = append(records, record)
	}

	return records, nil
}
