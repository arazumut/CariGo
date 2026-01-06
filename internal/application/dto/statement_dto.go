package dto

import "time"

type StatementItem struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"` // INVOICE or PAYMENT
	ReferenceID string    `json:"reference_id"`
	Description string    `json:"description"`
	Debt        float64   `json:"debt"`   // Borç (Fatura Tutarı)
	Credit      float64   `json:"credit"` // Alacak (Ödeme Tutarı)
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
}

type CustomerStatementDTO struct {
	Customer     CustomerDTO     `json:"customer"`
	Transactions []StatementItem `json:"transactions"`
	FinalBalance float64         `json:"final_balance"`
	Currency     string          `json:"currency"`
}
