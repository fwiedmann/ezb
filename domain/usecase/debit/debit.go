package debit

import (
	"context"
	"errors"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/checking_account"
)

var (
	ErrorDebitInvalidPinForAccount  = errors.New("invalid pin for checking account")
	ErrorDebitInvalidAmount         = errors.New("invalid deposit amount. Amount has to be positive")
	ErrorDebitExceedsOverdraftLimit = errors.New("debit amount exceeds the overdraft limit")
)

func NewUseCase(cm checking_account.Manager) *UseCase {
	return &UseCase{
		checkingAccountManager: cm,
	}
}

type UseCase struct {
	checkingAccountManager checking_account.Manager
}

func (uc *UseCase) Debit(ctx context.Context, checkingAccountNumber entity.ID, amount float64, pin string) error {
	if amount <= 0 {
		return ErrorDebitInvalidAmount
	}

	account, err := uc.checkingAccountManager.Get(ctx, checkingAccountNumber)
	if err != nil {
		return err
	}

	if !account.IsValidPin(pin) {
		return ErrorDebitInvalidPinForAccount
	}

	if !account.IsDebitAllowed(amount) {
		return ErrorDebitExceedsOverdraftLimit
	}
	account.Debit(amount)
	return uc.checkingAccountManager.Update(ctx, account, pin)
}
