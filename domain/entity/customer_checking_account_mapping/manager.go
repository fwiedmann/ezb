package customer_checking_account_mapping

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
)

func NewManager(r repository) Manager {
	return manager{repository: r}
}

type manager struct {
	repository repository
}

func (m manager) Get(ctx context.Context, customerID, checkingAccountNumber entity.ID) (Mapping, error) {
	return m.repository.Get(ctx, customerID, checkingAccountNumber)
}

func (m manager) Create(ctx context.Context, customerID, checkingAccountNumber entity.ID) error {
	return m.repository.Create(ctx, customerID, checkingAccountNumber)
}
