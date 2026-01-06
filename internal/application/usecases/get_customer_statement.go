package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"sort"
)

type GetCustomerStatementUseCase struct {
	custRepo ports.CustomerRepository
	invRepo  ports.InvoiceRepository
	payRepo  ports.PaymentRepository
}

func NewGetCustomerStatementUseCase(c ports.CustomerRepository, i ports.InvoiceRepository, p ports.PaymentRepository) *GetCustomerStatementUseCase {
	return &GetCustomerStatementUseCase{
		custRepo: c,
		invRepo:  i,
		payRepo:  p,
	}
}

func (uc *GetCustomerStatementUseCase) Execute(ctx context.Context, customerID string) (*dto.CustomerStatementDTO, error) {
	cid := domain.CustomerID(customerID)
	
	// 1. Get Customer
	customer, err := uc.custRepo.FindByID(ctx, cid)
	if err != nil {
		return nil, err
	}

	// 2. Get Invoices
	invoices, err := uc.invRepo.FindByCustomer(ctx, cid)
	if err != nil {
		return nil, err
	}

	// 3. Get Payments
	payments, err := uc.payRepo.FindByCustomer(ctx, cid)
	if err != nil {
		return nil, err
	}

	// 4. Merge and Sort
	var transactions []dto.StatementItem

	for _, inv := range invoices {
		transactions = append(transactions, dto.StatementItem{
			Date:        inv.IssueDate,
			Type:        "FATURA",
			ReferenceID: string(inv.ID),
			Description: "Satış Faturası",
			Debt:        float64(inv.TotalAmount.Amount()) / 100.0,
			Credit:      0,
			Currency:    inv.TotalAmount.Currency(),
		})
	}

	for _, pay := range payments {
		transactions = append(transactions, dto.StatementItem{
			Date:        pay.Date,
			Type:        "TAHSİLAT",
			ReferenceID: string(pay.ID),
			Description: "Ödeme Alındı",
			Debt:        0,
			Credit:      float64(pay.Amount.Amount()) / 100.0,
			Currency:    pay.Amount.Currency(),
		})
	}

	// Sort by Date
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date.Before(transactions[j].Date)
	})

	// 5. Calculate Running Balance
	balance := 0.0
	for i := range transactions {
		// Borç artırır, Alacak (Ödeme) düşürür
		balance += transactions[i].Debt
		balance -= transactions[i].Credit
		transactions[i].Balance = balance
	}

	return &dto.CustomerStatementDTO{
		Customer: dto.CustomerDTO{
			ID:    string(customer.ID),
			Name:  customer.Name,
			Email: customer.Email,
			TaxID: customer.TaxID,
		},
		Transactions: transactions,
		FinalBalance: balance,
		Currency:     "TRY", // MVP assumption or derived from transactions
	}, nil
}
