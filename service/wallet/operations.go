package wallet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	"maz/domain/wallet/entity"
)

// Reserve reserves funds inside an existing transaction.
func (s *Service) Reserve(ctx context.Context, tx *gorm.DB, p ReserveParams) (uuid.UUID, error) {
	wallets, ledgers := s.repos(tx)
	guildID, err := uuid.Parse(p.GuildID)
	if err != nil {
		return uuid.Nil, shared.ErrInvalidState
	}
	refID, err := uuid.Parse(p.ReferenceID)
	if err != nil {
		return uuid.Nil, shared.ErrInvalidState
	}
	if p.IdempotencyKey != "" {
		if existing, err := ledgers.FindByIdempotencyKey(ctx, p.IdempotencyKey); err == nil {
			return existing.ID, nil
		} else if err != shared.ErrNotFound {
			return uuid.Nil, err
		}
	}

	w, err := wallets.GetForUpdate(ctx, guildID)
	if err != nil {
		return uuid.Nil, err
	}
	guild, err := ensureGuild(ctx, wallets, w, guildID)
	if err != nil {
		return uuid.Nil, err
	}

	todaySpend, err := wallets.GetTodaySpend(ctx, guildID, p.Now)
	if err != nil {
		return uuid.Nil, err
	}
	if err := canCommit(w, guild, todaySpend, p.Amount); err != nil {
		return uuid.Nil, err
	}

	w.Reserved += p.Amount
	todaySpend += p.Amount

	var idem *string
	if p.IdempotencyKey != "" {
		idem = &p.IdempotencyKey
	}
	entry := entity.WalletLedger{
		ID:             uuid.New(),
		GuildID:        guildID,
		EntryType:      shared.LedgerEntryReserve,
		Amount:         p.Amount,
		ReferenceType:  p.ReferenceType,
		ReferenceID:    refID,
		IdempotencyKey: idem,
		CreatedAt:      p.Now,
	}
	if err := ledgers.Append(ctx, &entry); err != nil {
		return uuid.Nil, err
	}
	if err := wallets.Save(ctx, w); err != nil {
		return uuid.Nil, err
	}
	if err := wallets.UpsertTodaySpend(ctx, &entity.DailySpend{
		GuildID:   guildID,
		SpendDate: p.Now,
		Amount:    todaySpend,
	}); err != nil {
		return uuid.Nil, err
	}
	return entry.ID, nil
}

// Release releases a prior reservation inside an existing transaction.
func (s *Service) Release(ctx context.Context, tx *gorm.DB, p ReleaseParams) error {
	wallets, ledgers := s.repos(tx)
	guildID, err := uuid.Parse(p.GuildID)
	if err != nil {
		return shared.ErrInvalidState
	}
	refID, err := uuid.Parse(p.ReferenceID)
	if err != nil {
		return shared.ErrInvalidState
	}

	w, err := wallets.GetForUpdate(ctx, guildID)
	if err != nil {
		return err
	}
	if w.Reserved < p.Amount {
		return shared.ErrInvalidState
	}
	w.Reserved -= p.Amount

	todaySpend, err := wallets.GetTodaySpend(ctx, guildID, p.Now)
	if err != nil {
		return err
	}
	if todaySpend >= p.Amount {
		todaySpend -= p.Amount
	}

	if err := ledgers.Append(ctx, &entity.WalletLedger{
		ID:            uuid.New(),
		GuildID:       guildID,
		EntryType:     shared.LedgerEntryRelease,
		Amount:        p.Amount,
		ReferenceType: p.ReferenceType,
		ReferenceID:   refID,
		CreatedAt:     p.Now,
	}); err != nil {
		return err
	}
	if err := wallets.Save(ctx, w); err != nil {
		return err
	}
	return wallets.UpsertTodaySpend(ctx, &entity.DailySpend{
		GuildID:   guildID,
		SpendDate: p.Now,
		Amount:    todaySpend,
	})
}

// Debit converts a reservation into a final debit inside an existing transaction.
func (s *Service) Debit(ctx context.Context, tx *gorm.DB, p DebitParams) error {
	wallets, ledgers := s.repos(tx)
	guildID, err := uuid.Parse(p.GuildID)
	if err != nil {
		return shared.ErrInvalidState
	}
	refID, err := uuid.Parse(p.ReferenceID)
	if err != nil {
		return shared.ErrInvalidState
	}

	w, err := wallets.GetForUpdate(ctx, guildID)
	if err != nil {
		return err
	}
	if w.Reserved < p.Amount || w.Balance < p.Amount {
		return shared.ErrInsufficientFunds
	}
	w.Balance -= p.Amount
	w.Reserved -= p.Amount

	if err := ledgers.Append(ctx, &entity.WalletLedger{
		ID:            uuid.New(),
		GuildID:       guildID,
		EntryType:     shared.LedgerEntryDebit,
		Amount:        p.Amount,
		ReferenceType: p.ReferenceType,
		ReferenceID:   refID,
		CreatedAt:     p.Now,
	}); err != nil {
		return err
	}
	return wallets.Save(ctx, w)
}

// Credit adds gold inside an existing transaction.
func (s *Service) Credit(ctx context.Context, tx *gorm.DB, p CreditParams) error {
	wallets, ledgers := s.repos(tx)
	guildID, err := uuid.Parse(p.GuildID)
	if err != nil {
		return shared.ErrInvalidState
	}
	refID, err := uuid.Parse(p.ReferenceID)
	if err != nil {
		return shared.ErrInvalidState
	}

	w, err := wallets.GetForUpdate(ctx, guildID)
	if err != nil {
		return err
	}
	w.Balance += p.Amount
	if err := ledgers.Append(ctx, &entity.WalletLedger{
		ID:            uuid.New(),
		GuildID:       guildID,
		EntryType:     shared.LedgerEntryCredit,
		Amount:        p.Amount,
		ReferenceType: p.ReferenceType,
		ReferenceID:   refID,
		CreatedAt:     p.Now,
	}); err != nil {
		return err
	}
	return wallets.Save(ctx, w)
}

func canCommit(w *entity.Wallet, guild *entity.Guild, todaySpend, amount int64) error {
	if w.AvailableBalance() < amount {
		return shared.ErrInsufficientFunds
	}
	if todaySpend+amount > guild.DailyPurchaseCap {
		return shared.ErrDailyCapExceeded
	}
	return nil
}

func ensureGuild(ctx context.Context, wallets interface {
	GetGuild(ctx context.Context, guildID uuid.UUID) (*entity.Guild, error)
}, w *entity.Wallet, guildID uuid.UUID) (*entity.Guild, error) {
	if w.Guild != nil {
		return w.Guild, nil
	}
	g, err := wallets.GetGuild(ctx, guildID)
	if err != nil {
		return nil, fmt.Errorf("load guild: %w", err)
	}
	w.Guild = g
	return g, nil
}
