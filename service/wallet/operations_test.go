package wallet_test

import (
	"testing"

	"maz/domain/wallet/entity"
	swallet "maz/service/wallet"
)

func TestCanCommit_InsufficientFunds(t *testing.T) {
	w := &entity.Wallet{Balance: 100, Reserved: 50}
	g := &entity.Guild{DailyPurchaseCap: 1_000}
	err := swallet.CanCommitForTest(w, g, 0, 60)
	if err == nil {
		t.Fatal("expected insufficient funds")
	}
}

func TestCanCommit_DailyCapExceeded(t *testing.T) {
	w := &entity.Wallet{Balance: 1_000_000, Reserved: 0}
	g := &entity.Guild{DailyPurchaseCap: 500_000}
	err := swallet.CanCommitForTest(w, g, 450_000, 100_000)
	if err == nil {
		t.Fatal("expected daily cap exceeded")
	}
}

func TestCanCommit_OK(t *testing.T) {
	w := &entity.Wallet{Balance: 1_000_000, Reserved: 0}
	g := &entity.Guild{DailyPurchaseCap: 500_000}
	if err := swallet.CanCommitForTest(w, g, 100_000, 50_000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
