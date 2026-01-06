package sqlite

import (
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"errors"
)

type InvoiceModel struct {
	ID          string `gorm:"primaryKey"`
	CustomerID  string `gorm:"index"`
	TotalAmount int64
	Currency    string
	PaidAmount  int64
	Status      string
	IssueDate   int64
	DueDate     int64
	CreatedAt   int64
	UpdatedAt   int64
}

// --- Invoice Repo ---

func (r *GormRepository) SaveInvoice(ctx context.Context, i *domain.Invoice) error {
	m := InvoiceModel{
		ID:          string(i.ID),
		CustomerID:  string(i.CustomerID),
		TotalAmount: i.TotalAmount.Amount(),
		Currency:    i.TotalAmount.Currency(),
		PaidAmount:  i.PaidAmount.Amount(),
		Status:      string(i.Status),
		IssueDate:   i.IssueDate.Unix(),
		DueDate:     i.DueDate.Unix(),
		CreatedAt:   i.CreatedAt.Unix(),
		UpdatedAt:   i.UpdatedAt.Unix(),
	}
	return r.getDB(ctx).Save(&m).Error
}

type InvoiceAdapter struct{ repo *GormRepository }

func (a *InvoiceAdapter) Save(ctx context.Context, i *domain.Invoice) error {
	return a.repo.SaveInvoice(ctx, i)
}
func (a *InvoiceAdapter) FindByID(ctx context.Context, id domain.InvoiceID) (*domain.Invoice, error) {
	return nil, errors.New("not implemented")
} 
func (a *InvoiceAdapter) FindOpenByCustomer(ctx context.Context, cid domain.CustomerID) ([]*domain.Invoice, error) {
	var models []InvoiceModel
	// Find invoices where status is OPEN or PARTIAL, ordered by DueDate ASC (FIFO)
	err := a.repo.getDB(ctx).
		Where("customer_id = ? AND status IN ?", string(cid), []string{string(domain.InvoiceStatusOpen), string(domain.InvoiceStatusPartial)}).
		Order("due_date asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	for _, m := range models {
		inv, err := a.mapToDomain(m)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (a *InvoiceAdapter) FindAll(ctx context.Context) ([]*domain.Invoice, error) {
	var models []InvoiceModel
	err := a.repo.getDB(ctx).Order("created_at desc").Find(&models).Error
	if err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	for _, m := range models {
		inv, err := a.mapToDomain(m)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (a *InvoiceAdapter) FindByCustomer(ctx context.Context, cid domain.CustomerID) ([]*domain.Invoice, error) {
	var models []InvoiceModel
	err := a.repo.getDB(ctx).
		Where("customer_id = ?", string(cid)).
		Order("created_at asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	for _, m := range models {
		inv, err := a.mapToDomain(m)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (a *InvoiceAdapter) mapToDomain(m InvoiceModel) (*domain.Invoice, error) {
	total, err := domain.NewMoney(m.TotalAmount, m.Currency)
	if err != nil {
		return nil, err
	}
	inv, err := domain.NewInvoice(domain.InvoiceID(m.ID), domain.CustomerID(m.CustomerID), total, parseTime(m.IssueDate), parseTime(m.DueDate))
	if err != nil {
		return nil, err
	}
	
	paid, _ := domain.NewMoney(m.PaidAmount, m.Currency)
	inv.PaidAmount = paid
	inv.Status = domain.InvoiceStatus(m.Status)
	inv.CreatedAt = parseTime(m.CreatedAt)
	inv.UpdatedAt = parseTime(m.UpdatedAt)
	
	return inv, nil
}

func (a *InvoiceAdapter) CountAllOpen(ctx context.Context) (int64, error) {
	var count int64
	err := a.repo.getDB(ctx).Model(&InvoiceModel{}).
		Where("status IN ?", []string{string(domain.InvoiceStatusOpen), string(domain.InvoiceStatusPartial)}).
		Count(&count).Error
	return count, err
}

var _ ports.InvoiceRepository = &InvoiceAdapter{}
