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

func (a *CustomerAdapter) FindAll(ctx context.Context) ([]*domain.Customer, error) {
	var models []CustomerModel
	if err := a.repo.getDB(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	var customers []*domain.Customer
	for _, m := range models {
		c, err := domain.NewCustomer(domain.CustomerID(m.ID), m.Name, m.Email, m.TaxID)
		if err != nil {
			return nil, err
		}
		c.CreatedAt = parseTime(m.CreatedAt)
		
		customers = append(customers, c)
	}
	return customers, nil
}

func (a *CustomerAdapter) Count(ctx context.Context) (int64, error) {
	var count int64
	err := a.repo.getDB(ctx).Model(&CustomerModel{}).Count(&count).Error
	return count, err
}

var _ ports.CustomerRepository = &CustomerAdapter{}
