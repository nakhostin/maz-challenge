package oracle_test

import (
	"testing"

	"maz/pkg/oracle"
)

func TestSanitizePrices_DropsInvalid(t *testing.T) {
	in := map[string]int64{
		"Valid":   100,
		"Zero":    0,
		"Negative": -5,
	}
	out := oracle.SanitizePrices(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 valid price, got %d", len(out))
	}
	if out["Valid"] != 100 {
		t.Fatalf("unexpected valid price: %d", out["Valid"])
	}
}

func TestSanitizePrices_EmptyInput(t *testing.T) {
	if oracle.SanitizePrices(nil) != nil {
		t.Fatal("expected nil for empty input")
	}
}
