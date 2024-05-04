package domain

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"github.com/totegamma/concurrent/core"
	"github.com/totegamma/concurrent/util"
)

var tracer = otel.Tracer("domain")

// Service is the domain service interface
type Handler interface {
	Get(c echo.Context) error
	List(c echo.Context) error
}

type handler struct {
	service core.DomainService
	config  util.Config
}

// NewHandler creates a new handler
func NewHandler(service core.DomainService, config util.Config) Handler {
	return &handler{service, config}
}

// Get returns a host by ID
func (h handler) Get(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "Domain.Handler.Get")
	defer span.End()

	id := c.Param("id")
	host, err := h.service.GetByFQDN(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Domain not found"})
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "content": host})

}

// List returns all hosts
func (h handler) List(c echo.Context) error {
	ctx, span := tracer.Start(c.Request().Context(), "Domain.Handler.List")
	defer span.End()

	hosts, err := h.service.List(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "content": hosts})
}
