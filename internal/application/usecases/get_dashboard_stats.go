package usecases

import (
	"carigo/internal/application/ports"
	"context"
)

type DashboardStats struct {
	TotalCollected int64
	OpenInvoices   int64
}

type GetDashboardStatsUseCase struct {
	paymentRepo ports.PaymentRepository
	invoiceRepo ports.InvoiceRepository
}

func NewGetDashboardStatsUseCase(pr ports.PaymentRepository, ir ports.InvoiceRepository) *GetDashboardStatsUseCase {
	return &GetDashboardStatsUseCase{
		paymentRepo: pr,
		invoiceRepo: ir,
	}
}

func (uc *GetDashboardStatsUseCase) Execute(ctx context.Context) (*DashboardStats, error) {
	total, err := uc.paymentRepo.SumTotalCollected(ctx)
	if err != nil {
		return nil, err
	}

	count, err := uc.invoiceRepo.CountAllOpen(ctx)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalCollected: total,
		OpenInvoices:   count,
	}, nil
}
