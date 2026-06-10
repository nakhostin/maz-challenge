package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
)

// Bid is an offer on a Legendary auction.
type Bid struct {
	ID                  uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AuctionID           uuid.UUID        `gorm:"type:uuid;not null;index:idx_bids_auction_created,priority:1"`
	GuildID             uuid.UUID        `gorm:"type:uuid;not null;index"`
	Amount              int64            `gorm:"not null;check:amount > 0"`
	Status              shared.BidStatus `gorm:"type:text;not null"`
	ReservationLedgerID *uuid.UUID       `gorm:"type:uuid;index"`
	IdempotencyKey      *string          `gorm:"type:text;uniqueIndex"`
	CreatedAt           time.Time        `gorm:"not null;autoCreateTime;index:idx_bids_auction_created,priority:2,sort:desc"`

	Auction *Auction `gorm:"foreignKey:AuctionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Bid) TableName() string { return "bids" }

func (b *Bid) BeforeCreate(_ *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// IsHighest reports whether this bid is the current winning offer.
func (b *Bid) IsHighest(auction *Auction) bool {
	if auction == nil || auction.HighestBidID == nil {
		return false
	}
	return b.ID == *auction.HighestBidID
}
