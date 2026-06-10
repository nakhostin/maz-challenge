package entity

import (
	"maz/domain/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Item is a marketplace good (Common, Rare, or Legendary).
type Item struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SellerGuildID uuid.UUID         `gorm:"type:uuid;not null;index"`
	OwnerGuildID  uuid.UUID         `gorm:"type:uuid;not null;index"`
	Name          string            `gorm:"type:text;not null"`
	ItemType      shared.ItemType   `gorm:"column:item_type;type:text;not null"`
	Status        shared.ItemStatus `gorm:"type:text;not null"`
	ListPrice      int64             `gorm:"not null;check:list_price > 0"`
	IdempotencyKey *string           `gorm:"type:text;uniqueIndex"`
	CreatedAt      time.Time         `gorm:"not null;autoCreateTime"`
}

func (Item) TableName() string { return "items" }

func (i *Item) BeforeCreate(_ *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// IsLegendary reports whether the item is a one-of-a-kind Legendary good.
func (i *Item) IsLegendary() bool {
	return i.ItemType == shared.ItemTypeLegendary
}
