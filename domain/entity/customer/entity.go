package customer

import (
	"time"

	"github.com/fwiedmann/ezb/domain/entity"
)

// Customer entity which represents unique customers
type Customer struct {
	ID                  entity.ID
	FirstName           string
	LastName            string
	Gender              string
	Birthdate           string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
}

func (c *Customer) GetCreationTimestamp() time.Time {
	return c.creationTimestamp
}

func (c *Customer) GetCLastUpdateTimestamp() time.Time {
	return c.lastUpdateTimestamp
}
