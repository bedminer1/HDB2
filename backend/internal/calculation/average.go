package calculation

import (
	"sort"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
)

func MonthlyStats(records []models.HDBRecord) []models.TimeBasedRecord {
	return CalculateXlyStats("2006-01", records)
}

func YearlyStats(records []models.HDBRecord) []models.TimeBasedRecord {
	return CalculateXlyStats("2006", records)
}

func CalculateXlyStats(dateFormat string, records []models.HDBRecord) []models.TimeBasedRecord {
	XlyData := make(map[string][]models.HDBRecord)
	for _, record := range records {
		XKey := record.Time.Format(dateFormat) // Format as "YYYY-MM"
		XlyData[XKey] = append(XlyData[XKey], record)
	}

	var XlyRecords []models.TimeBasedRecord

	for XKey, records := range XlyData {
		var totalUnits int
		var totalResalePrice float64
		var totalPricePerArea float64

		townSet := make(map[string]struct{})
		flatTypeSet := make(map[string]struct{})

		for _, record := range records {
			totalUnits++
			totalResalePrice += float64(record.ResalePrice)
			totalPricePerArea += record.PricePerArea

			townSet[record.Town] = struct{}{}
			flatTypeSet[record.FlatType] = struct{}{}
		}

		towns := make([]string, 0, len(townSet))
		for town := range townSet {
			towns = append(towns, town)
		}

		flatTypes := make([]string, 0, len(flatTypeSet))
		for flatType := range flatTypeSet {
			flatTypes = append(flatTypes, flatType)
		}

		averageResalePrice := totalResalePrice / float64(totalUnits)
		averagePricePerArea := totalPricePerArea / float64(totalUnits)
		XTime, _ := time.Parse(dateFormat, XKey)

		XlyRecords = append(XlyRecords, models.TimeBasedRecord{
			Time:                XTime,
			Towns:               towns,
			FlatTypes:           flatTypes,
			NumberOfUnits:       totalUnits,
			AverageResalePrice:  averageResalePrice,
			AveragePricePerArea: averagePricePerArea,
		})
	}

	sortByTime(XlyRecords)

	return XlyRecords
}

func sortByTime(records []models.TimeBasedRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.Before(records[j].Time)
	})
}
