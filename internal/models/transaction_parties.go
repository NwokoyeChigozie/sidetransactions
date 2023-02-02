package models

import (
	"fmt"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionParty struct {
	ID                   uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	TransactionPartiesID string    `gorm:"column:transaction_parties_id; type:varchar(255); not null" json:"transaction_parties_id"`
	TransactionID        string    `gorm:"column:transaction_id; type:varchar(255); not null; comment: 12 characters long string" json:"transaction_id"`
	AccountID            int       `gorm:"column:account_id; type:int" json:"account_id"`
	Role                 string    `gorm:"column:role; type:varchar(255); not null" json:"role"`
	DeletedAt            time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt            time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	RoleCapabilities     jsonmap   `gorm:"column:role_capabilities; type:varchar(250); default: '{\"view\"}'; comment: view|manage" json:"role_capabilities"`
	RoleDescription      string    `gorm:"column:role_description; type:text" json:"role_description"`
	Status               string    `gorm:"column:status; type:varchar(255); not null;default:created" json:"status"`
}

func (t *TransactionParty) CreateTransactionParty(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction party creation failed: %v", err.Error())
	}
	return nil
}

func (t *TransactionParty) GetAllByTransactionID(db *gorm.DB) ([]TransactionParty, error) {
	details := []TransactionParty{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_id = ? ", t.TransactionID)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (t *TransactionParty) GetAllByTransactionPartiesID(db *gorm.DB) ([]TransactionParty, error) {
	details := []TransactionParty{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_parties_id = ? ", t.TransactionPartiesID)
	if err != nil {
		return details, err
	}
	return details, nil
}
