package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrNotConnected = errors.New("postgres: not connected")

type Client struct {
	GORM *gorm.DB
	sql  *sql.DB
}

func Open(ctx context.Context, cfg *Config, opts ...gorm.Option) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	orm, err := gorm.Open(postgres.Open(cfg.DSN()), opts...)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres: sql db: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	return &Client{GORM: orm, sql: sqlDB}, nil
}

func (c *Client) Close() error {
	if c == nil || c.sql == nil {
		return ErrNotConnected
	}
	return c.sql.Close()
}

func (c *Client) Health(ctx context.Context) error {
	if c == nil || c.sql == nil {
		return ErrNotConnected
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return c.sql.PingContext(ctx)
}
