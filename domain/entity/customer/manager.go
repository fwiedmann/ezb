package customer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/fwiedmann/ezb/domain/entity"
)

// birthdateLayout represents a birthdate in this format: day-month-year
const birthdateLayout = "02-01-2006"

var (
	ErrorEmptyUserInput   = errors.New("user input contains empty fields")
	ErrorInvalidBirthdate = errors.New("users birthdate is invalid. It should be in the format of day-month-year")
)

func NewManager(r repository) Manager {
	return manager{repo: r}
}

type manager struct {
	repo repository
}

func (m manager) Create(ctx context.Context, c Customer) (entity.ID, error) {
	if err := validateCustomer(c); err != nil {
		return [16]byte{}, err
	}
	c.ID = entity.NewID()
	timestamp := time.Now()
	c.creationTimestamp, c.lastUpdateTimestamp = timestamp, timestamp

	if err := m.repo.Create(ctx, c); err != nil {
		return [16]byte{}, err
	}
	return c.ID, nil
}

func (m manager) Update(ctx context.Context, c Customer) error {
	if err := validateCustomer(c); err != nil {
		return err
	}

	c.lastUpdateTimestamp = time.Now()

	if err := m.repo.Update(ctx, c); err != nil {
		return err
	}
	return nil
}

func (m manager) Get(ctx context.Context, id entity.ID) (Customer, error) {
	return m.repo.Get(ctx, id)
}

func validateCustomer(c Customer) error {
	if c.FirstName == "" {
		return ErrorEmptyUserInput
	}

	if c.LastName == "" {
		return ErrorEmptyUserInput
	}

	if c.Gender == "" {
		return ErrorEmptyUserInput
	}

	if c.Birthdate == "" {
		return ErrorEmptyUserInput
	}

	if err := validateBirthdate(c.Birthdate); err != nil {
		return err
	}
	return nil
}

func hashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), 14)
	return string(bytes), err
}

func validateBirthdate(date string) error {
	if _, err := time.Parse(birthdateLayout, date); err != nil {
		return fmt.Errorf("%w, error: %s", ErrorInvalidBirthdate, err)
	}
	return nil
}
