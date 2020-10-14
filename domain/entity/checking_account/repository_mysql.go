package checking_account

import (
	"context"
	"database/sql"
	"time"

	"github.com/fwiedmann/ezb/domain/entity"
	"github.com/google/uuid"
)

type MySqlRepository struct {
	client *sql.DB
}

func NewMySqlRepository(client *sql.DB) (*MySqlRepository, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}

	if _, err := client.Exec(`CREATE TABLE IF NOT EXISTS ezb.checkingAccount (Number varchar(255) NOT NULL, Name varchar(255) NOT NULL, Balance FLOAT NOT NULL, OverDraftLimit FLOAT NOT NULL, PinHash varchar(255), CreationTimestamp int NOT NULL, LastUpdateTimestamp int NOT NULL, PRIMARY KEY (Number))`); err != nil {

	}
	return &MySqlRepository{
		client: client,
	}, nil
}

func (m *MySqlRepository) Create(ctx context.Context, c CheckingAccount) error {
	_, err := m.client.ExecContext(ctx, `INSERT INTO ezb.checkingAccount (Number, Name, Balance, OverDraftLimit, PinHash, CreationTimestamp, LastUpdateTimestamp) VALUES (?, ?, ?, ?, ?, ?, ?)`, c.Number.String(), c.Name, c.balance, c.OverDraftLimit, c.hashedPin, c.creationTimestamp.Unix(), c.lastUpdateTimestamp.Unix())
	if err != nil {
		return err
	}
	return nil
}

func (m *MySqlRepository) Update(ctx context.Context, c CheckingAccount) error {
	_, err := m.client.ExecContext(ctx, `UPDATE ezb.checkingAccount SET Name = ?, Balance = ?, OverDraftLimit = ?, PinHash = ?, LastUpdateTimestamp = ? WHERE Number = ?`, c.Name, c.balance, c.OverDraftLimit, c.hashedPin, c.lastUpdateTimestamp.Unix(), c.Number)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySqlRepository) Get(ctx context.Context, id entity.ID) (CheckingAccount, error) {
	var c CheckingAccount
	var creationTimestamp int64
	var lastUpdateTimestamp int64
	var ID string
	err := m.client.QueryRowContext(ctx, `SELECT Number, Name, Balance, OverDraftLimit, PinHash, CreationTimestamp,LastUpdateTimestamp FROM ezb.checkingAccount WHERE ezb.checkingAccount.Number = ?`, id).Scan(&ID, &c.Name, &c.balance, &c.OverDraftLimit, &c.hashedPin, &creationTimestamp, &lastUpdateTimestamp)
	if err != nil {
		return CheckingAccount{}, err
	}
	c.Number, err = uuid.Parse(ID)
	if err != nil {
		return CheckingAccount{}, err
	}
	c.creationTimestamp = time.Unix(creationTimestamp, 0)
	c.lastUpdateTimestamp = time.Unix(lastUpdateTimestamp, 0)
	return c, nil
}
