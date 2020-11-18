package checking_account_management

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/checking_account"
)

// NewUseCase inits a new checking account management usecase
func NewUseCase(cm checking_account.Manager) *UseCase {
	return &UseCase{cm: cm}
}

// UseCase manages the ezb available CheckinAccount
type UseCase struct {
	cm checking_account.Manager
}

// Create a new CheckingAccount
func (u *UseCase) Create(ctx context.Context, c checking_account.CheckingAccount, pin string) (entity.ID, error) {
	return u.cm.Create(ctx, c, pin)
}

// Update an existing CheckingAccount
func (u *UseCase) Update(ctx context.Context, c checking_account.CheckingAccount, pin string) error {
	return u.cm.Update(ctx, c, pin)
}

// Get an existing CheckingAccount
func (u *UseCase) Get(ctx context.Context, id entity.ID) (checking_account.CheckingAccount, error) {
	return u.cm.Get(ctx, id)
}
