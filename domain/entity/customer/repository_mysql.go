package customer

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/fwiedmann/ezb/domain/entity"
)

type MySqlRepository struct {
	client *sql.DB
}

func NewMySqlRepository(client *sql.DB) (*MySqlRepository, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}

	if _, err := client.Exec(`CREATE TABLE IF NOT EXISTS ezb.customer (ID varchar(255) NOT NULL, FirstName varchar(255) NOT NULL, LastName varchar(255) NOT NULL, Birthdate varchar(255) NOT NULL, Gender varchar(255), CreationTimestamp int NOT NULL, LastUpdateTimestamp int NOT NULL, PRIMARY KEY (ID))`); err != nil {

	}
	return &MySqlRepository{
		client: client,
	}, nil
}

func (m MySqlRepository) Create(ctx context.Context, c Customer) error {
	_, err := m.client.ExecContext(ctx, `INSERT INTO ezb.customer (ID, FirstName, LastName, Birthdate, Gender, CreationTimestamp, LastUpdateTimestamp) VALUES (?, ?, ?, ?, ?, ?, ?)`, c.ID.String(), c.FirstName, c.LastName, c.Birthdate, c.Gender, c.creationTimestamp.Unix(), c.lastUpdateTimestamp.Unix())
	if err != nil {
		return err
	}
	return nil
}

func (m MySqlRepository) Update(ctx context.Context, c Customer) error {
	_, err := m.client.ExecContext(ctx, `UPDATE ezb.customer SET FirstName = ?, LastName = ?, Birthdate = ?, Gender = ?, LastUpdateTimestamp = ? WHERE ID = ?`, c.FirstName, c.LastName, c.Birthdate, c.Gender, c.lastUpdateTimestamp.Unix(), c.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m MySqlRepository) Get(ctx context.Context, id entity.ID) (Customer, error) {
	var c Customer
	var creationTimestamp int64
	var lastUpdateTimestamp int64
	var ID string
	err := m.client.QueryRowContext(ctx, `SELECT ezb.customer.ID, ezb.customer.FirstName, ezb.customer.LastName, ezb.customer.Birthdate, ezb.customer.Gender, ezb.customer.CreationTimestamp, ezb.customer.LastUpdateTimestamp FROM ezb.customer WHERE ezb.customer.ID = ?`, id).Scan(&ID, &c.FirstName, &c.LastName, &c.Birthdate, &c.Gender, &creationTimestamp, &lastUpdateTimestamp)
	if err != nil {
		return Customer{}, err
	}
	c.ID, err = uuid.Parse(ID)
	if err != nil {
		return Customer{}, err
	}
	c.creationTimestamp = time.Unix(creationTimestamp, 0)
	c.lastUpdateTimestamp = time.Unix(lastUpdateTimestamp, 0)
	return c, nil
}

// c0b87604-2442-4d19-96d8-e44e026d3481
