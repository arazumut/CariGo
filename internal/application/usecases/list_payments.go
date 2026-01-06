package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"context"
)

type ListPaymentsUseCase struct {
	repo ports.PaymentRepository
}

func NewListPaymentsUseCase(r ports.PaymentRepository) *ListPaymentsUseCase {
	return &ListPaymentsUseCase{repo: r}
}

func (uc *ListPaymentsUseCase) Execute(ctx context.Context) ([]dto.PaymentDTO, error) {
	payments, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.PaymentDTO, len(payments))
	for i, p := range payments {
		dtos[i] = dto.PaymentDTO{
			ID:              string(p.ID),
			CustomerID:      string(p.CustomerID),
			Amount:          float64(p.Amount.Amount()) / 100.0,
			AvailableAmount: float64(p.AvailableAmount.Amount()) / 100.0,
			Currency:        p.Amount.Currency(),
			Date:            p.Date.Format("2006-01-02"),
		}
	}
	return dtos, nil
}
