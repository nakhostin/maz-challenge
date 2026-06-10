package wallet

import (
	"maz/domain/shared"
	"maz/domain/wallet/entity"
)

// CanCommitForTest exposes wallet commitment checks for unit tests.
func CanCommitForTest(w *entity.Wallet, guild *entity.Guild, todaySpend, amount int64) error {
	return canCommit(w, guild, todaySpend, amount)
}

var _ = shared.ErrInsufficientFunds
