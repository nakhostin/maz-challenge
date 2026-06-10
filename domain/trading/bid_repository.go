package trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	"maz/domain/trading/entity"
)

type bidRepository struct {
	db *gorm.DB
}

// NewBidRepository returns a GORM-backed BidRepository.
func NewBidRepository(db *gorm.DB) BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bid, error) {
	var bid entity.Bid
	err := r.db.WithContext(ctx).First(&bid, "id = ?", id).Error
	if err != nil {
		return nil, mapBidErr(err)
	}
	return &bid, nil
}

func (r *bidRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.Bid, error) {
	if key == "" {
		return nil, shared.ErrNotFound
	}
	var bid entity.Bid
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&bid).Error
	if err != nil {
		return nil, mapBidErr(err)
	}
	return &bid, nil
}

func mapBidErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.ErrNotFound
	}
	return fmt.Errorf("bid repository: %w", err)
}
