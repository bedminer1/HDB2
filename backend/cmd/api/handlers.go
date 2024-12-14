package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func initHandler() *handler {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate()

	return &handler{DB: db}
}

func (h *handler) healthCheckHandler(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"status": "available",
		"system_info": map[string]string{
			"version":     version,
		},
	})
}