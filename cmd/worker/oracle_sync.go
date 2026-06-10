package worker

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	dmarket "maz/domain/marketplace"
	ment "maz/domain/marketplace/entity"
	"maz/pkg/oracle"
)

// RunOracleSync periodically refreshes cached reference prices from the oracle client.
func RunOracleSync(ctx context.Context, db *gorm.DB, client oracle.Client, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	repo := dmarket.NewOraclePriceRepository(db)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	syncOnce(ctx, repo, client)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			syncOnce(ctx, repo, client)
		}
	}
}

func syncOnce(ctx context.Context, repo dmarket.OraclePriceRepository, client oracle.Client) {
	prices, err := client.FetchPrices(ctx)
	if err != nil {
		log.Printf("oracle sync: fetch: %v", err)
		return
	}
	valid := oracle.SanitizePrices(prices)
	if len(valid) == 0 {
		log.Printf("oracle sync: no valid prices in response")
		return
	}
	now := time.Now().UTC()
	for name, price := range valid {
		if err := repo.Upsert(ctx, &ment.OraclePrice{
			ItemName:  name,
			BasePrice: price,
			UpdatedAt: now,
		}); err != nil {
			log.Printf("oracle sync: upsert %q: %v", name, err)
		}
	}
}
