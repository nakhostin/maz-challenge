package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"maz/cmd/handler"
	"maz/cmd/router"
	"maz/cmd/worker"
	mmarket "maz/domain/marketplace/mock"
	tmock "maz/domain/trading/mock"
	wmock "maz/domain/wallet/mock"
	migration "maz/migrations"
	"maz/pkg/postgres"
	smarket "maz/service/marketplace"
	strading "maz/service/trading"
	swallet "maz/service/wallet"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	dbConfig := postgres.ConfigFromEnv()
	db, err := postgres.Open(ctx, dbConfig)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Health(ctx); err != nil {
		log.Fatalf("database health: %v", err)
	}

	if err := migration.AutoMigrate(ctx, db.GORM); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	auctionDuration := envDuration("AUCTION_DURATION", 24*time.Hour)
	if err := seed(ctx, db.GORM, auctionDuration); err != nil {
		log.Fatalf("seed: %v", err)
	}

	walletSvc := swallet.NewService(db.GORM)
	tradingSvc := strading.NewService(db.GORM, walletSvc, auctionDuration)
	marketSvc := smarket.NewService(db.GORM, tradingSvc, auctionDuration)

	httpAddr := envString("HTTP_ADDR", ":8080")
	app := router.New(router.Deps{
		Items:    handler.NewItemHandler(marketSvc),
		Bids:     handler.NewBidHandler(tradingSvc),
		Auctions: handler.NewAuctionHandler(tradingSvc),
		Wallets:  handler.NewWalletHandler(walletSvc),
	})

	workerCtx, stopWorker := context.WithCancel(ctx)
	defer stopWorker()
	go worker.RunAuctionCloser(workerCtx, tradingSvc, envDuration("AUCTION_CLOSER_INTERVAL", time.Minute))

	go func() {
		log.Printf("listening on %s", httpAddr)
		if err := app.Listen(httpAddr); err != nil {
			log.Fatalf("http server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stopWorker()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func seed(ctx context.Context, db *gorm.DB, auctionDuration time.Duration) error {
	if err := wmock.Seed(ctx, db); err != nil {
		return err
	}
	if err := mmarket.Seed(ctx, db); err != nil {
		return err
	}
	return tmock.Seed(ctx, db, time.Now().UTC(), auctionDuration)
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	if sec, err := strconv.Atoi(raw); err == nil {
		return time.Duration(sec) * time.Second
	}
	return fallback
}
