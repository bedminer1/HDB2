package calculation

import (
	"fmt"

	"github.com/bedminer1/hdb2/internal/models"
)

type HoltWintersParameters struct {
	Alpha        float64
	Beta         float64
	Gamma        float64
	SeasonLength int
}

func CalculateHoltWinters(records []models.TimeBasedRecord, monthsAhead int, params HoltWintersParameters) ([]models.SimplifiedTimeBasedRecord, []models.SimplifiedTimeBasedRecord, string) {
	if len(records) < params.SeasonLength+1 {
		panic("Not enough data points for given season length")
	}

	var y []float64
	var historicalData []models.SimplifiedTimeBasedRecord
	for _, record := range records {
		y = append(y, record.AverageResalePrice)
		historicalData = append(historicalData, models.SimplifiedTimeBasedRecord{
			Date:               record.Time,
			AverageResalePrice: record.AverageResalePrice,
		})
	}

	L := make([]float64, len(y))
	T := make([]float64, len(y))
	S := make([]float64, len(y)+params.SeasonLength)

	L[0] = y[0]
	T[0] = y[1] - y[0]
	for i := 0; i < params.SeasonLength; i++ {
		S[i] = y[i] / L[0]
	}

	// Holt-Winters smoothing
	for t := 1; t < len(y); t++ {
		seasonIdx := (t - 1) % params.SeasonLength
		L[t] = params.Alpha*(y[t]/S[seasonIdx]) + (1-params.Alpha)*(L[t-1]+T[t-1])
		T[t] = params.Beta*(L[t]-L[t-1]) + (1-params.Beta)*T[t-1]
		S[t+params.SeasonLength] = params.Gamma*(y[t]/L[t]) + (1-params.Gamma)*S[seasonIdx]
	}

	// Predictions
	lastIdx := len(y) - 1
	lastDate := records[lastIdx].Time
	var predictions []models.SimplifiedTimeBasedRecord
	for h := 1; h <= monthsAhead; h++ {
		seasonIdx := (lastIdx + h) % params.SeasonLength
		futureValue := (L[lastIdx] + float64(h)*T[lastIdx]) * S[seasonIdx]

		predictions = append(predictions, models.SimplifiedTimeBasedRecord{
			Date:               lastDate.AddDate(0, h, 0),
			AverageResalePrice: futureValue,
		})
	}

	// Model representation
	model := fmt.Sprintf("Holt-Winters Model: α=%.2f, β=%.2f, γ=%.2f, SeasonLength=%d", params.Alpha, params.Beta, params.Gamma, params.SeasonLength)

	return predictions, historicalData, model
}
