package marketplace

import (
	"time"

	"maz/domain/marketplace/entity"
)

// ItemView is the handler-facing item representation.
type ItemView struct {
	ID             string `json:"id"`
	SellerGuildID  string `json:"seller_guild_id"`
	OwnerGuildID   string `json:"owner_guild_id"`
	Name           string `json:"name"`
	ItemType       string `json:"item_type"`
	Status         string `json:"status"`
	ListPrice      int64  `json:"list_price"`
	ReferencePrice *int64 `json:"reference_price,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// RegisterItemCommand is input for POST /items.
type RegisterItemCommand struct {
	SellerGuildID  string
	Name           string
	ItemType       string
	ListPrice      int64
	IdempotencyKey string
	Now            time.Time
}

func toItemView(item *entity.Item, oracle map[string]int64) ItemView {
	var ref *int64
	if p, ok := oracle[item.Name]; ok {
		ref = &p
	}
	return ItemView{
		ID:             item.ID.String(),
		SellerGuildID:  item.SellerGuildID.String(),
		OwnerGuildID:   item.OwnerGuildID.String(),
		Name:           item.Name,
		ItemType:       string(item.ItemType),
		Status:         string(item.Status),
		ListPrice:      item.ListPrice,
		ReferencePrice: ref,
		CreatedAt:      item.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toItemViews(items []entity.Item, oracle map[string]int64) []ItemView {
	views := make([]ItemView, 0, len(items))
	for i := range items {
		views = append(views, toItemView(&items[i], oracle))
	}
	return views
}
