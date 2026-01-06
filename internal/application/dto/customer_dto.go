package dto

import "time"

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	TaxID string `json:"tax_id" binding:"required"`
}

type CreateCustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	TaxID     string    `json:"tax_id"`
	CreatedAt time.Time `json:"created_at"`
}
