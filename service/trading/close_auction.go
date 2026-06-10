package trading

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	tentity "maz/domain/trading/entity"
	swallet "maz/service/wallet"
)

// CloseExpiredAuctions settles or cancels auctions past their end time.
func (s *Service) CloseExpiredAuctions(ctx context.Context, now time.Time) (CloseExpiredResult, error) {
	var result CloseExpiredResult

	auctions, err := s.auctions(s.db).ListExpiredActiveForUpdate(ctx, now)
	if err != nil {
		return result, err
	}

	for i := range auctions {
		auction := &auctions[i]
		if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return s.closeOne(ctx, tx, auction.ID, now, &result)
		}); err != nil {
			return result, err
		}
	}
	return result, nil
}

func (s *Service) closeOne(ctx context.Context, tx *gorm.DB, auctionID uuid.UUID, now time.Time, result *CloseExpiredResult) error {
	auctionRepo := s.auctions(tx)
	itemRepo := s.items(tx)

	auction, err := auctionRepo.GetForUpdate(ctx, auctionID)
	if err != nil {
		return err
	}
	if auction.Status != shared.AuctionStatusActive {
		return nil
	}
	if now.Before(auction.EndsAt) {
		return nil
	}

	result.Processed++

	item, err := itemRepo.GetForUpdate(ctx, auction.ItemID)
	if err != nil {
		return err
	}

	if !auction.HasBids() {
		auction.Status = shared.AuctionStatusCancelled
		item.Status = shared.ItemStatusAvailable
		if err := auctionRepo.Save(ctx, auction); err != nil {
			return err
		}
		if err := itemRepo.Save(ctx, item); err != nil {
			return err
		}
		result.Cancelled++
		return nil
	}

	winner := findWinningBid(auction)
	if winner == nil {
		return shared.ErrInvalidState
	}
	winner.Status = shared.BidStatusWon
	auction.Status = shared.AuctionStatusSettled
	item.Status = shared.ItemStatusSold
	item.OwnerGuildID = winner.GuildID

	if err := s.wallet.Debit(ctx, tx, swallet.DebitParams{
		GuildID:       winner.GuildID.String(),
		Amount:        winner.Amount,
		ReferenceType: "auction_settle",
		ReferenceID:   auction.ID.String(),
		Now:           now,
	}); err != nil {
		return err
	}
	if err := s.wallet.Credit(ctx, tx, swallet.CreditParams{
		GuildID:       auction.SellerGuildID.String(),
		Amount:        winner.Amount,
		ReferenceType: "auction_sale",
		ReferenceID:   auction.ID.String(),
		Now:           now,
	}); err != nil {
		return err
	}
	if err := auctionRepo.Save(ctx, auction); err != nil {
		return err
	}
	if err := itemRepo.Save(ctx, item); err != nil {
		return err
	}
	result.Settled++
	return nil
}

func findWinningBid(auction *tentity.Auction) *tentity.Bid {
	if auction.HighestBidID == nil {
		return nil
	}
	for i := range auction.Bids {
		if auction.Bids[i].ID == *auction.HighestBidID {
			return &auction.Bids[i]
		}
	}
	return nil
}
