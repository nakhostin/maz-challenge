package mock

import (
	"time"

	"maz/domain/marketplace/entity"
	"maz/domain/shared"
)

// Items returns PRD-aligned sample marketplace catalog.
func Items(now time.Time) []entity.Item {
	if now.IsZero() {
		now = time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	}

	return []entity.Item{
		// Common — limit orders (listed, high supply)
		{
			ID:            ItemHealingDraught,
			SellerGuildID: GuildIronVanguard,
			OwnerGuildID:  GuildIronVanguard,
			Name:          "Healing Draught",
			ItemType:      shared.ItemTypeCommon,
			Status:        shared.ItemStatusListed,
			ListPrice:     500,
			CreatedAt:     now.Add(-48 * time.Hour),
		},
		{
			ID:            ItemArcaneThread,
			SellerGuildID: GuildShadowSyndicate,
			OwnerGuildID:  GuildShadowSyndicate,
			Name:          "Arcane Thread",
			ItemType:      shared.ItemTypeCommon,
			Status:        shared.ItemStatusListed,
			ListPrice:     1_200,
			CreatedAt:     now.Add(-36 * time.Hour),
		},
		{
			ID:            ItemIronwoodShield,
			SellerGuildID: GuildCrystalForge,
			OwnerGuildID:  GuildCrystalForge,
			Name:          "Ironwood Shield",
			ItemType:      shared.ItemTypeCommon,
			Status:        shared.ItemStatusAvailable,
			ListPrice:     3_500,
			CreatedAt:     now.Add(-12 * time.Hour),
		},

		// Rare — limit orders (limited supply, fixed price at listing)
		{
			ID:            ItemMoonsteelIngot,
			SellerGuildID: GuildCrystalForge,
			OwnerGuildID:  GuildCrystalForge,
			Name:          "Moonsteel Ingot",
			ItemType:      shared.ItemTypeRare,
			Status:        shared.ItemStatusListed,
			ListPrice:     25_000,
			CreatedAt:     now.Add(-72 * time.Hour),
		},
		{
			ID:            ItemPhoenixFeatherCloak,
			SellerGuildID: GuildIronVanguard,
			OwnerGuildID:  GuildIronVanguard,
			Name:          "Phoenix Feather Cloak",
			ItemType:      shared.ItemTypeRare,
			Status:        shared.ItemStatusListed,
			ListPrice:     45_000,
			CreatedAt:     now.Add(-24 * time.Hour),
		},

		// Legendary — globally unique, auction-only (in_auction when listed)
		{
			ID:            ItemSoulReaver,
			SellerGuildID: GuildShadowSyndicate,
			OwnerGuildID:  GuildShadowSyndicate,
			Name:          "Soul Reaver",
			ItemType:      shared.ItemTypeLegendary,
			Status:        shared.ItemStatusInAuction,
			ListPrice:     500_000,
			CreatedAt:     now.Add(-6 * time.Hour),
		},
		{
			ID:            ItemEyeOfTheDragon,
			SellerGuildID: GuildCrystalForge,
			OwnerGuildID:  GuildCrystalForge,
			Name:          "Eye of the Dragon",
			ItemType:      shared.ItemTypeLegendary,
			Status:        shared.ItemStatusInAuction,
			ListPrice:     750_000,
			CreatedAt:     now.Add(-3 * time.Hour),
		},
	}
}
