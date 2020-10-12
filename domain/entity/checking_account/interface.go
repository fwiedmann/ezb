package checking_account

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
)

type Manager interface {
	Create(ctx context.Context, c CheckingAccount, pin string) (entity.ID, error)
	Update(ctx context.Context, c CheckingAccount, pin string) error
	Get(ctx context.Context, id entity.ID) (CheckingAccount, error)
}

type repository interface {
	Create(ctx context.Context, c CheckingAccount) error
	Update(ctx context.Context, c CheckingAccount) error
	Get(ctx context.Context, id entity.ID) (CheckingAccount, error)
}
