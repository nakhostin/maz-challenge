package mock

import "github.com/google/uuid"

// Guild IDs match dev seed in docs/DATABASE_ER.md.
var (
	GuildIronVanguard   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	GuildShadowSyndicate = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	GuildCrystalForge   = uuid.MustParse("33333333-3333-3333-3333-333333333333")
)

// Item IDs are fixed for predictable dev and tests.
var (
	ItemHealingDraught     = uuid.MustParse("a1000001-0000-4000-8000-000000000001")
	ItemArcaneThread       = uuid.MustParse("a1000001-0000-4000-8000-000000000002")
	ItemIronwoodShield     = uuid.MustParse("a1000001-0000-4000-8000-000000000003")
	ItemMoonsteelIngot     = uuid.MustParse("a2000002-0000-4000-8000-000000000001")
	ItemPhoenixFeatherCloak = uuid.MustParse("a2000002-0000-4000-8000-000000000002")
	ItemSoulReaver         = uuid.MustParse("a3000003-0000-4000-8000-000000000001")
	ItemEyeOfTheDragon      = uuid.MustParse("a3000003-0000-4000-8000-000000000002")
)
