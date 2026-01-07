package dto

type InvoiceDTO struct {
	ID          string  `json:"id"`
	CustomerID  string  `json:"customer_id"`
	TotalAmount float64 `json:"total_amount"`
	PaidAmount  float64 `json:"paid_amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	IssueDate   string  `json:"issue_date"`
	DueDate     string  `json:"due_date"`
}
