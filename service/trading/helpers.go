package trading

import (
	"time"

	"github.com/google/uuid"

	"maz/domain/shared"
	tentity "maz/domain/trading/entity"
)

func minRequiredBid(auction *tentity.Auction) int64 {
	if !auction.HasBids() {
		return auction.StartingPrice
	}
	inc := *auction.HighestAmount * shared.MinBidIncreasePercent / 100
	if inc < 1 {
		inc = 1
	}
	return *auction.HighestAmount + inc
}

func shouldExtend(endsAt, now time.Time) bool {
	return endsAt.Sub(now) <= time.Duration(shared.AntiSnipeWindow)*time.Minute
}

func extendEndsAt(now time.Time) time.Time {
	return now.Add(time.Duration(shared.AntiSnipeExtend) * time.Minute)
}

func markHighestOutbid(auction *tentity.Auction) *tentity.Bid {
	if auction.HighestBidID == nil {
		return nil
	}
	for i := range auction.Bids {
		if auction.Bids[i].ID == *auction.HighestBidID && auction.Bids[i].Status == shared.BidStatusActive {
			auction.Bids[i].Status = shared.BidStatusOutbid
			return &auction.Bids[i]
		}
	}
	return nil
}

func findBid(auction *tentity.Auction, bidID uuid.UUID) (*tentity.Bid, error) {
	for i := range auction.Bids {
		if auction.Bids[i].ID == bidID {
			return &auction.Bids[i], nil
		}
	}
	return nil, shared.ErrNotFound
}

func toBidView(b *tentity.Bid, itemID uuid.UUID) BidView {
	v := BidView{
		ID:        b.ID.String(),
		AuctionID: b.AuctionID.String(),
		GuildID:   b.GuildID.String(),
		Amount:    b.Amount,
		Status:    string(b.Status),
		CreatedAt: b.CreatedAt.UTC().Format(time.RFC3339),
	}
	if itemID != uuid.Nil {
		v.ItemID = itemID.String()
	}
	return v
}

func toAuctionView(a *tentity.Auction, includeBids bool) AuctionView {
	v := AuctionView{
		ID:            a.ID.String(),
		ItemID:        a.ItemID.String(),
		SellerGuildID: a.SellerGuildID.String(),
		Status:        string(a.Status),
		StartingPrice: a.StartingPrice,
		HighestAmount: a.HighestAmount,
		EndsAt:        a.EndsAt.UTC().Format(time.RFC3339),
		CreatedAt:     a.CreatedAt.UTC().Format(time.RFC3339),
	}
	if a.HighestBidID != nil {
		s := a.HighestBidID.String()
		v.HighestBidID = &s
	}
	if includeBids {
		v.Bids = make([]BidView, 0, len(a.Bids))
		for i := range a.Bids {
			v.Bids = append(v.Bids, toBidView(&a.Bids[i], a.ItemID))
		}
	}
	return v
}
