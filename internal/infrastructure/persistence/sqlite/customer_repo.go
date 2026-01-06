package sqlite

import (
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
)

type CustomerModel struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string
	TaxID     string
	CreatedAt int64
	UpdatedAt int64
}

// --- Customer Repo ---

func (r *GormRepository) SaveCustomer(ctx context.Context, c *domain.Customer) error {
	m := CustomerModel{
		ID:        string(c.ID),
		Name:      c.Name,
		Email:     c.Email,
		TaxID:     c.TaxID,
		CreatedAt: c.CreatedAt.Unix(),
		UpdatedAt: c.UpdatedAt.Unix(),
	}
	return r.getDB(ctx).Save(&m).Error
}

func (r *GormRepository) FindCustomerByID(ctx context.Context, id domain.CustomerID) (*domain.Customer, error) {
	var m CustomerModel
	if err := r.getDB(ctx).First(&m, "id = ?", string(id)).Error; err != nil {
		return nil, err
	}
	return domain.NewCustomer(domain.CustomerID(m.ID), m.Name, m.Email, m.TaxID)
}

type CustomerAdapter struct{ repo *GormRepository }

func (a *CustomerAdapter) Save(ctx context.Context, c *domain.Customer) error {
	return a.repo.SaveCustomer(ctx, c)
}
func (a *CustomerAdapter) FindByID(ctx context.Context, id domain.CustomerID) (*domain.Customer, error) {
	return a.repo.FindCustomerByID(ctx, id)
}

var _ ports.CustomerRepository = &CustomerAdapter{}
