package assign_customer_to_checking_account

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/customer_checking_account_mapping"
)

func NewUseCase(mappingManager customer_checking_account_mapping.Manager) *UseCase {
	return &UseCase{mappingManager: mappingManager}
}

type UseCase struct {
	mappingManager customer_checking_account_mapping.Manager
}

func (uc *UseCase) CreateMapping(ctx context.Context, customerID, checkingAccountNumber entity.ID) error {
	return uc.mappingManager.Create(ctx, customerID, checkingAccountNumber)
}
