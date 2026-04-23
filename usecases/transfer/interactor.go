package transfer

import (
	"context"
	"errors"
	"fmt"
	"github.com/toptyanach/backend-assessment/contracts"
	"github.com/toptyanach/backend-assessment/domain"
)

type TransferRequest struct {
	FromAccountID domain.AccountID
	ToAccountID   domain.AccountID
	Amount        int64
}

type Interactor struct {
	repo contracts.AccountRepository
}

func NewInteractor(repo contracts.AccountRepository) *Interactor {
	return &Interactor{repo: repo}
}

func (uc *Interactor) Execute(ctx context.Context, req *TransferRequest) (*contracts.Plan, error) {
	if req.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}
	if req.FromAccountID == req.ToAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	source, err := uc.repo.Retrieve(ctx, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("retrieve source account: %w", err)
	}

	dest, err := uc.repo.Retrieve(ctx, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("retrieve destination account: %w", err)
	}

	if err := source.Withdraw(req.Amount); err != nil {
		return nil, err
	}

	if err := dest.Deposit(req.Amount); err != nil {
		return nil, err
	}

	plan := contracts.NewPlan()
	plan.Add(uc.repo.UpdateMut(source))
	plan.Add(uc.repo.UpdateMut(dest))

	return plan, nil
}
