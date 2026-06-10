package migration

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	ment "maz/domain/marketplace/entity"
	tent "maz/domain/trading/entity"
	went "maz/domain/wallet/entity"
)

// AutoMigrate creates or updates schema for all domain entities.
func AutoMigrate(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).AutoMigrate(
		&went.Guild{},
		&went.Wallet{},
		&went.WalletLedger{},
		&went.DailySpend{},
		&ment.Item{},
		&ment.OraclePrice{},
		&tent.Auction{},
		&tent.Bid{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	return nil
}
