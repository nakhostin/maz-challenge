package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"maz/domain/shared"
	swallet "maz/service/wallet"
)

type WalletHandler struct {
	wallet *swallet.Service
}

func NewWalletHandler(wallet *swallet.Service) *WalletHandler {
	return &WalletHandler{wallet: wallet}
}

func (h *WalletHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeError(c, shared.ErrInvalidState)
	}
	view, err := h.wallet.GetWallet(c.Context(), id)
	if err != nil {
		return writeError(c, err)
	}
	return writeJSON(c, http.StatusOK, view)
}
