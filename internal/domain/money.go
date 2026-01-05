package domain

// Money, floating point hatalarından kaçınmak için parayı "cents" (kuruş) cinsinden tutan Value Object'tir.
// Immutable çalışır. Her işlem yeni bir Money objesi döndürür.
type Money struct {
	amount   int64
	currency string
}

// NewMoney creates a new Money object.
func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}
	if currency == "" {
		return Money{}, ErrInvalidCurrency
	}
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount returns the raw amount in cents.
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code.
func (m Money) Currency() string {
	return m.currency
}

// Add adds another Money to this one.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// Subtract subtracts another Money from this one.
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

// IsZero returns true if the amount is 0.
func (m Money) IsZero() bool {
	return m.amount == 0
}

// GreaterThan checks if m > other.
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, ErrCurrencyMismatch
	}
	return m.amount > other.amount, nil
}

// Equals checks if m == other.
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
