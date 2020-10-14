package customer_checking_account_mapping

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/fwiedmann/ezb/domain/entity"
)

func NewMySqlRepository(client *sql.DB) (*MySqlRepository, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}

	if _, err := client.Exec(`CREATE TABLE IF NOT EXISTS ezb.customerCheckinAccountMapping (CustomerID varchar(255), CheckingAccountNumber varchar(255), PRIMARY KEY (CustomerID, CheckingAccountNumber), FOREIGN KEY (CustomerID) REFERENCES customer(ID), FOREIGN KEY (CheckingAccountNumber) REFERENCES checkingAccount(Number))`); err != nil {

	}
	return &MySqlRepository{
		client: client,
	}, nil
}

type MySqlRepository struct {
	client *sql.DB
}

func (m MySqlRepository) Create(ctx context.Context, customerID, checkingAccountNumber entity.ID) error {
	_, err := m.client.ExecContext(ctx, `INSERT INTO ezb.customerCheckinAccountMapping (CustomerID, CheckingAccountNumber) VALUES (?, ?)`, customerID.String(), checkingAccountNumber.String())
	if err != nil {
		return err
	}
	return nil
}

func (m MySqlRepository) Get(ctx context.Context, customerID, checkingAccountNumber entity.ID) (Mapping, error) {
	var customerIDResp string
	var checkingAccountNumberResp string
	err := m.client.QueryRowContext(ctx, `SELECT ezb.customerCheckinAccountMapping.CustomerID, ezb.customerCheckinAccountMapping.CheckingAccountNumber FROM ezb.customerCheckinAccountMapping WHERE ezb.customerCheckinAccountMapping.CustomerID = ? AND  ezb.customerCheckinAccountMapping.CheckingAccountNumber = ?`, customerID.String(), checkingAccountNumber.String()).Scan(&customerIDResp, checkingAccountNumberResp)
	if err != nil {
		return Mapping{}, err
	}

	var mapping Mapping
	mapping.CustomerID, err = uuid.Parse(customerIDResp)
	if err != nil {
		return Mapping{}, err
	}

	mapping.CheckingAccountNumber, err = uuid.Parse(checkingAccountNumberResp)
	if err != nil {
		return Mapping{}, err
	}

	return mapping, nil
}
