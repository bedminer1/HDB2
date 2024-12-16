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
	e.GET("/linear_regression", h.handleGetLinearRegressionPrediction)
	e.GET("/town_stats", h.handleGetTownBasedStats)

	e.Logger.Fatal(e.Start(":4000"))
}
