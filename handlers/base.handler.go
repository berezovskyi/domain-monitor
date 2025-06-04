package handlers

import (
	"github.com/berezovskyi/domain-monitor/views/layout"
	"github.com/labstack/echo/v4"
)

type BaseHandler struct {
	IncludeConfiguration bool
}

func (bh *BaseHandler) HandlerShowBase(c echo.Context) error {
	return View(c, layout.Base(bh.IncludeConfiguration))
}
