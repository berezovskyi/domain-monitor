package handlers

import (
	"log"

	"github.com/berezovskyi/domain-monitor/service"
	"github.com/berezovskyi/domain-monitor/views/configuration"
	"github.com/labstack/echo/v4"
)

func HandlerRenderConfiguration(c echo.Context) error {
	configuration := configuration.Configuration()

	return View(c, configuration)
}

type ConfigurationHandler struct {
	ConfigurationService service.ConfigurationService
}

func NewConfigurationHandler(cs service.ConfigurationService) *ConfigurationHandler {
	return &ConfigurationHandler{
		ConfigurationService: cs,
	}
}

// Support getting the value for a particular configuration settings and key.
//
// The possible keys are:
// - app
// - smtp
// - scheduler
// - alerts
//
// The possible keys for each section are represented by the keys in the ConfigurationFile struct.
func (h *ConfigurationHandler) GetSectionKey(c echo.Context) error {
	section := c.Param("section")
	key := c.Param("key")

	value, err := h.ConfigurationService.GetConfigurationValue(section, key)
	if err != nil {
		return err
	}

	return c.JSON(200, value)
}

// / Set the value for a particular configuration settings and key.
func (h *ConfigurationHandler) SetSectionKey(c echo.Context) error {
	section := c.Param("section")
	key := c.Param("key")

	value := c.FormValue("value")

	err := h.ConfigurationService.SetConfigurationValue(section, key, value)
	if err != nil {
		log.Printf("🚨 Error setting configuration value: %s", err.Error())
		return err
	}

	return c.NoContent(201)
}

// Render the domain configuration page.
func (h *ConfigurationHandler) RenderDomainConfiguration(c echo.Context) error {
	return View(c, configuration.DomainTab())
}

// Render the app configuration page.
func (h *ConfigurationHandler) RenderAppConfiguration(c echo.Context) error {
	return View(c, configuration.AppTab(h.ConfigurationService.GetAppConfiguration()))
}

// Render the smtp configuration page.
func (h *ConfigurationHandler) RenderSmtpConfiguration(c echo.Context) error {
	return View(c, configuration.SmtpTab(h.ConfigurationService.GetSMTPConfiguration()))
}

// Render the scheduler configuration page.
func (h *ConfigurationHandler) RenderSchedulerConfiguration(c echo.Context) error {
	return View(c, configuration.SchedulerTab(h.ConfigurationService.GetSchedulerConfiguration()))
}

// Render the alerts configuration page.
func (h *ConfigurationHandler) RenderAlertsConfiguration(c echo.Context) error {
	return View(c, configuration.AlertsTab(h.ConfigurationService.GetAlertsConfiguration()))
}
