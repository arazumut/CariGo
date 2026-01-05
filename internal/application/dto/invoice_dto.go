package dto

import "time"

type CreateInvoiceRequest struct {
	CustomerID string   `json:"customer_id" binding:"required"`
	Amount     int64    `json:"amount" binding:"required,gt=0"` // Cents
	Currency   string   `json:"currency" binding:"required,len=3"`
	DueDate    time.Time `json:"due_date" binding:"required"`
}

type CreateInvoiceResponse struct {
	InvoiceID   string    `json:"invoice_id"`
	TotalAmount int64     `json:"total_amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
}
