package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"fmt"
)

type CreateInvoiceUseCase struct {
	invoiceRepo ports.InvoiceRepository
	clock       ports.Clock
}

func NewCreateInvoiceUseCase(ir ports.InvoiceRepository, clk ports.Clock) *CreateInvoiceUseCase {
	return &CreateInvoiceUseCase{
		invoiceRepo: ir,
		clock:       clk,
	}
}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, req dto.CreateInvoiceRequest) (*dto.CreateInvoiceResponse, error) {
	// 1. Create Money Value Object
	total, err := domain.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// 2. Generate ID
	id := domain.InvoiceID(fmt.Sprintf("INV-%d", uc.clock.Now().UnixNano()))

	// 3. Create Entity
	// IssueDate is Now
	inv, err := domain.NewInvoice(id, domain.CustomerID(req.CustomerID), total, uc.clock.Now(), req.DueDate)
	if err != nil {
		return nil, err
	}

	// 4. Persist
	if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
		return nil, err
	}

	return &dto.CreateInvoiceResponse{
		InvoiceID:   string(inv.ID),
		TotalAmount: inv.TotalAmount.Amount(),
		Currency:    inv.TotalAmount.Currency(),
		Status:      string(inv.Status),
		DueDate:     inv.DueDate,
	}, nil
}
