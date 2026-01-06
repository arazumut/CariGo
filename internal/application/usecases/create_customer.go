package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"fmt"
	"time"
)

type CreateCustomerUseCase struct {
	repo ports.CustomerRepository
}

func NewCreateCustomerUseCase(repo ports.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{repo: repo}
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CreateCustomerResponse, error) {
	// Generate simple ID like "CUST-TIMESTAMP"
	id := domain.CustomerID(fmt.Sprintf("CUST-%d", time.Now().UnixNano()))
	
	customer, err := domain.NewCustomer(id, req.Name, req.Email, req.TaxID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, customer); err != nil {
		return nil, err
	}

	return &dto.CreateCustomerResponse{
		ID:    string(customer.ID),
		Name:  customer.Name,
		Email: customer.Email,
	}, nil
}
