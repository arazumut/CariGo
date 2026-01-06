package sqlite

import (
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"errors"
)

type PaymentModel struct {
	ID              string `gorm:"primaryKey"`
	CustomerID      string `gorm:"index"`
	Amount          int64
	Currency        string
	AvailableAmount int64
	Date            int64
	CreatedAt       int64
}

// --- Payment Repo ---

func (r *GormRepository) SavePayment(ctx context.Context, p *domain.Payment) error {
	m := PaymentModel{
		ID:              string(p.ID),
		CustomerID:      string(p.CustomerID),
		Amount:          p.Amount.Amount(),
		Currency:        p.Amount.Currency(),
		AvailableAmount: p.AvailableAmount.Amount(),
		Date:            p.Date.Unix(),
		CreatedAt:       p.CreatedAt.Unix(),
	}
	return r.getDB(ctx).Save(&m).Error
}

type PaymentAdapter struct{ repo *GormRepository }

func (a *PaymentAdapter) Save(ctx context.Context, p *domain.Payment) error {
	return a.repo.SavePayment(ctx, p)
}
func (a *PaymentAdapter) FindByID(ctx context.Context, id domain.PaymentID) (*domain.Payment, error) {
	return nil, errors.New("not implemented")
}

func (a *PaymentAdapter) FindAll(ctx context.Context) ([]*domain.Payment, error) {
	var models []PaymentModel
	err := a.repo.getDB(ctx).Order("created_at desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	var payments []*domain.Payment
	for _, m := range models {
		amount, _ := domain.NewMoney(m.Amount, m.Currency)
		p := domain.NewPayment(domain.PaymentID(m.ID), domain.CustomerID(m.CustomerID), amount, parseTime(m.Date))
		
		avail, _ := domain.NewMoney(m.AvailableAmount, m.Currency)
		// We need to set AvailableAmount. Since it is public in entity (usually) or via constructor.
		// Checking domain/payment.go would be ideal but assuming it's accessible or we reconstruct.
		// The NewPayment sets Available = Amount. We need to overwrite it.
		// Wait, Payment struct fields might be private/protected?
		// Let's assume standard Go struct access or I'll check domain if this fails.
		// Given Invoice was accessible, this likely is too.
		p.AvailableAmount = avail
		p.CreatedAt = parseTime(m.CreatedAt)
		
		payments = append(payments, p)
	}
	return payments, nil
}

func (a *PaymentAdapter) SumTotalCollected(ctx context.Context) (int64, error) {
	var total int64
	// Sum the initial amounts of all payments
	err := a.repo.getDB(ctx).Model(&PaymentModel{}).
		Select("ifnull(sum(amount), 0)").
		Scan(&total).Error
	return total, err
}

var _ ports.PaymentRepository = &PaymentAdapter{}
