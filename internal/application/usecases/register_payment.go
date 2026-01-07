package usecases

import (
	"carigo/internal/application/dto"
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"fmt"
)

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

func (uc *RegisterPaymentUseCase) Execute(ctx context.Context, req dto.RegisterPaymentRequest) (*dto.RegisterPaymentResponse, error) {
	amount, err := domain.NewMoney(req.Amount, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("invalid money: %w", err)
	}

	date := req.Date
	if date.IsZero() {
		date = uc.clock.Now()
	}

	paymentID := domain.PaymentID(fmt.Sprintf("PAY-%d", date.UnixNano()))
	payment := domain.NewPayment(paymentID, domain.CustomerID(req.CustomerID), amount, date)
	
	allocatedItems := []dto.AllocatedInvoiceParams{}
	totalAllocated := int64(0)

	err = uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.paymentRepo.Save(ctx, payment); err != nil {
			return err
		}
		invoices, err := uc.invoiceRepo.FindOpenByCustomer(ctx, domain.CustomerID(req.CustomerID))
		if err != nil {
			return err
		}

		for _, inv := range invoices {
			if payment.AvailableAmount.IsZero() {
				break
			}
			remainingDebt := inv.RemainingAmount()
			
			available := payment.AvailableAmount

			var allocationAmount domain.Money
			
			if remainingDebt.Currency() != available.Currency() {
				continue
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

			allocID := domain.AllocationID(fmt.Sprintf("AL-%s-%s", payment.ID, inv.ID))
			allocation, err := domain.NewAllocation(allocID, payment, inv, allocationAmount)
			if err != nil {
				return err
			}

			if err := uc.invoiceRepo.Save(ctx, inv); err != nil {
				return err
			}
			if err := uc.paymentRepo.Save(ctx, payment); err != nil {
				return err
			}

			if err := uc.allocationRepo.Save(ctx, allocation); err != nil {
				return err
			}

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
