package main

import (
	"github.com/bedminer1/hdb2/internal/calculation"
	"github.com/bedminer1/hdb2/internal/db"
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
	db.AutoMigrate()

	return &handler{DB: db}
}

func (h *handler) handleGetRecords(c echo.Context) error {
	// QUERY PARAMS
	start := c.QueryParam("start")
	if start == "" {
		start = "2018-01"
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2021-01"
	}
	town := c.QueryParam("town")
	flatType := c.QueryParam("flatType")

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, town, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"records":           records,
	})
}

func (h *handler) handleGetMonthlyStats(c echo.Context) error {
	// QUERY PARAMS
	start := c.QueryParam("start")
	if start == "" {
		start = "2018-01"
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2021-01"
	}
	town := c.QueryParam("town")
	flatType := c.QueryParam("flatType")

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, town, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	monthlyStats := calculation.MonthlyStats(records)
	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"number of months":  len(monthlyStats),
		"monthly_stats":     monthlyStats,
	})
}

func (h *handler) handleGetYearlyStats(c echo.Context) error {
	// QUERY PARAMS
	start := c.QueryParam("start")
	if start == "" {
		start = "2018-01"
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2021-01"
	}
	town := c.QueryParam("town")
	flatType := c.QueryParam("flatType")

	// CALL FETCH FROM DB PACKAGE
	records, err := db.Fetch(start, end, town, flatType, h.DB)
	if err != nil {
		return c.JSON(400, echo.Map{"error": err.Error()})
	}

	yearlyStats := calculation.YearlyStats(records)
	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"number of years":  len(yearlyStats),
		"yearly_stats":     yearlyStats,
	})
}
