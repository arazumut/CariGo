package domain

type Money struct {
	amount   int64
	currency string
}

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

func (m Money) Amount() int64 {
	return m.amount
}
func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, ErrCurrencyMismatch
	}
	return m.amount > other.amount, nil
}

func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
