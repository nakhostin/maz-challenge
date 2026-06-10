package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"maz/domain/shared"
	"maz/domain/wallet/entity"
)

type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository returns a GORM-backed WalletRepository.
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetByGuildID(ctx context.Context, guildID uuid.UUID) (*entity.Wallet, error) {
	return r.findWallet(ctx, guildID, false)
}

func (r *walletRepository) GetForUpdate(ctx context.Context, guildID uuid.UUID) (*entity.Wallet, error) {
	return r.findWallet(ctx, guildID, true)
}

func (r *walletRepository) Save(ctx context.Context, wallet *entity.Wallet) error {
	if wallet == nil {
		return shared.ErrInvalidState
	}
	if err := r.db.WithContext(ctx).Save(wallet).Error; err != nil {
		return fmt.Errorf("save wallet: %w", err)
	}
	return nil
}

func (r *walletRepository) GetGuild(ctx context.Context, guildID uuid.UUID) (*entity.Guild, error) {
	var guild entity.Guild
	err := r.db.WithContext(ctx).First(&guild, "id = ?", guildID).Error
	if err != nil {
		return nil, mapWalletErr(err)
	}
	return &guild, nil
}

func (r *walletRepository) GetTodaySpend(ctx context.Context, guildID uuid.UUID, date time.Time) (int64, error) {
	spendDate := truncateDate(date)
	var row entity.DailySpend
	err := r.db.WithContext(ctx).
		Where("guild_id = ? AND spend_date = ?", guildID, spendDate).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get today spend: %w", err)
	}
	return row.Amount, nil
}

func (r *walletRepository) UpsertTodaySpend(ctx context.Context, spend *entity.DailySpend) error {
	if spend == nil {
		return shared.ErrInvalidState
	}
	spend.SpendDate = truncateDate(spend.SpendDate)
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}, {Name: "spend_date"}},
			DoUpdates: clause.AssignmentColumns([]string{"amount"}),
		}).
		Create(spend).Error
	if err != nil {
		return fmt.Errorf("upsert today spend: %w", err)
	}
	return nil
}

func (r *walletRepository) findWallet(ctx context.Context, guildID uuid.UUID, forUpdate bool) (*entity.Wallet, error) {
	q := r.db.WithContext(ctx)
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var w entity.Wallet
	err := q.Preload("Guild").First(&w, "guild_id = ?", guildID).Error
	if err != nil {
		return nil, mapWalletErr(err)
	}
	return &w, nil
}

func truncateDate(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func mapWalletErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shared.ErrNotFound
	}
	return err
}
