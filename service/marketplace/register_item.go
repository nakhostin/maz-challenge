package marketplace

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/marketplace/entity"
	"maz/domain/shared"
	dtrading "maz/domain/trading"
	tentity "maz/domain/trading/entity"
)

// RegisterItem lists a new item for sale (Common/Rare) or starts a Legendary auction.
func (s *Service) RegisterItem(ctx context.Context, cmd RegisterItemCommand) (*ItemView, error) {
	sellerID, err := uuid.Parse(cmd.SellerGuildID)
	if err != nil {
		return nil, shared.ErrInvalidState
	}
	itemType, err := parseItemType(cmd.ItemType)
	if err != nil {
		return nil, err
	}
	if cmd.ListPrice <= 0 || strings.TrimSpace(cmd.Name) == "" {
		return nil, shared.ErrInvalidState
	}

	var view *ItemView
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		items := s.items(tx)
		if itemType == shared.ItemTypeLegendary {
			exists, err := items.ExistsLegendaryName(ctx, cmd.Name)
			if err != nil {
				return err
			}
			if exists {
				return shared.ErrDuplicateLegendary
			}
		}

		item := &entity.Item{
			ID:            uuid.New(),
			SellerGuildID: sellerID,
			OwnerGuildID:  sellerID,
			Name:          strings.TrimSpace(cmd.Name),
			ItemType:      itemType,
			Status:        shared.ItemStatusAvailable,
			ListPrice:     cmd.ListPrice,
			CreatedAt:     cmd.Now,
		}

		switch itemType {
		case shared.ItemTypeCommon, shared.ItemTypeRare:
			item.Status = shared.ItemStatusListed
		case shared.ItemTypeLegendary:
			item.Status = shared.ItemStatusInAuction
		}

		if err := items.Create(ctx, item); err != nil {
			return err
		}

		if itemType == shared.ItemTypeLegendary {
			auction := &tentity.Auction{
				ID:            uuid.New(),
				ItemID:        item.ID,
				SellerGuildID: sellerID,
				Status:        shared.AuctionStatusActive,
				StartingPrice: cmd.ListPrice,
				EndsAt:        cmd.Now.Add(s.auctionDuration),
				CreatedAt:     cmd.Now,
			}
			if err := dtrading.NewAuctionRepository(tx).Create(ctx, auction); err != nil {
				return err
			}
		}

		oracle, err := s.oracle(tx).List(ctx)
		if err != nil {
			return err
		}
		v := toItemView(item, oracle)
		view = &v
		return nil
	})
	if err != nil {
		return nil, err
	}
	return view, nil
}

func parseItemType(raw string) (shared.ItemType, error) {
	switch shared.ItemType(strings.ToLower(strings.TrimSpace(raw))) {
	case shared.ItemTypeCommon, shared.ItemTypeRare, shared.ItemTypeLegendary:
		return shared.ItemType(strings.ToLower(strings.TrimSpace(raw))), nil
	default:
		return "", shared.ErrInvalidItemType
	}
}
