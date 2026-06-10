package trading

import (
	"context"
	"time"

	"github.com/google/uuid"

	"maz/domain/trading/entity"
)

// AuctionRepository persists the Auction aggregate root and its bids.
type AuctionRepository interface {
	Create(ctx context.Context, auction *entity.Auction) error
	Save(ctx context.Context, auction *entity.Auction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Auction, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*entity.Auction, error)
	GetByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error)
	GetActiveByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error)
	GetForUpdateByItemID(ctx context.Context, itemID uuid.UUID) (*entity.Auction, error)
	ListActive(ctx context.Context) ([]entity.Auction, error)
	ListExpiredActiveForUpdate(ctx context.Context, now time.Time) ([]entity.Auction, error)
}

// BidRepository provides bid lookups used by write paths (e.g. idempotency).
type BidRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Bid, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*entity.Bid, error)
}
