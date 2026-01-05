package domain_test

import (
	"testing"
	"carigo/internal/domain"
)

func TestNewMoney(t *testing.T) {
	m, err := domain.NewMoney(100, "TRY")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.Amount() != 100 {
		t.Errorf("expected 100, got %d", m.Amount())
	}
	if m.Currency() != "TRY" {
		t.Errorf("expected TRY, got %s", m.Currency())
	}

	_, err = domain.NewMoney(-50, "TRY")
	if err != domain.ErrNegativeAmount {
		t.Errorf("expected ErrNegativeAmount, got %v", err)
	}
}

func TestMoney_Add(t *testing.T) {
	m1, _ := domain.NewMoney(100, "USD")
	m2, _ := domain.NewMoney(50, "USD")
	
	sum, err := m1.Add(m2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.Amount() != 150 {
		t.Errorf("expected 150, got %d", sum.Amount())
	}

	m3, _ := domain.NewMoney(50, "EUR")
	_, err = m1.Add(m3)
	if err != domain.ErrCurrencyMismatch {
		t.Errorf("expected mismatch error, got %v", err)
	}
}

func TestMoney_Subtract(t *testing.T) {
	m1, _ := domain.NewMoney(100, "USD")
	m2, _ := domain.NewMoney(40, "USD")

	res, err := m1.Subtract(m2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Amount() != 60 {
		t.Errorf("expected 60, got %d", res.Amount())
	}
}
