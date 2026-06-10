package entity

import "time"

// OraclePrice caches the last valid reference price from the external oracle (display only).
type OraclePrice struct {
	ItemName  string    `gorm:"type:text;primaryKey"`
	BasePrice int64     `gorm:"not null;check:base_price > 0"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (OraclePrice) TableName() string { return "oracle_prices" }
