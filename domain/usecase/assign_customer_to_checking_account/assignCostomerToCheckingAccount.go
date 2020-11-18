package assign_customer_to_checking_account

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/fwiedmann/ezb/domain/entity/customer_checking_account_mapping"
)

// NewUseCase inits the useCase manager to create mapping of CheckingAccounts with Customers
func NewUseCase(mappingManager customer_checking_account_mapping.Manager) *UseCase {
	return &UseCase{mappingManager: mappingManager}
}

// UseCase manages mappings of CheckingAccounts with Customers
type UseCase struct {
	mappingManager customer_checking_account_mapping.Manager
}

// CreateMapping for the given CheckingAccount and Customer
func (uc *UseCase) CreateMapping(ctx context.Context, customerID, checkingAccountNumber entity.ID) error {
	return uc.mappingManager.Create(ctx, customerID, checkingAccountNumber)
}
