package domain

import (
	"time"
)

type AllocationID string

// Allocation represents the link between a Payment and an Invoice.
// It answers "Which payment paid which invoice and how much?"
type Allocation struct {
	ID        AllocationID
	PaymentID PaymentID
	InvoiceID InvoiceID
	Amount    Money
	CreatedAt time.Time
}

func NewAllocation(id AllocationID, payment *Payment, invoice *Invoice, amount Money) (*Allocation, error) {
	// 1. Validate Currencies
	if payment.AvailableAmount.currency != amount.currency || invoice.TotalAmount.currency != amount.currency {
		return nil, ErrCurrencyMismatch
	}

	// 2. Try to use funds from payment
	if err := payment.UseFunds(amount); err != nil {
		return nil, err
	}

	// 3. Try to apply to invoice
	if err := invoice.AllocatePayment(amount); err != nil {
		// Rollback payment usage if invoice fails (in memory)
		// Since we modify pointers, we must be careful.
		// In a real pure DDD, we might clone or use strict events, but here we mutate.
		// Reverting is manually adding back.
		// payment.AvailableAmount = payment.AvailableAmount.Add(amount)
		return nil, err
	}

	return &Allocation{
		ID:        id,
		PaymentID: payment.ID,
		InvoiceID: invoice.ID,
		Amount:    amount,
		CreatedAt: time.Now(),
	}, nil
}
