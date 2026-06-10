package trading

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"maz/domain/shared"
	"maz/domain/trading/entity"
)

type auctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository returns a GORM-backed AuctionRepository.
// Pass a transaction handle from db.Transaction to scope writes atomically.
func NewAuctionRepository(db *gorm.DB) AuctionRepository {
	return &auctionRepository{db: db}
}

func (r *auctionRepository) Create(ctx context.Context, auction *entity.Auction) error {
	if auction == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Create(auction).Error; err != nil {
		return fmt.Errorf("create auction: %w", err)
	}
	return nil
}

func (r *auctionRepository) Save(ctx context.Context, auction *entity.Auction) error {
	if auction == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(auction).Error; err != nil {
		return fmt.Errorf("save auction: %w", err)
	}
	return nil
}

func (r *auctionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Auction, error) {
	return r.findOne(ctx, id, false)
}

func (r *auctionRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*entity.Auction, error) {
	return r.findOne(ctx, id, true)
}

func (r *auctionRepository) GetByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error) {
	return r.findByItemID(ctx, itemID, false, false)
}

func (r *auctionRepository) GetActiveByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error) {
	return r.findByItemID(ctx, itemID, true, false)
}

func (r *auctionRepository) GetForUpdateByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error) {
	return r.findByItemID(ctx, itemID, true, true)
}

func (r *auctionRepository) ListActive(ctx context.Context) ([]entity.Auction, error) {
	var auctions []entity.Auction
	err := r.db.WithContext(ctx).
		Preload("Bids", bidOrder).
		Where("status = ?", shared.AuctionStatusActive).
		Order("ends_at ASC").
		Find(&auctions).Error
	if err != nil {
		return nil, fmt.Errorf("list active auctions: %w", err)
	}
	return auctions, nil
}

func (r *auctionRepository) ListExpiredActiveForUpdate(ctx context.Context, now time.Time) ([]entity.Auction, error) {
	var auctions []entity.Auction
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Preload("Bids", bidOrder).
		Where("status = ? AND ends_at <= ?", shared.AuctionStatusActive, now).
		Find(&auctions).Error
	if err != nil {
		return nil, fmt.Errorf("list expired active auctions: %w", err)
	}
	return auctions, nil
}

func (r *auctionRepository) findOne(ctx context.Context, id uuid.UUID, forUpdate bool) (*entity.Auction, error) {
	q := r.db.WithContext(ctx)
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var auction entity.Auction
	err := q.Preload("Bids", bidOrder).First(&auction, "id = ?", id).Error
	if err != nil {
		return nil, mapAuctionErr(err)
	}
	return &auction, nil
}

func (r *auctionRepository) findByItemID(ctx context.Context, itemID uuid.UUID, activeOnly, forUpdate bool) (*entity.Auction, error) {
	q := r.db.WithContext(ctx).Where("item_id = ?", itemID)
	if activeOnly {
		q = q.Where("status = ?", shared.AuctionStatusActive)
	}
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var auction entity.Auction
	err := q.Preload("Bids", bidOrder).Order("created_at DESC").First(&auction).Error
	if err != nil {
		return nil, mapAuctionErr(err)
	}
	return &auction, nil
}

func bidOrder(db *gorm.DB) *gorm.DB {
	return db.Order("created_at ASC")
}

func mapAuctionErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.ErrNotFound
	}
	return err
}
