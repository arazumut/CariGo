package ports

import (
	"carigo/internal/domain"
	"context"
)

// InvoiceRepository defines access to Invoice storage.
type InvoiceRepository interface {
	Save(ctx context.Context, invoice *domain.Invoice) error
	FindByID(ctx context.Context, id domain.InvoiceID) (*domain.Invoice, error)
	// FindOpenByCustomer returns all non-PAID/VOID invoices for a customer, typically ordered by DueDate (FIFO).
	FindOpenByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Invoice, error)
	FindAll(ctx context.Context) ([]*domain.Invoice, error)
	FindByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Invoice, error)
	CountAllOpen(ctx context.Context) (int64, error)
	SumTotalAmount(ctx context.Context) (int64, error)
}

// PaymentRepository defines access to Payment storage.
type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) error
	FindByID(ctx context.Context, id domain.PaymentID) (*domain.Payment, error)
	FindAll(ctx context.Context) ([]*domain.Payment, error)
	FindByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Payment, error)
	SumTotalCollected(ctx context.Context) (int64, error)
}

// CustomerRepository defines access to Customer storage.
type CustomerRepository interface {
	Save(ctx context.Context, customer *domain.Customer) error
	FindByID(ctx context.Context, id domain.CustomerID) (*domain.Customer, error)
	FindAll(ctx context.Context) ([]*domain.Customer, error)
	Count(ctx context.Context) (int64, error)
}

// AllocationRepository defines access to Allocation storage.
type AllocationRepository interface {
	Save(ctx context.Context, allocation *domain.Allocation) error
}

// TransactionManager handles database transactions.
// It allows UseCases to wrap multiple repo calls in a single atomic block.
type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
