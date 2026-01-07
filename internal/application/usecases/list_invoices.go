package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"time"
)

type ListInvoicesUseCase struct {
	repo ports.InvoiceRepository
}

func NewListInvoicesUseCase(r ports.InvoiceRepository) *ListInvoicesUseCase {
	return &ListInvoicesUseCase{repo: r}
}

func (uc *ListInvoicesUseCase) Execute(ctx context.Context) ([]dto.InvoiceDTO, error) {
	invoices, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.InvoiceDTO, len(invoices))
	for i, inv := range invoices {
		dtos[i] = dto.InvoiceDTO{
			ID:          string(inv.ID),
			CustomerID:  string(inv.CustomerID),
			TotalAmount: float64(inv.TotalAmount.Amount()) / 100.0,
			PaidAmount:  float64(inv.PaidAmount.Amount()) / 100.0,
			Currency:    inv.TotalAmount.Currency(),
			Status:      string(inv.Status),
			IssueDate:   inv.IssueDate.Format("2006-01-02"),
			DueDate:     inv.DueDate.Format("2006-01-02"),
		}
		
		if inv.Status == domain.InvoiceStatusOpen && inv.DueDate.Before(time.Now()) {
			
		}
	}
	return dtos, nil
}
