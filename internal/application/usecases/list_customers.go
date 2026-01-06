package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"context"
)

type ListCustomersUseCase struct {
	repo ports.CustomerRepository
}

func NewListCustomersUseCase(repo ports.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{repo: repo}
}

func (uc *ListCustomersUseCase) Execute(ctx context.Context) ([]dto.CustomerDTO, error) {
	customers, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.CustomerDTO, len(customers))
	for i, c := range customers {
		dtos[i] = dto.CustomerDTO{
			ID:        string(c.ID),
			Name:      c.Name,
			Email:     c.Email,
			TaxID:     c.TaxID,
			CreatedAt: c.CreatedAt,
		}
	}
	return dtos, nil
}
