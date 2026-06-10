package marketplace

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"maz/domain/shared"
	swallet "maz/service/wallet"
)

// PurchaseItemCommand is input for POST /items/{id}/purchase.
type PurchaseItemCommand struct {
	ItemID         string
	BuyerGuildID   string
	IdempotencyKey string
	Now            time.Time
}

// PurchaseItem buys a listed Common/Rare item at its fixed list price.
func (s *Service) PurchaseItem(ctx context.Context, cmd PurchaseItemCommand) (*ItemView, error) {
	itemID, err := uuid.Parse(cmd.ItemID)
	if err != nil {
		return nil, shared.ErrInvalidState
	}
	buyerID, err := uuid.Parse(cmd.BuyerGuildID)
	if err != nil {
		return nil, shared.ErrInvalidState
	}

	var view *ItemView
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		items := s.items(tx)

		if cmd.IdempotencyKey != "" {
			if entry, err := s.ledgers(tx).FindByIdempotencyKey(ctx, cmd.IdempotencyKey); err == nil {
				item, err := items.GetByID(ctx, entry.ReferenceID)
				if err != nil {
					return err
				}
				oracle, err := s.oracle(tx).List(ctx)
				if err != nil {
					return err
				}
				v := toItemView(item, oracle)
				view = &v
				return nil
			} else if err != shared.ErrNotFound {
				return err
			}
		}

		item, err := items.GetForUpdate(ctx, itemID)
		if err != nil {
			return err
		}
		if item.IsLegendary() {
			return shared.ErrInvalidItemType
		}
		if item.Status != shared.ItemStatusListed {
			return shared.ErrInvalidState
		}
		if buyerID == item.SellerGuildID {
			return shared.ErrSelfPurchase
		}

		if err := s.wallet.Spend(ctx, tx, swallet.SpendParams{
			GuildID:        buyerID.String(),
			Amount:         item.ListPrice,
			ReferenceType:  "purchase",
			ReferenceID:    item.ID.String(),
			IdempotencyKey: cmd.IdempotencyKey,
			Now:            cmd.Now,
		}); err != nil {
			return err
		}
		if err := s.wallet.Credit(ctx, tx, swallet.CreditParams{
			GuildID:       item.SellerGuildID.String(),
			Amount:        item.ListPrice,
			ReferenceType: "sale",
			ReferenceID:   item.ID.String(),
			Now:           cmd.Now,
		}); err != nil {
			return err
		}

		item.Status = shared.ItemStatusSold
		item.OwnerGuildID = buyerID
		if err := items.Save(ctx, item); err != nil {
			return err
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
