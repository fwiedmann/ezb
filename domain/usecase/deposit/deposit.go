package deposit

import (
	"context"
	"errors"

	"github.com/fwiedmann/ezb/domain/entity/checking_account"

	"github.com/fwiedmann/ezb/domain/entity"
)

var (
	ErrorDepositInvalidPinForAccount = errors.New("invalid pin for checking account")
	ErrorDepositInvalidAmount        = errors.New("invalid deposit amount. Amount has to be positive")
)

func NewUseCase(cm checking_account.Manager) *UseCase {
	return &UseCase{
		checkingAccountManager: cm,
	}
}

type UseCase struct {
	checkingAccountManager checking_account.Manager
}

func (uc *UseCase) Deposit(ctx context.Context, checkingAccountNumber entity.ID, amount float64, pin string) error {

	if amount <= 0 {
		return ErrorDepositInvalidAmount
	}

	account, err := uc.checkingAccountManager.Get(ctx, checkingAccountNumber)
	if err != nil {
		return err
	}

	if !account.IsValidPin(pin) {
		return ErrorDepositInvalidPinForAccount
	}

	account.Deposit(amount)
	return uc.checkingAccountManager.Update(ctx, account, pin)
}
