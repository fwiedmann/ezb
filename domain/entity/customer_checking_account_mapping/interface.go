package customer_checking_account_mapping

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
)

type Manager interface {
	Create(ctx context.Context, customerID, checkingAccountNumber entity.ID) error
	Get(ctx context.Context, customerID, checkingAccountNumber entity.ID) (Mapping, error)
}

type repository interface {
	Manager
}
