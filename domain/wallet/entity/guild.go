package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Guild is a trading party in the Dragon Market.
type Guild struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name             string    `gorm:"type:text;not null;uniqueIndex"`
	DailyPurchaseCap int64     `gorm:"not null;check:daily_purchase_cap > 0"`
	CreatedAt        time.Time `gorm:"not null;autoCreateTime"`

	Wallet       *Wallet        `gorm:"foreignKey:GuildID;references:ID;constraint:OnDelete:CASCADE"`
	Ledger       []WalletLedger `gorm:"foreignKey:GuildID;references:ID"`
	DailySpends  []DailySpend   `gorm:"foreignKey:GuildID;references:ID"`
}

func (Guild) TableName() string { return "guilds" }

func (g *Guild) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
