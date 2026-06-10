package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"maz/domain/shared"
)

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(c *fiber.Ctx, status int, payload any) error {
	return c.Status(status).JSON(payload)
}

func writeError(c *fiber.Ctx, err error) error {
	status, msg := mapError(err)
	return c.Status(status).JSON(errorBody{Error: msg})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, shared.ErrInsufficientFunds):
		return http.StatusPaymentRequired, err.Error()
	case errors.Is(err, shared.ErrDailyCapExceeded):
		return http.StatusPaymentRequired, err.Error()
	case errors.Is(err, shared.ErrSelfBid):
		return http.StatusConflict, err.Error()
	case errors.Is(err, shared.ErrSelfPurchase):
		return http.StatusConflict, err.Error()
	case errors.Is(err, shared.ErrBidTooLow):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, shared.ErrCannotWithdrawBid):
		return http.StatusConflict, err.Error()
	case errors.Is(err, shared.ErrAuctionNotActive):
		return http.StatusConflict, err.Error()
	case errors.Is(err, shared.ErrDuplicateLegendary):
		return http.StatusConflict, err.Error()
	case errors.Is(err, shared.ErrInvalidItemType):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, shared.ErrInvalidState):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
