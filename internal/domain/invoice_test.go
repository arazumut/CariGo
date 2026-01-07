package domain_test

import (
	"testing"
	"time"
	"carigo/internal/domain"
)

func TestInvoice_Lifecycle(t *testing.T) {
	total, _ := domain.NewMoney(1000, "TRY")
	inv, err := domain.NewInvoice("INV-001", "CUST-001", total, time.Now(), time.Now().Add(24*time.Hour))
	
	if err != nil {
		t.Fatalf("failed to create invoice: %v", err)
	}

	if inv.Status != domain.InvoiceStatusOpen {
		t.Errorf("expected status OPEN, got %s", inv.Status)
	}

	payment1, _ := domain.NewMoney(400, "TRY")
	err = inv.AllocatePayment(payment1)
	if err != nil {
		t.Fatalf("failed to allocate payment: %v", err)
	}

	if inv.Status != domain.InvoiceStatusPartial {
		t.Errorf("expected status PARTIAL, got %s", inv.Status)
	}
	if inv.RemainingAmount().Amount() != 600 {
		t.Errorf("expected remaining 600, got %d", inv.RemainingAmount().Amount())
	}

	payment2, _ := domain.NewMoney(600, "TRY")
	err = inv.AllocatePayment(payment2)
	if err != nil {
		t.Fatalf("failed to allocate payment: %v", err)
	}

	if inv.Status != domain.InvoiceStatusPaid {
		t.Errorf("expected status PAID, got %s[", inv.Status)
	}

	extra, _ := domain.NewMoney(1, "TRY")
	err = inv.AllocatePayment(extra)
	if err != domain.ErrInvoiceAlreadyPaid {
		t.Errorf("expected ErrInvoiceAlreadyPaid, got %v", err)
	}
}
