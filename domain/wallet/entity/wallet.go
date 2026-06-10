package entity

import (
	"github.com/google/uuid"
)

// Wallet holds a guild's gold balance and active bid reservations.
type Wallet struct {
	GuildID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Balance  int64     `gorm:"not null;default:0;check:balance >= 0"`
	Reserved int64     `gorm:"not null;default:0;check:reserved >= 0"`

	Guild *Guild `gorm:"foreignKey:GuildID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Wallet) TableName() string { return "wallets" }

// AvailableBalance returns spendable gold (balance minus active reserves).
func (w *Wallet) AvailableBalance() int64 {
	return w.Balance - w.Reserved
}
