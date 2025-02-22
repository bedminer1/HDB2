package main

import (
	"net/http"
	"strconv"

	"github.com/bedminer1/hdb2/internal/calculation"
	"github.com/bedminer1/hdb2/internal/db"
	"github.com/bedminer1/hdb2/internal/llm"
	"github.com/bedminer1/hdb2/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func initHandler() *handler {
	db, err := gorm.Open(sqlite.Open("../../hdb.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.HDBRecord{})

	return &handler{DB: db}
}

func (h *handler) handleGetRecords(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, _, _, _ := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"records":           records,
	})
}

// ================= //
// DATE SORTED STATS //
// ================= //

func (h *handler) handleGetMonthlyStats(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, _, _, _ := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// SORT AND CALCULATE
	monthlyStats := calculation.MonthlyStats(records)

	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"number of months":  len(monthlyStats),
		"monthly_stats":     monthlyStats,
	})
}

func (h *handler) handleGetYearlyStats(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, _, _, _ := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// SORT AND CALCULATE
	yearlyStats := calculation.YearlyStats(records)

	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"number of years":   len(yearlyStats),
		"yearly_stats":      yearlyStats,
	})
}

// ================= //
// TOWN SORTED STATS //
// ================= //

func (h *handler) handleGetTownBasedStats(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, _, dateFormat, _ := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// SORT AND CALCULATE
	townBasedRecords := calculation.CalculateTownStats(records, dateFormat)
	return c.JSON(200, echo.Map{
		"records": townBasedRecords,
	})
}

func (h *handler) handleGetTownBasedPredictions(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, dateBasis, dateFormat, timeAhead := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// SORT AND CALCULATE
	townBasedPredictions := calculation.CalculateTownTrends(records, timeAhead, dateBasis, dateFormat)
	return c.JSON(200, echo.Map{
		"records": townBasedPredictions,
	})

}

// =============== //
// PREDICTED STATS //
// =============== //

func (h *handler) handleGetLinearRegressionPrediction(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, dateBasis, dateFormat, timeAhead := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// DO CALCULATIONS
	xlyStats := calculation.CalculateXlyStats(dateFormat, records)
	predictions, historicalData, model := calculation.CalculateLinearRegression(xlyStats, timeAhead, dateBasis)

	return c.JSON(200, echo.Map{
		"model":           model,
		"predictions":     predictions,
		"historical_data": historicalData,
	})
}

func (h *handler) handleGetPolynomialRegressionPrediction(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, dateBasis, dateFormat, timeAhead := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// DO CALCULATIONS
	xlyStats := calculation.CalculateXlyStats(dateFormat, records)
	predictions, historicalData, model := calculation.CalculatePolynomialRegression(xlyStats, 4, timeAhead, dateBasis)

	return c.JSON(200, echo.Map{
		"model":           model,
		"predictions":     predictions,
		"historical_data": historicalData,
	})
}

func (h *handler) handleGetHoltWinters(c echo.Context) error {
	// QUERY PARAMS
	start, end, towns, flatType, _, dateFormat, timeAhead := parseQueryParams(c)

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, towns, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	// DO CALCULATIONS
	params := calculation.HoltWintersParameters{
		Alpha:        0.2,
		Beta:         0.1,
		Gamma:        0.3,
		SeasonLength: 12,
	}
	xlyStats := calculation.CalculateXlyStats(dateFormat, records)
	predictions, historicalData, model := calculation.CalculateHoltWinters(xlyStats, timeAhead, params)

	return c.JSON(200, echo.Map{
		"model":           model,
		"predictions":     predictions,
		"historical_data": historicalData,
	})
}

func parseQueryParams(c echo.Context) (string, string, []string, string, string, string, int) {
	// QUERY PARAMS
	start := c.QueryParam("start")
	if start == "" {
		start = "2018-01"
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2021-01"
	}
	towns := c.QueryParams()["towns"]
	flatType := c.QueryParam("flatType")

	dateBasis := c.QueryParam("dateBasis")
	var dateFormat string
	switch dateBasis {
	case "yearly":
		dateFormat = "2006"
	default: // defaults to 'monthly'
		dateFormat = "2006-01"
	}

	timeAheadStr := c.QueryParam("timeAhead")
	timeAhead, _ := strconv.Atoi(timeAheadStr)
	if timeAhead == 0 {
		timeAhead = 5
	}

	return start, end, towns, flatType, dateBasis, dateFormat, timeAhead
}

// ============ //
// LLM ANALYSIS //
// ============ //

func (h *handler) handleGetLLMAnalysis(c echo.Context) error {
	town := c.QueryParam("town") 
	if town == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "town field cannot be empty",
		})
	}

	bullThesis, _ := llm.FetchBullAnalysis(town)
	bearThesis, _ := llm.FetchBearAnalysis(town)

	return c.JSON(200, echo.Map{
		"bull_thesis": bullThesis,
		"bear_thesis": bearThesis,
	})
}