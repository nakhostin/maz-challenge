package marketplace

import (
	"time"

	"gorm.io/gorm"

	dmarket "maz/domain/marketplace"
	strading "maz/service/trading"
)

// Service orchestrates marketplace use cases.
type Service struct {
	db              *gorm.DB
	trading         *strading.Service
	auctionDuration time.Duration
}

func NewService(db *gorm.DB, trading *strading.Service, auctionDuration time.Duration) *Service {
	if auctionDuration <= 0 {
		auctionDuration = 24 * time.Hour
	}
	return &Service{
		db:              db,
		trading:         trading,
		auctionDuration: auctionDuration,
	}
}

func (s *Service) items(db *gorm.DB) dmarket.ItemRepository {
	return dmarket.NewItemRepository(db)
}

func (s *Service) oracle(db *gorm.DB) dmarket.OraclePriceRepository {
	return dmarket.NewOraclePriceRepository(db)
}
