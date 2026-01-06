package sqlite

import (
	"carigo/internal/application/ports"
	"context"
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

// Helper
func parseTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// Verify interface compliance
var _ ports.TransactionManager = &GormRepository{}

// Factory to return all ports
func NewRepositories(dsn string) (*GormRepository, *CustomerAdapter, *InvoiceAdapter, *PaymentAdapter, *AllocationAdapter, error) {
	base, err := NewGormRepository(dsn)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return base, &CustomerAdapter{base}, &InvoiceAdapter{base}, &PaymentAdapter{base}, &AllocationAdapter{base}, nil
}
