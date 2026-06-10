package mock

import (
	"time"

	"maz/domain/marketplace/entity"
)
type Data struct {
	Items        []entity.Item
	OraclePrices []entity.OraclePrice
}

// New returns PRD-aligned mock marketplace data at the given time.
func New(now time.Time) Data {
	return Data{
		Items:        Items(now),
		OraclePrices: OraclePrices(now),
	}
}

// ItemNames returns every mock item name (for oracle mock service registration).
func ItemNames() []string {
	names := make([]string, 0, len(Items(time.Time{})))
	for _, item := range Items(time.Time{}) {
		names = append(names, item.Name)
	}
	return names
}
