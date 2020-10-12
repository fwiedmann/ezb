package checking_account

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/fwiedmann/ezb/domain/entity"
)

// birthdateLayout represents a birthdate in this format: day-month-year
const birthdateLayout = "02-01-2006"

var (
	ErrorEmptyUserInput        = errors.New("user input contains empty fields")
	ErrorInvalidOverdraftLimit = errors.New("overdraft limit has to be a negative number")
)

func NewManager(r repository) Manager {
	return &manager{repo: r}
}

type manager struct {
	repo repository
}

func (m *manager) Create(ctx context.Context, c CheckingAccount, pin string) (entity.ID, error) {
	if err := validateCheckingAccount(c); err != nil {
		return [16]byte{}, err
	}

	hashedPin, err := hashPin(pin)
	if err != nil {
		return [16]byte{}, err
	}
	c.hashedPin = hashedPin
	c.Number = entity.NewID()

	timestamp := time.Now()
	c.creationTimestamp, c.lastUpdateTimestamp = timestamp, timestamp
	if err := m.repo.Create(ctx, c); err != nil {
		return [16]byte{}, err
	}
	return c.Number, nil
}

func (m *manager) Update(ctx context.Context, c CheckingAccount, pin string) error {
	if err := validateCheckingAccount(c); err != nil {
		return err
	}

	hashedPin, err := hashPin(pin)
	if err != nil {
		return err
	}
	c.hashedPin = hashedPin
	c.lastUpdateTimestamp = time.Now()
	return m.repo.Update(ctx, c)
}

func (m *manager) Get(ctx context.Context, id entity.ID) (CheckingAccount, error) {
	return m.repo.Get(ctx, id)
}

func validateCheckingAccount(c CheckingAccount) error {
	if c.Name == "" {
		return ErrorEmptyUserInput
	}

	if c.OverDraftLimit > 0 {
		return ErrorInvalidOverdraftLimit
	}
	return nil
}

func hashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), 14)
	return string(bytes), err
}
