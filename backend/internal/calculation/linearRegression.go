package calculation

import (
	"fmt"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
	"gonum.org/v1/gonum/stat"
)

func CalculateLinearRegression(records []models.TimeBasedRecord, monthsAhead int, dateBasis string) ([]models.SimplifiedTimeBasedRecord, []models.SimplifiedTimeBasedRecord, string) {
	var x []float64
	var y []float64
	var historicalData []models.SimplifiedTimeBasedRecord

	startDate := records[0].Time

	for _, record := range records {
		monthsDiff := float64((record.Time.Year()-startDate.Year())*12 + int(record.Time.Month()-startDate.Month()))
		x = append(x, monthsDiff)
		y = append(y, record.AverageResalePrice)

		historicalData = append(historicalData, models.SimplifiedTimeBasedRecord{
			Date:               record.Time,
			AverageResalePrice: record.AverageResalePrice,
		})
	}

	intercept, slope := stat.LinearRegression(x, y, nil, false)
	lastDate := records[len(records)-1].Time
	lastX := x[len(x)-1]

	// Predict prices
	var predictions []models.SimplifiedTimeBasedRecord
	for i := 1; i <= monthsAhead; i++ {
		futureX := lastX + float64(i)
		predictedPrice := intercept + slope*futureX

		var futureDate time.Time
		switch dateBasis {
		case "yearly":
			futureDate = lastDate.AddDate(i, 0, 0)
		default:
			futureDate = lastDate.AddDate(0, i, 0)
		}

		predictions = append(predictions, models.SimplifiedTimeBasedRecord{
			Date:               futureDate,
			AverageResalePrice: predictedPrice,
		})
	}

	dateUnits := "Months"
	if dateBasis == "yearly" {
		dateUnits = "Years"
	}

	model := fmt.Sprintf("Price = %.2f + %.2f * %s\n", intercept, slope, dateUnits)

	return predictions, historicalData, model
}
