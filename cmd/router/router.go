package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	appmw "maz/cmd/middleware"
	"maz/cmd/handler"
)

type Deps struct {
	Items    *handler.ItemHandler
	Bids     *handler.BidHandler
	Auctions *handler.AuctionHandler
	Wallets  *handler.WalletHandler
}

func New(deps Deps) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Dragon Market",
		ErrorHandler: handler.ErrorHandler,
	})

	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/items", deps.Items.List)
	app.Get("/items/:id", deps.Items.Get)
	app.Post("/items", appmw.RequireGuildID, deps.Items.Register)
	app.Post("/items/:id/bid", appmw.RequireGuildID, deps.Bids.Place)
	app.Delete("/items/:id/bid/:bid_id", appmw.RequireGuildID, deps.Bids.Withdraw)

	app.Get("/auctions", deps.Auctions.List)
	app.Get("/auctions/:id", deps.Auctions.Get)

	app.Get("/guilds/:id/wallet", deps.Wallets.Get)

	return app
}
