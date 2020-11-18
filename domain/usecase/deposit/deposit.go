package deposit

import (
	"context"
	"errors"

	"github.com/fwiedmann/ezb/domain/entity/checking_account"

	"github.com/fwiedmann/ezb/domain/entity"
)

var (
	// ErrorDepositInvalidPinForAccount if the given pin is not equals with the stored hashed pin
	ErrorDepositInvalidPinForAccount = errors.New("invalid pin for checking account")
	// ErrorDepositInvalidAmount if the given deposit operation amount is a negative number
	ErrorDepositInvalidAmount = errors.New("invalid deposit amount. Amount has to be positive")
)

// NewUseCase inits a deposit UseCase
func NewUseCase(cm checking_account.Manager) *UseCase {
	return &UseCase{
		checkingAccountManager: cm,
	}
}

// UseCase manages deposit operations for CheckingAccounts
type UseCase struct {
	checkingAccountManager checking_account.Manager
}

// Deposit operation for the given CheckingAccount
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
