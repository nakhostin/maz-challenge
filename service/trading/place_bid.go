package trading

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	tentity "maz/domain/trading/entity"
	swallet "maz/service/wallet"
)

// PlaceBid places a bid on the active auction for an item.
func (s *Service) PlaceBid(ctx context.Context, cmd PlaceBidCommand) (*BidView, error) {
	itemID, err := uuid.Parse(cmd.ItemID)
	if err != nil {
		return nil, shared.ErrInvalidState
	}
	guildID, err := uuid.Parse(cmd.GuildID)
	if err != nil {
		return nil, shared.ErrInvalidState
	}
	if cmd.Amount <= 0 {
		return nil, shared.ErrInvalidState
	}

	var view *BidView
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		bidRepo := s.bids(tx)
		if cmd.IdempotencyKey != "" {
			if existing, err := bidRepo.FindByIdempotencyKey(ctx, cmd.IdempotencyKey); err == nil {
				v := toBidView(existing, itemID)
				view = &v
				return nil
			} else if err != shared.ErrNotFound {
				return err
			}
		}

		auctionRepo := s.auctions(tx)
		auction, err := auctionRepo.GetForUpdateByItemID(ctx, itemID)
		if err != nil {
			return err
		}
		if !auction.IsActive(cmd.Now) {
			return shared.ErrAuctionNotActive
		}
		if guildID == auction.SellerGuildID {
			return shared.ErrSelfBid
		}
		if cmd.Amount < minRequiredBid(auction) {
			return shared.ErrBidTooLow
		}

		outbid := markHighestOutbid(auction)

		newBid := tentity.Bid{
			ID:        uuid.New(),
			AuctionID: auction.ID,
			GuildID:   guildID,
			Amount:    cmd.Amount,
			Status:    shared.BidStatusActive,
			CreatedAt: cmd.Now,
		}
		if cmd.IdempotencyKey != "" {
			key := cmd.IdempotencyKey
			newBid.IdempotencyKey = &key
		}

		ledgerID, err := s.wallet.Reserve(ctx, tx, swallet.ReserveParams{
			GuildID:       guildID.String(),
			Amount:        cmd.Amount,
			ReferenceType: "bid",
			ReferenceID:   newBid.ID.String(),
			Now:           cmd.Now,
		})
		if err != nil {
			return err
		}
		newBid.ReservationLedgerID = &ledgerID

		if outbid != nil {
			if err := s.wallet.Release(ctx, tx, swallet.ReleaseParams{
				GuildID:       outbid.GuildID.String(),
				Amount:        outbid.Amount,
				ReferenceType: "bid_outbid",
				ReferenceID:   newBid.ID.String(),
				Now:           cmd.Now,
			}); err != nil {
				return err
			}
		}

		auction.Bids = append(auction.Bids, newBid)
		auction.HighestBidID = &newBid.ID
		auction.HighestAmount = &cmd.Amount
		if shouldExtend(auction.EndsAt, cmd.Now) {
			auction.EndsAt = extendEndsAt(cmd.Now)
		}

		if err := auctionRepo.Save(ctx, auction); err != nil {
			return err
		}
		v := toBidView(&newBid, itemID)
		view = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return view, nil
}
