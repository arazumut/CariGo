package sqlite

import (
	"carigo/internal/application/ports"
	"context"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

func (r *GormRepository) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

type txKey struct{}

func (r *GormRepository) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db.WithContext(ctx)
}

func parseTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

var _ ports.TransactionManager = &GormRepository{}
func NewRepositories(dsn string) (*GormRepository, *CustomerAdapter, *InvoiceAdapter, *PaymentAdapter, *AllocationAdapter, error) {
	base, err := NewGormRepository(dsn)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return base, &CustomerAdapter{base}, &InvoiceAdapter{base}, &PaymentAdapter{base}, &AllocationAdapter{base}, nil
}
