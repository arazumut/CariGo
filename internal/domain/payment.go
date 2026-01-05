package domain

import (
	"time"
)

type PaymentID string

// Payment represents a money transfer from a customer.
type Payment struct {
	ID              PaymentID
	CustomerID      CustomerID
	Amount          Money
	AvailableAmount Money // Amount not yet allocated to invoices
	Date            time.Time
	Notes           string
	CreatedAt       time.Time
}

func NewPayment(id PaymentID, customerID CustomerID, amount Money, date time.Time) *Payment {
	return &Payment{
		ID:              id,
		CustomerID:      customerID,
		Amount:          amount,
		AvailableAmount: amount, // Initially all available
		Date:            date,
		CreatedAt:       time.Now(),
	}
}

// UseFunds deducts from AvailableAmount.
func (p *Payment) UseFunds(amount Money) error {
	if amount.currency != p.AvailableAmount.currency {
		return ErrCurrencyMismatch
	}

	// Check if enough funds
	greater, _ := amount.GreaterThan(p.AvailableAmount)
	if greater {
		return ErrInsufficientPaymentBalance
	}

	newAvailable, err := p.AvailableAmount.Subtract(amount)
	if err != nil {
		return err
	}
	p.AvailableAmount = newAvailable
	return nil
}
