package oracle

import "context"

// Client fetches reference prices from an external oracle service.
type Client interface {
	FetchPrices(ctx context.Context) (map[string]int64, error)
}
