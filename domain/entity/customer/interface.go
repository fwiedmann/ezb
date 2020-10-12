package customer

import (
	"context"

	"github.com/fwiedmann/ezb/domain/entity"
)

type Manager interface {
	Create(ctx context.Context, c Customer) (entity.ID, error)
	Update(ctx context.Context, c Customer) error
	Get(ctx context.Context, id entity.ID) (Customer, error)
}

type repository interface {
	Create(ctx context.Context, c Customer) error
	Update(ctx context.Context, c Customer) error
	Get(ctx context.Context, id entity.ID) (Customer, error)
}
