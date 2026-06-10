package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"maz/domain/shared"
	strading "maz/service/trading"
)

type AuctionHandler struct {
	trading *strading.Service
}

func NewAuctionHandler(trading *strading.Service) *AuctionHandler {
	return &AuctionHandler{trading: trading}
}

func (h *AuctionHandler) List(c *fiber.Ctx) error {
	auctions, err := h.trading.ListActiveAuctions(c.Context())
	if err != nil {
		return writeError(c, err)
	}
	if auctions == nil {
		auctions = []strading.AuctionView{}
	}
	return writeJSON(c, http.StatusOK, auctions)
}

func (h *AuctionHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeError(c, shared.ErrInvalidState)
	}
	view, err := h.trading.GetAuction(c.Context(), id)
	if err != nil {
		return writeError(c, err)
	}
	return writeJSON(c, http.StatusOK, view)
}
