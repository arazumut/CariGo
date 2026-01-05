package domain

import (
	"time"
)

type InvoiceStatus string

const (
	InvoiceStatusOpen    InvoiceStatus = "OPEN"
	InvoiceStatusPartial InvoiceStatus = "PARTIAL"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusVoid    InvoiceStatus = "VOID"
)

type InvoiceID string

// Invoice is the core entity.
// It holds the state of debt.
type Invoice struct {
	ID          InvoiceID
	CustomerID  CustomerID
	TotalAmount Money
	PaidAmount  Money // Amount already paid/allocated
	IssueDate   time.Time
	DueDate     time.Time
	Status      InvoiceStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewInvoice creates a fresh invoice in OPEN state.
func NewInvoice(id InvoiceID, customerID CustomerID, total Money, issueDate, dueDate time.Time) (*Invoice, error) {
	if total.IsZero() || total.amount < 0 {
		return nil, ErrNegativeAmount
	}
	
	zeroMoney, _ := NewMoney(0, total.Currency())

	return &Invoice{
		ID:          id,
		CustomerID:  customerID,
		TotalAmount: total,
		PaidAmount:  zeroMoney,
		IssueDate:   issueDate,
		DueDate:     dueDate,
		Status:      InvoiceStatusOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// RemainingAmount calculates how much is left to pay.
func (i *Invoice) RemainingAmount() Money {
	remaining, _ := i.TotalAmount.Subtract(i.PaidAmount)
	return remaining
}

// AllocatePayment attempts to apply a payment amount to this invoice.
// It updates the status automatically.
func (i *Invoice) AllocatePayment(amount Money) error {
	if i.Status == InvoiceStatusPaid || i.Status == InvoiceStatusVoid {
		return ErrInvoiceAlreadyPaid
	}

	if amount.currency != i.TotalAmount.currency {
		return ErrCurrencyMismatch
	}

	remaining := i.RemainingAmount()
	
	// Check for overpayment attempts if strict
	// For now, we allow full allocation logic to handle logic, but Entity should define boundaries.
	// If allocation amount > remaining, we should probably error here or cap it?
	// The business rule says: "Overpayment support". 
	// However, usually an Invoice entity itself shouldn't store more PaidAmount than Total.
	// Any excess should be handled by the Service/UseCase as "Account Credit".
	// So strictly, Invoice accepts up to Remaining.
	
	isOverpayment, _ := amount.GreaterThan(remaining)
	if isOverpayment {
		return ErrOverPaymentNotAllowed
	}

	newPaid, err := i.PaidAmount.Add(amount)
	if err != nil {
		return err
	}
	i.PaidAmount = newPaid
	i.updateStatus()
	i.UpdatedAt = time.Now()
	
	return nil
}

func (i *Invoice) updateStatus() {
	if i.PaidAmount.Equals(i.TotalAmount) {
		i.Status = InvoiceStatusPaid
	} else if i.PaidAmount.amount > 0 {
		i.Status = InvoiceStatusPartial
	} else {
		i.Status = InvoiceStatusOpen
	}
}
