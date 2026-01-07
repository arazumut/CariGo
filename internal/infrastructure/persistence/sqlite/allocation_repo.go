package sqlite

import (
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
)

type AllocationModel struct {
	ID        string `gorm:"primaryKey"`
	PaymentID string `gorm:"index"`
	InvoiceID string `gorm:"index"`
	Amount    int64
	Currency  string
	CreatedAt int64
}

func (r *GormRepository) SaveAllocation(ctx context.Context, a *domain.Allocation) error {
	m := AllocationModel{
		ID:        string(a.ID),
		PaymentID: string(a.PaymentID),
		InvoiceID: string(a.InvoiceID),
		Amount:    a.Amount.Amount(),
		Currency:  a.Amount.Currency(),
		CreatedAt: a.CreatedAt.Unix(),
	}
	return r.getDB(ctx).Save(&m).Error
}

type AllocationAdapter struct{ repo *GormRepository }

func (a *AllocationAdapter) Save(ctx context.Context, al *domain.Allocation) error {
	return a.repo.SaveAllocation(ctx, al)
}

var _ ports.AllocationRepository = &AllocationAdapter{}
