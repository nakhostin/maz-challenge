package shared

// ItemType classifies marketplace goods.
type ItemType string

const (
	ItemTypeCommon    ItemType = "common"
	ItemTypeRare      ItemType = "rare"
	ItemTypeLegendary ItemType = "legendary"
)

// ItemStatus tracks item lifecycle.
type ItemStatus string

const (
	ItemStatusAvailable ItemStatus = "available"
	ItemStatusListed    ItemStatus = "listed"
	ItemStatusInAuction ItemStatus = "in_auction"
	ItemStatusSold      ItemStatus = "sold"
)

// LedgerEntryType is a wallet ledger operation.
type LedgerEntryType string

const (
	LedgerEntryCredit  LedgerEntryType = "credit"
	LedgerEntryDebit   LedgerEntryType = "debit"
	LedgerEntryReserve LedgerEntryType = "reserve"
	LedgerEntryRelease LedgerEntryType = "release"
)

// AuctionStatus tracks auction lifecycle.
type AuctionStatus string

const (
	AuctionStatusActive    AuctionStatus = "active"
	AuctionStatusClosed    AuctionStatus = "closed"
	AuctionStatusSettled   AuctionStatus = "settled"
	AuctionStatusCancelled AuctionStatus = "cancelled"
)

// BidStatus tracks bid lifecycle.
type BidStatus string

const (
	BidStatusActive    BidStatus = "active"
	BidStatusOutbid    BidStatus = "outbid"
	BidStatusWithdrawn BidStatus = "withdrawn"
	BidStatusWon       BidStatus = "won"
)
