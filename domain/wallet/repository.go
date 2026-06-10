package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"

	"maz/domain/wallet/entity"
)

// WalletRepository persists guild wallet state and daily spend tracking.
type WalletRepository interface {
	GetByGuildID(ctx context.Context, guildID uuid.UUID) (*entity.Wallet, error)
	GetForUpdate(ctx context.Context, guildID uuid.UUID) (*entity.Wallet, error)
	Save(ctx context.Context, wallet *entity.Wallet) error
	GetGuild(ctx context.Context, guildID uuid.UUID) (*entity.Guild, error)
	GetTodaySpend(ctx context.Context, guildID uuid.UUID, date time.Time) (int64, error)
	UpsertTodaySpend(ctx context.Context, spend *entity.DailySpend) error
}

// LedgerRepository appends wallet audit entries.
type LedgerRepository interface {
	Append(ctx context.Context, entry *entity.WalletLedger) error
	AppendMany(ctx context.Context, entries []entity.WalletLedger) error
	FindByIdempotencyKey(ctx context.Context, key string) (*entity.WalletLedger, error)
}
