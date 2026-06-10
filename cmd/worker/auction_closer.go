package worker

import (
	"context"
	"log"
	"time"

	strading "maz/service/trading"
)

// RunAuctionCloser periodically settles expired auctions.
func RunAuctionCloser(ctx context.Context, trading *strading.Service, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			result, err := trading.CloseExpiredAuctions(ctx, now.UTC())
			if err != nil {
				log.Printf("auction closer: %v", err)
				continue
			}
			if result.Processed > 0 {
				log.Printf("auction closer: processed=%d settled=%d cancelled=%d",
					result.Processed, result.Settled, result.Cancelled)
			}
		}
	}
}
