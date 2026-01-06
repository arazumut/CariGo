package usecases

import (
	"carigo/internal/application/ports"
	"context"
)

type DashboardStats struct {
	TotalCollected int64
	OpenInvoices   int64
	TotalRevenue   int64 // Toplam Ciro (Kesilen Faturalar)
	TotalCustomers int64 // Müşteri Sayısı
	PendingBalance int64 // Bekleyen Alacak (TotalRevenue - TotalCollected kabaca veya (TotalAmount - PaidAmount))
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
	// 1. Total Collected
	totalCollected, err := uc.payRepo.SumTotalCollected(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Count Open Invoices
	openInvoices, err := uc.invRepo.CountAllOpen(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Total Revenue (Total Invoiced Amount)
	totalRevenue, err := uc.invRepo.SumTotalAmount(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Total Customers
	totalCustomers, err := uc.custRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// 5. Pending Balance (Basitçe Ciro - Tahsilat demek doğru olmaz, çünkü açık faturaların toplam kalan tutarı daha doğrudur)
	// Ancak şimdilik PRATİK olsun diye: Toplam Kesilen - Toplam Tahsilat = Piyasada Kalan Para diyebiliriz (kabaca)
	// Veya daha doğrusu Open Invoices toplamını çekmek olurdu ama repo methodu yok.
	// Şimdilik: TotalRevenue - TotalCollected
	pendingBalance := totalRevenue - totalCollected
	if pendingBalance < 0 {
		pendingBalance = 0 // Overpayment situation
	}

	return &DashboardStats{
		TotalCollected: totalCollected,
		OpenInvoices:   openInvoices,
		TotalRevenue:   totalRevenue,
		TotalCustomers: totalCustomers,
		PendingBalance: pendingBalance,
	}, nil
}
