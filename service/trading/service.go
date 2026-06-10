package trading

import (
	"time"

	"gorm.io/gorm"

	dmarket "maz/domain/marketplace"
	dtrading "maz/domain/trading"
	swallet "maz/service/wallet"
)

// Service orchestrates trading use cases.
type Service struct {
	db              *gorm.DB
	wallet          *swallet.Service
	auctionDuration time.Duration
}

func NewService(db *gorm.DB, wallet *swallet.Service, auctionDuration time.Duration) *Service {
	if auctionDuration <= 0 {
		auctionDuration = 24 * time.Hour
	}
	return &Service{
		db:              db,
		wallet:          wallet,
		auctionDuration: auctionDuration,
	}
}

func (s *Service) auctions(tx *gorm.DB) dtrading.AuctionRepository {
	return dtrading.NewAuctionRepository(tx)
}

func (s *Service) bids(tx *gorm.DB) dtrading.BidRepository {
	return dtrading.NewBidRepository(tx)
}

func (s *Service) items(tx *gorm.DB) dmarket.ItemRepository {
	return dmarket.NewItemRepository(tx)
}

func (s *Service) AuctionDuration() time.Duration {
	return s.auctionDuration
}
