package handlers

import (
	"github.com/berezovskyi/domain-monitor/views/dashboard"
	"github.com/labstack/echo/v4"
)

func HandlerRenderDashboard(c echo.Context) error {
	dashboard := dashboard.Dashboard()

	return View(c, dashboard)
}
