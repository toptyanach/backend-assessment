package transfer

import (
	"context"
	"errors"
	"github.com/toptyanach/backend-assessment/contracts"
	"github.com/toptyanach/backend-assessment/domain"
	"testing"
)

type mockRepo struct {
	accounts map[domain.AccountID]*domain.Account
}

func (m *mockRepo) Retrieve(ctx context.Context, id domain.AccountID) (*domain.Account, error) {
	acc, exists := m.accounts[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (m *mockRepo) UpdateMut(account *domain.Account) *contracts.Mutation {
	if len(account.Changes) == 0 {
		return nil
	}
	return &contracts.Mutation{
		Table: "accounts",
		ID:    string(account.ID()),
	}
}

func TestInteractor_Execute(t *testing.T) {
	tests := []struct {
		name        string
		req         *TransferRequest
		setupRepo   func() contracts.AccountRepository
		expectedErr error
		expectedMut int
	}{
		{
			name: "Success transfer",
			req: &TransferRequest{
				FromAccountID: "acc1",
				ToAccountID:   "acc2",
				Amount:        100,
			},
			setupRepo: func() contracts.AccountRepository {
				return &mockRepo{
					accounts: map[domain.AccountID]*domain.Account{
						"acc1": domain.NewAccount("acc1", 500, domain.StatusActive),
						"acc2": domain.NewAccount("acc2", 100, domain.StatusActive),
					},
				}
			},
			expectedErr: nil,
			expectedMut: 2,
		},
		{
			name: "Insufficient funds",
			req: &TransferRequest{
				FromAccountID: "acc1",
				ToAccountID:   "acc2",
				Amount:        1000,
			},
			setupRepo: func() contracts.AccountRepository {
				return &mockRepo{
					accounts: map[domain.AccountID]*domain.Account{
						"acc1": domain.NewAccount("acc1", 500, domain.StatusActive),
						"acc2": domain.NewAccount("acc2", 100, domain.StatusActive),
					},
				}
			},
			expectedErr: domain.ErrInsufficientFunds,
			expectedMut: 0,
		},
		{
			name: "Transfer to same account",
			req: &TransferRequest{
				FromAccountID: "acc1",
				ToAccountID:   "acc1",
				Amount:        100,
			},
			setupRepo: func() contracts.AccountRepository {
				return &mockRepo{}
			},
			expectedErr: errors.New("cannot transfer to the same account"),
			expectedMut: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewInteractor(tt.setupRepo())
			plan, err := uc.Execute(context.Background(), tt.req)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			if err == nil {
				if len(plan.GetMutations()) != tt.expectedMut {
					t.Errorf("expected %d mutations, got %d", tt.expectedMut, len(plan.GetMutations()))
				}
			}
		})
	}
}
