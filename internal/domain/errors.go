package domain

import "errors"

var (
	// Money Errors
	ErrExample = errors.New("example error")
	ErrNegativeAmount = errors.New("amount cannot be negative")
	ErrCurrencyMismatch = errors.New("cannot operate on different currencies")
	ErrInvalidCurrency = errors.New("invalid currency")

	// Invoice Errors
	ErrInvalidInvoiceState = errors.New("invalid invoice state transition")
	ErrInvoiceAlreadyPaid = errors.New("invoice is already paid")
)
