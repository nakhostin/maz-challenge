package marketplace

import (
	"context"

	"github.com/google/uuid"

	"maz/domain/shared"
)

// ListItems returns all items with oracle reference prices when available.
func (s *Service) ListItems(ctx context.Context) ([]ItemView, error) {
	items, err := s.items(s.db).List(ctx)
	if err != nil {
		return nil, err
	}
	oracle, err := s.oracle(s.db).List(ctx)
	if err != nil {
		return nil, err
	}
	return toItemViews(items, oracle), nil
}

// GetItem returns one item by ID.
func (s *Service) GetItem(ctx context.Context, itemID uuid.UUID) (*ItemView, error) {
	item, err := s.items(s.db).GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	oracle, err := s.oracle(s.db).List(ctx)
	if err != nil {
		return nil, err
	}
	view := toItemView(item, oracle)
	return &view, nil
}

// Ensure shared errors are referenced for lint stability in queries.
var _ = shared.ErrNotFound
