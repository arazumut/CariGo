package sqlite

import (
	"carigo/internal/application/ports"
	"carigo/internal/domain"
	"context"
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormRepository is a base struct that implements TransactionManager and holds DB connection.
type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(dsn string) (*GormRepository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	
	// Auto Migrate for MVP speed (In prod, use strict migrations)
	err = db.AutoMigrate(
		&CustomerModel{},
		&InvoiceModel{},
		&PaymentModel{},
		&AllocationModel{},
	)
	if err != nil {
		return nil, err
	}

	return &GormRepository{db: db}, nil
}

// TransactionManager Implementation
func (r *GormRepository) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject tx into context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// Helper to extract DB/TX from context
type txKey struct{}

func (r *GormRepository) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

// ---------------------------------------------------------
// REPOSITORY IMPLEMENTATIONS (Combining into one file for phase 1 speed, usually separate)
// ---------------------------------------------------------

// --- MODELS (Infra specific representations) ---

type CustomerModel struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Email     string
	TaxID     string
	CreatedAt int64
	UpdatedAt int64
}

type InvoiceModel struct {
	ID             string `gorm:"primaryKey"`
	CustomerID     string `gorm:"index"`
	TotalAmount    int64
	Currency       string
	PaidAmount     int64
	Status         string
	IssueDate      int64
	DueDate        int64
	CreatedAt      int64
	UpdatedAt      int64
}

type PaymentModel struct {
	ID              string `gorm:"primaryKey"`
	CustomerID      string `gorm:"index"`
	Amount          int64
	Currency        string
	AvailableAmount int64
	Date            int64
	CreatedAt       int64
}

type AllocationModel struct {
	ID        string `gorm:"primaryKey"`
	PaymentID string `gorm:"index"`
	InvoiceID string `gorm:"index"`
	Amount    int64
	Currency  string
	CreatedAt int64
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
	// Note: created_at/updated_at mapping omitted for brevity, logic usually doesn't need them critical
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

func (r *GormRepository) FindOpenByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Invoice, error) {
	var models []InvoiceModel
	// Find invoices where status is OPEN or PARTIAL, ordered by DueDate ASC (FIFO)
	err := r.getDB(ctx).
		Where("customer_id = ? AND status IN ?", string(customerID), []string{string(domain.InvoiceStatusOpen), string(domain.InvoiceStatusPartial)}).
		Order("due_date asc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	for _, m := range models {
		total, _ := domain.NewMoney(m.TotalAmount, m.Currency)
		// Reconstruct entity (simplified, normally use a factory or mapper)
		inv, _ := domain.NewInvoice(domain.InvoiceID(m.ID), domain.CustomerID(m.CustomerID), total, parseTime(m.IssueDate), parseTime(m.DueDate))
		
		// Set internal state
		paid, _ := domain.NewMoney(m.PaidAmount, m.Currency)
		// We have to use reflection or modify entity to set private/protected fields if we strictly follow Go struct rules.
		// Or we trust the domain logic is "NewInvoice creates initial state", but here we are loading from DB.
		// For MVP, we will assume we can reconstruct or simple assign if fields were public (they are public in our struct).
		// Wait, our domain fields are public. Invoice fields are detailed in domain/invoice.go.
		inv.PaidAmount = paid
		inv.Status = domain.InvoiceStatus(m.Status)
		inv.CreatedAt = parseTime(m.CreatedAt)
		inv.UpdatedAt = parseTime(m.UpdatedAt)
		
		invoices = append(invoices, inv)
	}
	return invoices, nil
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

// --- Allocation Repo ---

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

// Helper
func parseTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// Verify interface compliance
var _ ports.InvoiceRepository = &InvoiceAdapter{}
var _ ports.PaymentRepository = &PaymentAdapter{}
var _ ports.CustomerRepository = &CustomerAdapter{}
var _ ports.AllocationRepository = &AllocationAdapter{}
var _ ports.TransactionManager = &GormRepository{}

// Adapers to separate Repository Structs if needed, or simple direct methods on GormRepository.
// Since Go doesn't support method overloading and we want "r.invoiceRepo.Save", we need adapters or separate structs sharing the DB.

type CustomerAdapter struct{ repo *GormRepository }
func (a *CustomerAdapter) Save(ctx context.Context, c *domain.Customer) error { return a.repo.SaveCustomer(ctx, c) }
func (a *CustomerAdapter) FindByID(ctx context.Context, id domain.CustomerID) (*domain.Customer, error) { return a.repo.FindCustomerByID(ctx, id) }

type InvoiceAdapter struct{ repo *GormRepository }
func (a *InvoiceAdapter) Save(ctx context.Context, i *domain.Invoice) error { return a.repo.SaveInvoice(ctx, i) }
func (a *InvoiceAdapter) FindByID(ctx context.Context, id domain.InvoiceID) (*domain.Invoice, error) { return nil, errors.New("not implemented") } // Not needed for current UseCase
func (a *InvoiceAdapter) FindOpenByCustomer(ctx context.Context, cid domain.CustomerID) ([]*domain.Invoice, error) { return a.repo.FindOpenByCustomer(ctx, cid) }
func (a *InvoiceAdapter) CountAllOpen(ctx context.Context) (int64, error) {
	var count int64
	err := a.repo.getDB(ctx).Model(&InvoiceModel{}).
		Where("status IN ?", []string{string(domain.InvoiceStatusOpen), string(domain.InvoiceStatusPartial)}).
		Count(&count).Error
	return count, err
}

type PaymentAdapter struct{ repo *GormRepository }
func (a *PaymentAdapter) Save(ctx context.Context, p *domain.Payment) error { return a.repo.SavePayment(ctx, p) }
func (a *PaymentAdapter) FindByID(ctx context.Context, id domain.PaymentID) (*domain.Payment, error) { return nil, errors.New("not implemented") }
func (a *PaymentAdapter) SumTotalCollected(ctx context.Context) (int64, error) {
	var total int64
	// Sum the initial amounts of all payments
	err := a.repo.getDB(ctx).Model(&PaymentModel{}).
		Select("ifnull(sum(amount), 0)").
		Scan(&total).Error
	return total, err
}

type AllocationAdapter struct{ repo *GormRepository }
func (a *AllocationAdapter) Save(ctx context.Context, al *domain.Allocation) error { return a.repo.SaveAllocation(ctx, al) }

// Factory to return all ports
func NewRepositories(dsn string) (*GormRepository, *CustomerAdapter, *InvoiceAdapter, *PaymentAdapter, *AllocationAdapter, error) {
	base, err := NewGormRepository(dsn)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return base, &CustomerAdapter{base}, &InvoiceAdapter{base}, &PaymentAdapter{base}, &AllocationAdapter{base}, nil
}
