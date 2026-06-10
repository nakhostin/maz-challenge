package shared

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrDailyCapExceeded    = errors.New("daily purchase cap exceeded")
	ErrInvalidState        = errors.New("invalid state")
	ErrSelfBid             = errors.New("cannot bid on own item")
	ErrSelfPurchase        = errors.New("cannot purchase own item")
	ErrBidTooLow           = errors.New("bid must be at least 5% higher than current highest")
	ErrCannotWithdrawBid   = errors.New("highest bidder cannot withdraw bid")
	ErrAuctionNotActive    = errors.New("auction is not active")
	ErrDuplicateLegendary  = errors.New("legendary item name already exists")
	ErrInvalidItemType     = errors.New("invalid item type for operation")
)

const MinBidIncreasePercent = 5

const (
	AntiSnipeWindow = 5 // minutes
	AntiSnipeExtend = 5 // minutes
)
