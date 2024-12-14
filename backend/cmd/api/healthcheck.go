package main

import "github.com/labstack/echo/v4"

func (h *handler) handleHealthCheck(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"status": "available",
		"system_info": map[string]string{
			"version": version,
		},
	})
}
