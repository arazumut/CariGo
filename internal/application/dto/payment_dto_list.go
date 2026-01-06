package dto

type PaymentDTO struct {
	ID              string  `json:"id"`
	CustomerID      string  `json:"customer_id"`
	Amount          float64 `json:"amount"`
	AvailableAmount float64 `json:"available_amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
}
