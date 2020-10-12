package checking_account

import (
	"time"

	"github.com/fwiedmann/ezb/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type CheckingAccount struct {
	Number              entity.ID
	Name                string
	OverDraftLimit      float64
	balance             float64
	hashedPin           string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
}

func (c *CheckingAccount) IsValidPin(pin string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.hashedPin), []byte(pin)) == nil
}

func (c *CheckingAccount) GetCurrentBalance() float64 {
	return c.balance
}

func (c *CheckingAccount) IsDebitAllowed(debitValue float64) bool {
	return c.balance-debitValue > c.OverDraftLimit
}

func (c *CheckingAccount) Debit(debitValue float64) {
	c.balance -= debitValue
}

func (c *CheckingAccount) Deposit(debitValue float64) {
	c.balance += debitValue
}

func (c *CheckingAccount) GetCreationTimestamp() time.Time {
	return c.creationTimestamp
}

func (c *CheckingAccount) GetCLastUpdateTimestamp() time.Time {
	return c.lastUpdateTimestamp
}
