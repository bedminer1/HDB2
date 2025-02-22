package main

import (
	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

func main() {
	e := echo.New()
	h := initHandler()
	
	e.GET("/healthcheck", h.handleHealthCheck)
	e.GET("/records", h.handleGetRecords)
	e.GET("/monthly_stats", h.handleGetMonthlyStats)
	e.GET("/yearly_stats", h.handleGetYearlyStats)

	e.GET("/town_stats", h.handleGetTownBasedStats)
	e.GET("/town_predictions", h.handleGetTownBasedPredictions)
	
	e.GET("/linear_regression", h.handleGetLinearRegressionPrediction)
	e.GET("/polynomial_regression", h.handleGetPolynomialRegressionPrediction)
	e.GET("/holt_winters", h.handleGetHoltWinters)

	e.GET("/llm_analysis", h.handleGetLLMAnalysis)

	e.Logger.Fatal(e.Start(":4000"))
}