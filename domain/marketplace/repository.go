package marketplace

import (
	"context"

	"github.com/google/uuid"

	"maz/domain/marketplace/entity"
)

// ItemRepository persists marketplace items.
type ItemRepository interface {
	Create(ctx context.Context, item *entity.Item) error
	Save(ctx context.Context, item *entity.Item) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Item, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*entity.Item, error)
	List(ctx context.Context) ([]entity.Item, error)
	ExistsLegendaryName(ctx context.Context, name string) (bool, error)
}

// OraclePriceRepository reads cached oracle reference prices.
type OraclePriceRepository interface {
	List(ctx context.Context) (map[string]int64, error)
}
