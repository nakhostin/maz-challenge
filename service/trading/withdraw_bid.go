package trading

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	swallet "maz/service/wallet"
)

// WithdrawBid cancels a non-winning active bid.
func (s *Service) WithdrawBid(ctx context.Context, cmd WithdrawBidCommand) error {
	itemID, err := uuid.Parse(cmd.ItemID)
	if err != nil {
		return shared.ErrInvalidState
	}
	bidID, err := uuid.Parse(cmd.BidID)
	if err != nil {
		return shared.ErrInvalidState
	}
	guildID, err := uuid.Parse(cmd.GuildID)
	if err != nil {
		return shared.ErrInvalidState
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		auctionRepo := s.auctions(tx)
		auction, err := auctionRepo.GetForUpdateByItemID(ctx, itemID)
		if err != nil {
			return err
		}
		if !auction.IsActive(cmd.Now) {
			return shared.ErrAuctionNotActive
		}

		bid, err := findBid(auction, bidID)
		if err != nil {
			return err
		}
		if bid.GuildID != guildID {
			return shared.ErrNotFound
		}
		if auction.HighestBidID != nil && bid.ID == *auction.HighestBidID {
			return shared.ErrCannotWithdrawBid
		}
		if bid.Status != shared.BidStatusActive {
			return shared.ErrInvalidState
		}

		if err := s.wallet.Release(ctx, tx, swallet.ReleaseParams{
			GuildID:       bid.GuildID.String(),
			Amount:        bid.Amount,
			ReferenceType: "bid_withdraw",
			ReferenceID:   bid.ID.String(),
			Now:           cmd.Now,
		}); err != nil {
			return err
		}

		bid.Status = shared.BidStatusWithdrawn
		return auctionRepo.Save(ctx, auction)
	})
}
