package trading

import (
	"context"
	"time"

	"github.com/google/uuid"

	"maz/domain/shared"
	tentity "maz/domain/trading/entity"
)

// CreateAuction opens an auction for a legendary item.
func (s *Service) CreateAuction(ctx context.Context, itemID, sellerGuildID uuid.UUID, startingPrice int64, now time.Time) (*AuctionView, error) {
	if startingPrice <= 0 {
		return nil, shared.ErrInvalidState
	}
	auction := &tentity.Auction{
		ID:            uuid.New(),
		ItemID:        itemID,
		SellerGuildID: sellerGuildID,
		Status:        shared.AuctionStatusActive,
		StartingPrice: startingPrice,
		EndsAt:        now.Add(s.auctionDuration),
		CreatedAt:     now,
	}
	if err := s.auctions(s.db).Create(ctx, auction); err != nil {
		return nil, err
	}
	view := toAuctionView(auction, false)
	return &view, nil
}
