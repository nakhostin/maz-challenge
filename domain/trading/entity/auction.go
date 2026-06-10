package entity

import (
	"maz/domain/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Auction is a time-boxed sale for a Legendary item.
type Auction struct {
	ID            uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ItemID        uuid.UUID            `gorm:"type:uuid;not null;index"`
	SellerGuildID uuid.UUID            `gorm:"type:uuid;not null;index"`
	Status        shared.AuctionStatus `gorm:"type:text;not null;index"`
	StartingPrice int64                `gorm:"not null;check:starting_price > 0"`
	HighestBidID  *uuid.UUID           `gorm:"type:uuid;index"`
	HighestAmount *int64               `gorm:""`
	EndsAt        time.Time            `gorm:"not null;index"`
	CreatedAt     time.Time            `gorm:"not null;autoCreateTime"`

	Bids []Bid `gorm:"foreignKey:AuctionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Auction) TableName() string { return "auctions" }

func (a *Auction) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// IsActive reports whether the auction accepts bids at the given time.
func (a *Auction) IsActive(now time.Time) bool {
	return a.Status == shared.AuctionStatusActive && now.Before(a.EndsAt)
}

// HasBids reports whether any bid is currently winning.
func (a *Auction) HasBids() bool {
	return a.HighestBidID != nil && a.HighestAmount != nil
}
