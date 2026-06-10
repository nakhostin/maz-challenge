package entity

import (
	"time"

	"github.com/google/uuid"
)

// DailySpend tracks per-guild committed gold for a calendar day (daily cap enforcement).
type DailySpend struct {
	GuildID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	SpendDate time.Time `gorm:"type:date;primaryKey"`
	Amount    int64     `gorm:"not null;default:0;check:amount >= 0"`

	Guild *Guild `gorm:"foreignKey:GuildID;references:ID;constraint:OnDelete:CASCADE"`
}

func (DailySpend) TableName() string { return "daily_spend" }
