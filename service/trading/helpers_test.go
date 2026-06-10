package trading_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	tentity "maz/domain/trading/entity"
	strading "maz/service/trading"
)

func TestMinRequiredBid_FirstBidUsesStartingPrice(t *testing.T) {
	auction := &tentity.Auction{StartingPrice: 500_000}
	got := strading.MinRequiredBidForTest(auction)
	if got != 500_000 {
		t.Fatalf("expected 500000, got %d", got)
	}
}

func TestMinRequiredBid_RequiresFivePercentIncrease(t *testing.T) {
	highest := int64(500_000)
	auction := &tentity.Auction{
		StartingPrice: 500_000,
		HighestBidID:  ptrUUID(uuid.New()),
		HighestAmount: &highest,
	}
	got := strading.MinRequiredBidForTest(auction)
	want := int64(525_000)
	if got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestShouldExtend_InsideAntiSnipeWindow(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	endsAt := now.Add(3 * time.Minute)
	if !strading.ShouldExtendForTest(endsAt, now) {
		t.Fatal("expected anti-snipe extension")
	}
}

func TestShouldExtend_OutsideWindow(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	endsAt := now.Add(10 * time.Minute)
	if strading.ShouldExtendForTest(endsAt, now) {
		t.Fatal("expected no extension")
	}
}

func ptrUUID(id uuid.UUID) *uuid.UUID { return &id }
