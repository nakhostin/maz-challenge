package wallet

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"maz/domain/shared"
	"maz/domain/wallet/entity"
)

type ledgerRepository struct {
	db *gorm.DB
}

// NewLedgerRepository returns a GORM-backed LedgerRepository.
func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) Append(ctx context.Context, entry *entity.WalletLedger) error {
	if entry == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return fmt.Errorf("append ledger entry: %w", err)
	}
	return nil
}

func (r *ledgerRepository) AppendMany(ctx context.Context, entries []entity.WalletLedger) error {
	if len(entries) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&entries).Error; err != nil {
		return fmt.Errorf("append ledger entries: %w", err)
	}
	return nil
}

func (r *ledgerRepository) FindByIdempotencyKey(ctx context.Context, key string) (*entity.WalletLedger, error) {
	if key == "" {
		return nil, shared.ErrNotFound
	}
	var entry entity.WalletLedger
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&entry).Error
	if err != nil {
		return nil, mapLedgerErr(err)
	}
	return &entry, nil
}

func mapLedgerErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.ErrNotFound
	}
	return fmt.Errorf("ledger repository: %w", err)
}
