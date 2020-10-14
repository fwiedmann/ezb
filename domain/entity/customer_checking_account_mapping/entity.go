package customer_checking_account_mapping

import "github.com/fwiedmann/ezb/domain/entity"

type Mapping struct {
	CustomerID            entity.ID
	CheckingAccountNumber entity.ID
}
