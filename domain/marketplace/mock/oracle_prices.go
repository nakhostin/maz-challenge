package mock

import (
	"time"

	"maz/domain/marketplace/entity"
)

// OraclePrices returns cached reference prices for mock items (display only per PRD).
func OraclePrices(now time.Time) []entity.OraclePrice {
	if now.IsZero() {
		now = time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	}

	return []entity.OraclePrice{
		{ItemName: "Healing Draught", BasePrice: 480, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Arcane Thread", BasePrice: 1_150, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Ironwood Shield", BasePrice: 3_400, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Moonsteel Ingot", BasePrice: 24_500, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Phoenix Feather Cloak", BasePrice: 44_000, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Soul Reaver", BasePrice: 520_000, UpdatedAt: now.Add(-30 * time.Second)},
		{ItemName: "Eye of the Dragon", BasePrice: 780_000, UpdatedAt: now.Add(-30 * time.Second)},
	}
}

// OraclePriceMap returns item name → base price for list views.
func OraclePriceMap(now time.Time) map[string]int64 {
	prices := OraclePrices(now)
	out := make(map[string]int64, len(prices))
	for _, p := range prices {
		out[p.ItemName] = p.BasePrice
	}
	return out
}
