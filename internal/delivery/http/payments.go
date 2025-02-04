package http

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/service"
)

type PaymentsHandler struct {
	logger         *slog.Logger
	conf           config.Choco
	paymentService service.PaymentService
}

func NewPaymentsHandler(
	logger *slog.Logger,
	conf config.Choco,
	paymentService service.PaymentService,
) *PaymentsHandler {
	return &PaymentsHandler{
		logger:         logger,
		conf:           conf,
		paymentService: paymentService,
	}
}

func (ch *PaymentsHandler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api")

	api.GET("/get-company-customers", ch.getCompanyPayments())
}

func (ch *PaymentsHandler) getCompanyPayments() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(200, "Hello, World!")
	}
}
