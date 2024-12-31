package calculation

import (
	"fmt"
	"math"
	"time"

	"github.com/bedminer1/hdb2/internal/models"
	"gonum.org/v1/gonum/mat"
)

func CalculatePolynomialRegression(records []models.TimeBasedRecord, degree, monthsAhead int, dateBasis string) ([]models.SimplifiedTimeBasedRecord, []models.SimplifiedTimeBasedRecord, string) {
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

	coefficients, err := polynomialFit(x, y, degree)
	if err != nil {
		fmt.Println("error fitting polynomial: ", err)
		return nil, nil, "Error in polynomial fitting"
	}

	lastDate := records[len(records)-1].Time
	lastX := x[len(x)-1]
	var predictions []models.SimplifiedTimeBasedRecord

	for i := 1; i <= monthsAhead; i++ {
		futureX := lastX + float64(i)
		predictedPrice := evaluatePolynomial(coefficients, futureX)
		if predictedPrice < 0 { predictedPrice = 0 }

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

	model := "Price = "
	for i, coef := range coefficients {
		if i > 0 {
			model += fmt.Sprintf(" + %.4f * X^%d", coef, i)
		} else {
			model += fmt.Sprintf("%.4f", coef)
		}
	}

	return predictions, historicalData, model
}

func polynomialFit(X, Y []float64, degree int) ([]float64, error) {
	vandermonde := mat.NewDense(len(X), degree+1, nil)
	for i := 0; i < len(X); i++ {
		for j := 0; j <= degree; j++ {
			vandermonde.Set(i, j, math.Pow(X[i], float64(j)))
		}
	}

	yVector := mat.NewVecDense(len(Y), Y)
	coefficients := mat.NewDense(degree+1, 1, nil)

	var qr mat.QR
	qr.Factorize(vandermonde)
	err := qr.SolveTo(coefficients, false, yVector)
	if err != nil {
		return nil, err
	}

	coef := make([]float64, degree+1)
	for i := 0; i <= degree; i++ {
		coef[i] = coefficients.At(i, 0)
	}
	return coef, nil
}

func evaluatePolynomial(coefficients []float64, x float64) float64 {
	var result float64
	for i, coef := range coefficients {
		result += coef * math.Pow(x, float64(i))
	}
	return result
}
