package trading

import (
	"time"

	tentity "maz/domain/trading/entity"
)

// MinRequiredBidForTest exposes bid floor logic for unit tests.
func MinRequiredBidForTest(auction *tentity.Auction) int64 {
	return minRequiredBid(auction)
}

// ShouldExtendForTest exposes anti-snipe window logic for unit tests.
func ShouldExtendForTest(endsAt, now time.Time) bool {
	return shouldExtend(endsAt, now)
}
