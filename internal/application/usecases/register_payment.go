package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"fmt"
)

// RegisterPaymentUseCase handles receiving money and allocating it to open invoices.
type RegisterPaymentUseCase struct {
	paymentRepo    ports.PaymentRepository
	invoiceRepo    ports.InvoiceRepository
	allocationRepo ports.AllocationRepository
	txManager      ports.TransactionManager
	clock          ports.Clock
}

func NewRegisterPaymentUseCase(
	pr ports.PaymentRepository,
	ir ports.InvoiceRepository,
	ar ports.AllocationRepository,
	tm ports.TransactionManager,
	clk ports.Clock,
) *RegisterPaymentUseCase {
	return &RegisterPaymentUseCase{
		paymentRepo:    pr,
		invoiceRepo:    ir,
		allocationRepo: ar,
		txManager:      tm,
		clock:          clk,
	}
}

// Execute performs the payment registration and FIFO allocation.
func (uc *RegisterPaymentUseCase) Execute(ctx context.Context, req dto.RegisterPaymentRequest) (*dto.RegisterPaymentResponse, error) {
	// 1. Prepare Data
	amount, err := domain.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid money: %w", err)
	}

	date := req.Date
	if date.IsZero() {
		date = uc.clock.Now()
	}

	paymentID := domain.PaymentID(fmt.Sprintf("PAY-%d", date.UnixNano())) // Simple ID generation
	payment := domain.NewPayment(paymentID, domain.CustomerID(req.CustomerID), amount, date)
	
	allocatedItems := []dto.AllocatedInvoiceParams{}
	totalAllocated := int64(0)

	// 2. Transactional Block
	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		// A. Save Payment
		if err := uc.paymentRepo.Save(ctx, payment); err != nil {
			return err
		}

		// B. Fetch Open Invoices (FIFO by default from Repo)
		invoices, err := uc.invoiceRepo.FindOpenByCustomer(ctx, domain.CustomerID(req.CustomerID))
		if err != nil {
			return err
		}

		// C. Allocate Logic (FIFO)
		for _, inv := range invoices {
			// If payment is exhausted, stop.
			if payment.AvailableAmount.IsZero() {
				break
			}

			// How much can we allocate to this invoice?
			remainingDebt := inv.RemainingAmount()
			
			// Available funds
			available := payment.AvailableAmount

			// Determine allocation amount: min(remainingDebt, available)
			// But since we have strict Money types, logic requires comparison
			var allocationAmount domain.Money
			
			// We check strictly currency match first (handled by entities usually, but safe check)
			if remainingDebt.Currency() != available.Currency() {
				continue // Skip mismatch currencies
			}

			isDebtLarger, _ := remainingDebt.GreaterThan(available)
			if isDebtLarger {
				allocationAmount = available
			} else {
				allocationAmount = remainingDebt
			}

			if allocationAmount.IsZero() {
				continue
			}

			// Create Domain Allocation
			allocID := domain.AllocationID(fmt.Sprintf("AL-%s-%s", payment.ID, inv.ID))
			allocation, err := domain.NewAllocation(allocID, payment, inv, allocationAmount)
			if err != nil {
				return err // Should not happen given logic above, but handle it
			}

			// Save Invoice (Status update)
			if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
				return err
			}

			// Save Payment (AvailableAmount update)
			if err := uc.paymentRepo.Save(ctx, payment); err != nil {
				return err
			}

			// Save Allocation
			if err := uc.allocationRepo.Save(ctx, allocation); err != nil {
				return err
			}

			// Add to response
			allocatedItems = append(allocatedItems, dto.AllocatedInvoiceParams{
				InvoiceID: string(inv.ID),
				Amount:    allocationAmount.Amount(),
			})
			totalAllocated += allocationAmount.Amount()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.RegisterPaymentResponse{
		PaymentID:         string(payment.ID),
		AllocatedAmount:   totalAllocated,
		RemainingBalance:  payment.AvailableAmount.Amount(),
		AllocatedInvoices: allocatedItems,
	}, nil
}
