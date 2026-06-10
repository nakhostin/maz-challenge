package trading

import "time"

// BidView is returned after placing or querying a bid.
type BidView struct {
	ID        string `json:"id"`
	AuctionID string `json:"auction_id"`
	ItemID    string `json:"item_id,omitempty"`
	GuildID   string `json:"guild_id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// AuctionView is returned by auction endpoints.
type AuctionView struct {
	ID            string    `json:"id"`
	ItemID        string    `json:"item_id"`
	SellerGuildID string    `json:"seller_guild_id"`
	Status        string    `json:"status"`
	StartingPrice int64     `json:"starting_price"`
	HighestAmount *int64    `json:"highest_amount,omitempty"`
	HighestBidID  *string   `json:"highest_bid_id,omitempty"`
	EndsAt        string    `json:"ends_at"`
	CreatedAt     string    `json:"created_at"`
	Bids          []BidView `json:"bids,omitempty"`
}

// PlaceBidCommand is input for POST /items/{id}/bid.
type PlaceBidCommand struct {
	ItemID         string
	GuildID        string
	Amount         int64
	IdempotencyKey string
	Now            time.Time
}

// WithdrawBidCommand is input for DELETE /items/{id}/bid/{bid_id}.
type WithdrawBidCommand struct {
	ItemID  string
	BidID   string
	GuildID string
	Now     time.Time
}

// CloseExpiredResult summarizes auction closer worker output.
type CloseExpiredResult struct {
	Processed int
	Settled   int
	Cancelled int
}
