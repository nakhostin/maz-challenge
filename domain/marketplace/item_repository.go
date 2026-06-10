package marketplace

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"maz/domain/marketplace/entity"
	"maz/domain/shared"
)

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *entity.Item) error {
	if item == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create item: %w", err)
	}
	return nil
}

func (r *itemRepository) Save(ctx context.Context, item *entity.Item) error {
	if item == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("save item: %w", err)
	}
	return nil
}

func (r *itemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Item, error) {
	return r.find(ctx, id, false)
}

func (r *itemRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*entity.Item, error) {
	return r.find(ctx, id, true)
}

func (r *itemRepository) List(ctx context.Context) ([]entity.Item, error) {
	var items []entity.Item
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	return items, nil
}

func (r *itemRepository) ExistsLegendaryName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Item{}).
		Where("name = ? AND item_type = ?", name, shared.ItemTypeLegendary).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *itemRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.Item, error) {
	var item entity.Item
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) find(ctx context.Context, id uuid.UUID, forUpdate bool) (*entity.Item, error) {
	q := r.db.WithContext(ctx)
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var item entity.Item
	if err := q.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

type oraclePriceRepository struct {
	db *gorm.DB
}

func NewOraclePriceRepository(db *gorm.DB) OraclePriceRepository {
	return &oraclePriceRepository{db: db}
}

func (r *oraclePriceRepository) List(ctx context.Context) (map[string]int64, error) {
	var rows []entity.OraclePrice
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list oracle prices: %w", err)
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.ItemName] = row.BasePrice
	}
	return out, nil
}

func (r *oraclePriceRepository) Upsert(ctx context.Context, price *entity.OraclePrice) error {
	if price == nil || price.ItemName == "" || price.BasePrice <= 0 {
		return shared.ErrInvalidState
	}
	return r.db.WithContext(ctx).Save(price).Error
}
