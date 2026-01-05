package dto

import "time"

// RegisterPaymentRequest carries data to register a new payment.
type RegisterPaymentRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Amount     int64   `json:"amount" binding:"required,gt=0"` // Cents
	Currency   string  `json:"currency" binding:"required,len=3"`
	Date       time.Time `json:"date"` // Optional, defaults to now
	Notes      string  `json:"notes"`
}

// RegisterPaymentResponse is the result of the operation.
type RegisterPaymentResponse struct {
	PaymentID         string `json:"payment_id"`
	AllocatedAmount   int64  `json:"allocated_amount"`
	RemainingBalance  int64  `json:"remaining_balance"`
	AllocatedInvoices []AllocatedInvoiceParams `json:"allocated_invoices"`
}

type AllocatedInvoiceParams struct {
	InvoiceID string `json:"invoice_id"`
	Amount    int64  `json:"amount"`
}
