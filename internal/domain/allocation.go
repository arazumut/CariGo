package domain

import (
	"time"
)

type AllocationID string

type Allocation struct {
	ID        AllocationID
	PaymentID PaymentID
	InvoiceID InvoiceID
	Amount    Money
	CreatedAt time.Time
}

func NewAllocation(id AllocationID, payment *Payment, invoice *Invoice, amount Money) (*Allocation, error) {
	if payment.AvailableAmount.currency != amount.currency || invoice.TotalAmount.currency != amount.currency {
		return nil, ErrCurrencyMismatch
	}

	if err := payment.UseFunds(amount); err != nil {
		return nil, err
	}

	if err := invoice.AllocatePayment(amount); err != nil {
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
