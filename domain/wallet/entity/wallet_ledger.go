package entity

import (
	"maz/domain/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WalletLedger is an append-only audit log of wallet movements.
type WalletLedger struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GuildID        uuid.UUID              `gorm:"type:uuid;not null;index:idx_wallet_ledger_guild_created,priority:1"`
	EntryType      shared.LedgerEntryType `gorm:"type:text;not null"`
	Amount         int64                  `gorm:"not null;check:amount > 0"`
	ReferenceType  string                 `gorm:"type:text;not null"`
	ReferenceID    uuid.UUID              `gorm:"type:uuid;not null"`
	IdempotencyKey *string                `gorm:"type:text;uniqueIndex"`
	CreatedAt      time.Time              `gorm:"not null;autoCreateTime;index:idx_wallet_ledger_guild_created,priority:2,sort:desc"`

	Guild *Guild `gorm:"foreignKey:GuildID;references:ID;constraint:OnDelete:CASCADE"`
}

func (WalletLedger) TableName() string { return "wallet_ledger" }

func (e *WalletLedger) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
