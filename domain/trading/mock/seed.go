package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	mmarket "maz/domain/marketplace/mock"
	"maz/domain/shared"
	"maz/domain/trading/entity"
)

var (
	AuctionSoulReaver      = uuid.MustParse("b3000003-0000-4000-8000-000000000001")
	AuctionEyeOfTheDragon  = uuid.MustParse("b3000003-0000-4000-8000-000000000002")
)

// Seed inserts active auctions for mock legendary items. Idempotent on primary keys.
func Seed(ctx context.Context, db *gorm.DB, now time.Time, duration time.Duration) error {
	if duration <= 0 {
		duration = 24 * time.Hour
	}

	auctions := []entity.Auction{
		{
			ID:            AuctionSoulReaver,
			ItemID:        mmarket.ItemSoulReaver,
			SellerGuildID: mmarket.GuildShadowSyndicate,
			Status:        shared.AuctionStatusActive,
			StartingPrice: 500_000,
			EndsAt:        now.Add(duration),
			CreatedAt:     now.Add(-6 * time.Hour),
		},
		{
			ID:            AuctionEyeOfTheDragon,
			ItemID:        mmarket.ItemEyeOfTheDragon,
			SellerGuildID: mmarket.GuildCrystalForge,
			Status:        shared.AuctionStatusActive,
			StartingPrice: 750_000,
			EndsAt:        now.Add(duration),
			CreatedAt:     now.Add(-3 * time.Hour),
		},
	}

	for i := range auctions {
		a := auctions[i]
		if err := db.WithContext(ctx).Save(&a).Error; err != nil {
			return fmt.Errorf("seed auction for item %s: %w", a.ItemID, err)
		}
	}
	return nil
}
