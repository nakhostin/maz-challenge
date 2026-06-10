package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dwallet "maz/domain/wallet"
)

// Service orchestrates wallet use cases.
type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) repos(tx *gorm.DB) (dwallet.WalletRepository, dwallet.LedgerRepository) {
	return dwallet.NewWalletRepository(tx), dwallet.NewLedgerRepository(tx)
}

// GetWallet returns the guild wallet view for handlers.
func (s *Service) GetWallet(ctx context.Context, guildID uuid.UUID) (*WalletView, error) {
	wallets, _ := s.repos(s.db)
	w, err := wallets.GetByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}
	today, err := wallets.GetTodaySpend(ctx, guildID, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	guild := w.Guild
	if guild == nil {
		g, err := wallets.GetGuild(ctx, guildID)
		if err != nil {
			return nil, err
		}
		guild = g
	}
	return &WalletView{
		GuildID:    guildID.String(),
		GuildName:  guild.Name,
		Balance:    w.Balance,
		Reserved:   w.Reserved,
		Available:  w.AvailableBalance(),
		TodaySpent: today,
		DailyCap:   guild.DailyPurchaseCap,
	}, nil
}
