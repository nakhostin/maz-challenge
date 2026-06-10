package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"maz/cmd/middleware"
	"maz/pkg/fiberx"
	strading "maz/service/trading"
)

type BidHandler struct {
	trading *strading.Service
}

func NewBidHandler(trading *strading.Service) *BidHandler {
	return &BidHandler{trading: trading}
}

type placeBidRequest struct {
	Amount int64 `json:"amount"`
}

func (h *BidHandler) Place(c *fiber.Ctx) error {
	guildID, ok := middleware.GuildID(c)
	if !ok {
		return writeError(c, errMissingGuild)
	}

	body, err := fiberx.BodyParser[placeBidRequest](c)
	if err != nil {
		return writeError(c, err)
	}

	view, err := h.trading.PlaceBid(c.Context(), strading.PlaceBidCommand{
		ItemID:         c.Params("id"),
		GuildID:        guildID.String(),
		Amount:         body.Amount,
		IdempotencyKey: c.Get("Idempotency-Key"),
		Now:            time.Now().UTC(),
	})
	if err != nil {
		return writeError(c, err)
	}
	return writeJSON(c, http.StatusCreated, view)
}

func (h *BidHandler) Withdraw(c *fiber.Ctx) error {
	guildID, ok := middleware.GuildID(c)
	if !ok {
		return writeError(c, errMissingGuild)
	}

	err := h.trading.WithdrawBid(c.Context(), strading.WithdrawBidCommand{
		ItemID:  c.Params("id"),
		BidID:   c.Params("bid_id"),
		GuildID: guildID.String(),
		Now:     time.Now().UTC(),
	})
	if err != nil {
		return writeError(c, err)
	}
	return c.SendStatus(http.StatusNoContent)
}
