package usecases

import (
	"carigo/internal/application/ports"
	"context"
)

type DashboardStats struct {
	TotalCollected int64
	OpenInvoices   int64
	TotalRevenue   int64 
	TotalCustomers int64 
	PendingBalance int64 
}

type GetDashboardStatsUseCase struct {
	payRepo  ports.PaymentRepository
	invRepo  ports.InvoiceRepository
	custRepo ports.CustomerRepository
}

func NewGetDashboardStatsUseCase(pr ports.PaymentRepository, ir ports.InvoiceRepository, cr ports.CustomerRepository) *GetDashboardStatsUseCase {
	return &GetDashboardStatsUseCase{payRepo: pr, invRepo: ir, custRepo: cr}
}

func (uc *GetDashboardStatsUseCase) Execute(ctx context.Context) (*DashboardStats, error) {
	totalCollected, err := uc.payRepo.SumTotalCollected(ctx)
	if err != nil {
		return nil, err
	}
	openInvoices, err := uc.invRepo.CountAllOpen(ctx)
	if err != nil {
		return nil, err
	}

	totalRevenue, err := uc.invRepo.SumTotalAmount(ctx)
	if err != nil {
		return nil, err
	}

	totalCustomers, err := uc.custRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	pendingBalance := totalRevenue - totalCollected
	if pendingBalance < 0 {
		pendingBalance = 0 
	}

	return &DashboardStats{
		TotalCollected: totalCollected,
		OpenInvoices:   openInvoices,
		TotalRevenue:   totalRevenue,
		TotalCustomers: totalCustomers,
		PendingBalance: pendingBalance,
	}, nil
}
