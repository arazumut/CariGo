package dto

import "time"

type RegisterPaymentRequest struct {
	CustomerID string  `json:"customer_id" binding:"required"`
	Amount     int64   `json:"amount" binding:"required,gt=0"`
	Currency   string  `json:"currency" binding:"required,len=3"`
	Date       time.Time `json:"date"`
	Notes      string  `json:"notes"`
}

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
