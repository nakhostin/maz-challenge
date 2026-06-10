package mock

import (
	"context"
	"errors"
	"math/rand"
	"time"

	mmarket "maz/domain/marketplace/mock"
)

// ErrUnavailable simulates a flaky oracle dependency.
var ErrUnavailable = errors.New("oracle: temporarily unavailable")

// Client is a PRD-aligned mock oracle with occasional failures and bad values.
type Client struct {
	base map[string]int64
	rng  *rand.Rand
}

func NewClient(now time.Time) *Client {
	return &Client{
		base: mmarket.OraclePriceMap(now),
		rng:  rand.New(rand.NewSource(now.UnixNano())),
	}
}

func (c *Client) FetchPrices(_ context.Context) (map[string]int64, error) {
	if c.rng.Float64() < 0.05 {
		return nil, ErrUnavailable
	}

	out := make(map[string]int64, len(c.base))
	for name, base := range c.base {
		delta := int64(c.rng.Intn(21) - 10)
		price := base + delta
		if c.rng.Float64() < 0.1 {
			price = 0
		} else if c.rng.Float64() < 0.05 {
			price = -1
		}
		out[name] = price
	}
	return out, nil
}
