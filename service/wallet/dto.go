package wallet

import "time"

// WalletView is returned by GET /guilds/{id}/wallet.
type WalletView struct {
	GuildID    string `json:"guild_id"`
	GuildName  string `json:"guild_name"`
	Balance    int64  `json:"balance"`
	Reserved   int64  `json:"reserved"`
	Available  int64  `json:"available"`
	TodaySpent int64  `json:"today_spent"`
	DailyCap   int64  `json:"daily_cap"`
}

// ReserveParams configures a wallet reservation.
type ReserveParams struct {
	GuildID        string
	Amount         int64
	ReferenceType  string
	ReferenceID    string
	IdempotencyKey string
	Now            time.Time
}

// ReleaseParams configures a wallet release.
type ReleaseParams struct {
	GuildID       string
	Amount        int64
	ReferenceType string
	ReferenceID   string
	Now           time.Time
}

// DebitParams configures settling a prior reservation.
type DebitParams struct {
	GuildID       string
	Amount        int64
	ReferenceType string
	ReferenceID   string
	Now           time.Time
}

// CreditParams configures a wallet credit.
type CreditParams struct {
	GuildID       string
	Amount        int64
	ReferenceType string
	ReferenceID   string
	Now           time.Time
}

// SpendParams configures an immediate purchase debit (no prior reservation).
type SpendParams struct {
	GuildID        string
	Amount         int64
	ReferenceType  string
	ReferenceID    string
	IdempotencyKey string
	Now            time.Time
}
