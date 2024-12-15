package main

import (
	"time"

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
	db.AutoMigrate()

	return &handler{DB: db}
}

func (h *handler) handleGetRecords(c echo.Context) error {
	records := []models.HDBRecord{}

	start := c.QueryParam("start")
	if start == "" {
		start = "2018-01"
	}
	end := c.QueryParam("end")
	if end == "" {
		end = "2021-01"
	}

	startTime, err := time.Parse("2006-01", start)
	if err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid 'start' date format, use YYYY-MM"})
	}
	endTime, err := time.Parse("2006-01", end)
	if err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid 'end' date format, use YYYY-MM"})
	}

	h.DB.Where("time BETWEEN ? AND ?", startTime, endTime).Find(&records)

	return c.JSON(200, echo.Map{
		"number of records": len(records),
		"records": records,
	})
}
