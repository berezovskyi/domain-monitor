package handlers

import (
	"errors"
	"net/http"

	"github.com/berezovskyi/domain-monitor/configuration"
	"github.com/labstack/echo/v4"
)

type ApiDomainService interface {
	CreateDomain(domain configuration.Domain) (int, error)
	GetDomain(fqdn string) (configuration.Domain, error)
	GetDomains() ([]configuration.Domain, error)
	UpdateDomain(domain configuration.Domain) error
	DeleteDomain(fqdn string) error
	Flush()
}

func NewApiDomainHandler(ds ApiDomainService) *ApiDomainHandler {
	return &ApiDomainHandler{
		DomainService: ds,
	}
}

type ApiDomainHandler struct {
	DomainService ApiDomainService
}

func (h *ApiDomainHandler) HandleDomainCreate(c echo.Context) error {
	var domain configuration.Domain
	if err := c.Bind(&domain); err != nil {
		return err
	}

	id, err := h.DomainService.CreateDomain(domain)
	if err != nil {
		return err
	}

	defer h.DomainService.Flush()
	return c.JSON(http.StatusCreated, id)
}

func (h *ApiDomainHandler) HandleDomainShow(c echo.Context) error {
	fqdn := c.Param("fqdn")
	if len(fqdn) == 0 {
		return errors.New("Invalid domain to fetch (FQDN required)")
	}

	domain, err := h.DomainService.GetDomain(fqdn)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, domain)
}

func (h *ApiDomainHandler) HandleDomainList(c echo.Context) error {
	domains, err := h.DomainService.GetDomains()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, domains)
}

func (h *ApiDomainHandler) HandleDomainUpdate(c echo.Context) error {
	var domain configuration.Domain
	if err := c.Bind(&domain); err != nil {
		return err
	}

	err := h.DomainService.UpdateDomain(domain)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ApiDomainHandler) HandleDomainDelete(c echo.Context) error {
	fqdn := c.Param("fqdn")
	if len(fqdn) == 0 {
		return errors.New("FQDN is required")
	}

	err := h.DomainService.DeleteDomain(fqdn)
	if err != nil {
		return err
	}

	defer h.DomainService.Flush()
	return c.NoContent(http.StatusNoContent)
}
