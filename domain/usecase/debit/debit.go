package debit

import (
	"context"
	"errors"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/checking_account"
)

var (
	// ErrorDebitInvalidPinForAccount if the given pin is not equals with the stored hashed pin
	ErrorDebitInvalidPinForAccount = errors.New("invalid pin for checking account")
	// ErrorDebitInvalidAmount if the given debit operation amount is a negative number
	ErrorDebitInvalidAmount = errors.New("invalid debit amount. Amount has to be positive")
	// ErrorDebitExceedsOverdraftLimit if the given debit operation would exceed the overall checkingAccount balance
	ErrorDebitExceedsOverdraftLimit = errors.New("debit amount exceeds the overdraft limit")
)

// NewUseCase inits a new debit UseCase
func NewUseCase(cm checking_account.Manager) *UseCase {
	return &UseCase{
		checkingAccountManager: cm,
	}
}

// UseCase manages all debit operations
type UseCase struct {
	checkingAccountManager checking_account.Manager
}

// Debit operation for the given CheckingAccount
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
