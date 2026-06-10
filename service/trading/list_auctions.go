package trading

import (
	"context"

	"github.com/google/uuid"
)

// ListActiveAuctions returns all active auctions.
func (s *Service) ListActiveAuctions(ctx context.Context) ([]AuctionView, error) {
	auctions, err := s.auctions(s.db).ListActive(ctx)
	if err != nil {
		return nil, err
	}
	views := make([]AuctionView, 0, len(auctions))
	for i := range auctions {
		views = append(views, toAuctionView(&auctions[i], false))
	}
	return views, nil
}

// GetAuction returns a single auction with bids.
func (s *Service) GetAuction(ctx context.Context, auctionID uuid.UUID) (*AuctionView, error) {
	auction, err := s.auctions(s.db).GetByID(ctx, auctionID)
	if err != nil {
		return nil, err
	}
	view := toAuctionView(auction, true)
	return &view, nil
}

// GetActiveAuctionByItem returns the active auction for an item.
func (s *Service) GetActiveAuctionByItem(ctx context.Context, itemID uuid.UUID) (*AuctionView, error) {
	auction, err := s.auctions(s.db).GetActiveByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	view := toAuctionView(auction, true)
	return &view, nil
}
