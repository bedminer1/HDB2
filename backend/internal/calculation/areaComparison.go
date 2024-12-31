package calculation

import (
	"sort"

	"github.com/bedminer1/hdb2/internal/models"
)

func CalculateTownStats(records []models.HDBRecord, dateFormat string) []models.TownBasedRecord {
	townGroupedRecords := make(map[string][]models.HDBRecord)

	for _, record := range records {
		townGroupedRecords[record.Town] = append(townGroupedRecords[record.Town], record)
	}

	var townBasedRecords []models.TownBasedRecord

	for town, townRecords := range townGroupedRecords {
		timeBasedRecords := CalculateXlyStats(dateFormat, townRecords)

		for i := range timeBasedRecords {
			timeBasedRecords[i].Towns = nil
			timeBasedRecords[i].FlatTypes = nil
		}

		townRecord := models.TownBasedRecord{
			Town:             town,
			TimeBasedRecords: timeBasedRecords,
		}

		townBasedRecords = append(townBasedRecords, townRecord)
	}

	return townBasedRecords
}

func CalculateTownTrends(records []models.HDBRecord, monthsAhead int, dateBasis, dateFormat string) []models.TownPrediction {
	townGroupedRecords := make(map[string][]models.HDBRecord)

	for _, record := range records {
		townGroupedRecords[record.Town] = append(townGroupedRecords[record.Town], record)
	}

	var predictions []models.TownPrediction

	// Iterate through each town's records
	for town, townRecords := range townGroupedRecords {
		timeBasedRecords := CalculateXlyStats(dateFormat, townRecords)
		predictedData, historicalData, model := CalculatePolynomialRegression(timeBasedRecords, 4, monthsAhead, dateBasis)

		var expectedROI float64
		if len(historicalData) > 0 && len(predictedData) > 0 {
			initialPrice := historicalData[len(historicalData)-1].AverageResalePrice
			finalPredictedPrice := predictedData[len(predictedData)-1].AverageResalePrice
			if initialPrice > 0 {
				expectedROI = (finalPredictedPrice - initialPrice) / initialPrice * 100 // ROI in percentage
			}
		}

		// Add the result to predictions
		predictions = append(predictions, models.TownPrediction{
			Town:            town,
			HistoricalData:  historicalData,
			PredictedData:   predictedData,
			MostRecentPrice: historicalData[len(historicalData)-1].AverageResalePrice,
			FinalPredictedPrice: predictedData[len(predictedData)-1].AverageResalePrice,
			ExpectedROI:     expectedROI,
			PredictionModel: model,
		})
	}

	sortPredictionsByROI(predictions)

	return predictions
}

func sortPredictionsByROI(predictions []models.TownPrediction) {
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].ExpectedROI > predictions[j].ExpectedROI
	})
}
