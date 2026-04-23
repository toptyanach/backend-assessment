package repo

import (
	"github.com/toptyanach/backend-assessment/contracts"
	"github.com/toptyanach/backend-assessment/domain"
)

type AccountRepo struct{}

func (r *AccountRepo) UpdateMut(account *domain.Account) *contracts.Mutation {
	if account == nil || len(account.Changes) == 0 {
		return nil
	}

	updates := make(map[string]interface{})

	if account.Changes.Has("balance") {
		updates["balance"] = account.Balance()
	}
	if account.Changes.Has("status") {
		updates["status"] = account.Status()
	}

	return &contracts.Mutation{
		Table:   "accounts",
		ID:      string(account.ID()),
		Updates: updates,
	}
}
