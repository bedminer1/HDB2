package calculation

import "github.com/bedminer1/hdb2/internal/models"

func CalculateTownTrends(records []models.TimeBasedRecord, monthsAhead int, dateBasis string) []models.TownPrediction {
	townRecords := make(map[string][]models.TimeBasedRecord)

	for _, record := range records {
		for _, town := range record.Towns { // Iterate through each town in the record
			townRecords[town] = append(townRecords[town], record)
		}
	}

	var predictions []models.TownPrediction

	// Iterate through each town's records
	for town, townRecord := range townRecords {
		// Call CalculatePolynomialRegression for each town
		predictedData, historicalData, model := CalculatePolynomialRegression(townRecord, 2, monthsAhead, dateBasis)

		// Calculate ExpectedROI based on predictedData and historicalData
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
			ExpectedROI:     expectedROI,
			PredictionModel: model,
		})
	}

	return predictions
}