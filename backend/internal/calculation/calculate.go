package calculation

import (
	"sort"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
)

func MonthlyStats(records []models.HDBRecord) []models.MonthlyRecord {
	monthlyData := make(map[string][]models.HDBRecord)
	for _, record := range records {
		monthKey := record.Time.Format("2006-01") // Format as "YYYY-MM"
		monthlyData[monthKey] = append(monthlyData[monthKey], record)
	}

	var monthlyRecords []models.MonthlyRecord

	for monthKey, records := range monthlyData {
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
		monthTime, _ := time.Parse("2006-01", monthKey)

		monthlyRecords = append(monthlyRecords, models.MonthlyRecord{
			Time:                monthTime,
			Towns:               towns,
			FlatTypes:           flatTypes,
			NumberOfUnits:       totalUnits,
			AverageResalePrice:  averageResalePrice,
			AveragePricePerArea: averagePricePerArea,
		})
	}

	sortByTime(monthlyRecords)

	return monthlyRecords
}

func YearlyStats(records []models.HDBRecord) []models.MonthlyRecord {
	yearlyData := make(map[string][]models.HDBRecord)
	for _, record := range records {
		yearKey := record.Time.Format("2006") // Format as "YYYY-MM"
		yearlyData[yearKey] = append(yearlyData[yearKey], record)
	}

	var yearlyRecords []models.MonthlyRecord

	for yearKey, records := range yearlyData {
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
		yearTime, _ := time.Parse("2006", yearKey)

		yearlyRecords = append(yearlyRecords, models.MonthlyRecord{
			Time:                yearTime,
			Towns:               towns,
			FlatTypes:           flatTypes,
			NumberOfUnits:       totalUnits,
			AverageResalePrice:  averageResalePrice,
			AveragePricePerArea: averagePricePerArea,
		})
	}

	sortByTime(yearlyRecords)

	return yearlyRecords
}

func sortByTime(records []models.MonthlyRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.Before(records[j].Time)
	})
}
