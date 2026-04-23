package domain

import "errors"

type AccountID string
type AccountStatus string

const (
	StatusActive AccountStatus = "ACTIVE"
	StatusLocked AccountStatus = "LOCKED"
	StatusClosed AccountStatus = "CLOSED"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountInactive   = errors.New("account is not active")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
)

type ChangeTracker map[string]struct{}

func (c ChangeTracker) Mark(field string)     { c[field] = struct{}{} }
func (c ChangeTracker) Has(field string) bool { _, ok := c[field]; return ok }

type Account struct {
	id      AccountID
	balance int64 // в центах
	status  AccountStatus
	Changes ChangeTracker
}

func NewAccount(id AccountID, balance int64, status AccountStatus) *Account {
	return &Account{
		id:      id,
		balance: balance,
		status:  status,
		Changes: make(ChangeTracker),
	}
}

func (a *Account) ID() AccountID         { return a.id }
func (a *Account) Balance() int64        { return a.balance }
func (a *Account) Status() AccountStatus { return a.status }

func (a *Account) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.status != StatusActive {
		return ErrAccountInactive
	}
	if a.balance < amount {
		return ErrInsufficientFunds
	}

	a.balance -= amount
	a.Changes.Mark("balance")
	return nil
}

func (a *Account) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.status != StatusActive {
		return ErrAccountInactive
	}

	a.balance += amount
	a.Changes.Mark("balance")
	return nil
}
