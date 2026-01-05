package domain

import (
	"errors"
	"time"
)

// CustomerID identifies a customer unique.
type CustomerID string

// Customer represents the entity that owes money.
type Customer struct {
	ID        CustomerID
	Name      string
	Email     string
	TaxID     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCustomer creates a new customer.
// Validation logic plays here.
func NewCustomer(id CustomerID, name, email, taxID string) (*Customer, error) {
	if id == "" {
		return nil, errors.New("customer ID is required")
	}
	if name == "" {
		return nil, errors.New("customer name is required")
	}
	return &Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		TaxID:     taxID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
