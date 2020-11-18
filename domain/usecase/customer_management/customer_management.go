package customer_management

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/customer"
)

// NewUseCase inits a customer management UseCase
func NewUseCase(cm customer.Manager) *UseCase {
	return &UseCase{cm: cm}
}

// UseCase can manage the customers which are available in the ezb services
type UseCase struct {
	cm customer.Manager
}

// Create a new customer
func (u *UseCase) Create(ctx context.Context, c customer.Customer) (entity.ID, error) {
	return u.cm.Create(ctx, c)
}

// Update an existing customer
func (u *UseCase) Update(ctx context.Context, c customer.Customer) error {
	return u.cm.Update(ctx, c)
}

// Get an existing customer
func (u *UseCase) Get(ctx context.Context, id entity.ID) (customer.Customer, error) {
	return u.cm.Get(ctx, id)
}
