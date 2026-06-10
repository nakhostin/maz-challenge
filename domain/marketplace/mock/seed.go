package mock

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Seed inserts mock items and oracle prices. Idempotent on item/oracle primary keys.
func Seed(ctx context.Context, db *gorm.DB) error {
	data := New(time.Now().UTC())

	for i := range data.Items {
		item := data.Items[i]
		if err := db.WithContext(ctx).Save(&item).Error; err != nil {
			return fmt.Errorf("seed item %q: %w", item.Name, err)
		}
	}
	for i := range data.OraclePrices {
		price := data.OraclePrices[i]
		if err := db.WithContext(ctx).Save(&price).Error; err != nil {
			return fmt.Errorf("seed oracle price %q: %w", price.ItemName, err)
		}
	}
	return nil
}
