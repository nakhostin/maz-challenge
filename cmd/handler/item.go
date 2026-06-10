package handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"maz/cmd/middleware"
	"maz/domain/shared"
	"maz/pkg/fiberx"
	smarket "maz/service/marketplace"
)

type ItemHandler struct {
	marketplace *smarket.Service
}

func NewItemHandler(marketplace *smarket.Service) *ItemHandler {
	return &ItemHandler{marketplace: marketplace}
}

type registerItemRequest struct {
	Name      string `json:"name"`
	ItemType  string `json:"item_type"`
	ListPrice int64  `json:"list_price"`
}

func (h *ItemHandler) Register(c *fiber.Ctx) error {
	guildID, ok := middleware.GuildID(c)
	if !ok {
		return writeError(c, errMissingGuild)
	}

	body, err := fiberx.BodyParser[registerItemRequest](c)
	if err != nil {
		return writeError(c, err)
	}

	view, err := h.marketplace.RegisterItem(c.Context(), smarket.RegisterItemCommand{
		SellerGuildID:  guildID.String(),
		Name:           body.Name,
		ItemType:       body.ItemType,
		ListPrice:      body.ListPrice,
		IdempotencyKey: c.Get("Idempotency-Key"),
		Now:            time.Now().UTC(),
	})
	if err != nil {
		return writeError(c, err)
	}
	return writeJSON(c, http.StatusCreated, view)
}

func (h *ItemHandler) List(c *fiber.Ctx) error {
	items, err := h.marketplace.ListItems(c.Context())
	if err != nil {
		return writeError(c, err)
	}
	if items == nil {
		items = []smarket.ItemView{}
	}
	return writeJSON(c, http.StatusOK, items)
}

func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return writeError(c, shared.ErrInvalidState)
	}
	view, err := h.marketplace.GetItem(c.Context(), id)
	if err != nil {
		return writeError(c, err)
	}
	return writeJSON(c, http.StatusOK, view)
}
