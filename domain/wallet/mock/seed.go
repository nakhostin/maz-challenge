package mock

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	mmarket "maz/domain/marketplace/mock"
	"maz/domain/wallet/entity"
)

// Seed inserts dev guilds and wallets. Idempotent on primary keys.
func Seed(ctx context.Context, db *gorm.DB) error {
	guilds := []entity.Guild{
		{
			ID:               mmarket.GuildIronVanguard,
			Name:             "Iron Vanguard",
			DailyPurchaseCap: 500_000,
		},
		{
			ID:               mmarket.GuildShadowSyndicate,
			Name:             "Shadow Syndicate",
			DailyPurchaseCap: 500_000,
		},
		{
			ID:               mmarket.GuildCrystalForge,
			Name:             "Crystal Forge",
			DailyPurchaseCap: 500_000,
		},
	}

	for i := range guilds {
		g := guilds[i]
		if err := db.WithContext(ctx).Save(&g).Error; err != nil {
			return fmt.Errorf("seed guild %q: %w", g.Name, err)
		}
		wallet := entity.Wallet{
			GuildID:  g.ID,
			Balance:  1_000_000,
			Reserved: 0,
		}
		if err := db.WithContext(ctx).Save(&wallet).Error; err != nil {
			return fmt.Errorf("seed wallet for %q: %w", g.Name, err)
		}
	}
	return nil
}
