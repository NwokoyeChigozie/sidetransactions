package models

import (
	"fmt"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionFile struct {
	ID            uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	TransactionID string    `gorm:"column:transaction_id; type:varchar(255); not null; comment: 12 characters long string" json:"transaction_id"`
	AccountID     int       `gorm:"column:account_id; type:int" json:"account_id"`
	FileType      string    `gorm:"column:file_type; type:varchar(255)" json:"file_type"`
	FileUrl       string    `gorm:"column:file_url; type:varchar(255); not null" json:"file_url"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (t *TransactionFile) CreateTransactionFile(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction file creation failed: %v", err.Error())
	}
	return nil
}

func (t *TransactionFile) GetAllByTransactionID(db *gorm.DB) ([]TransactionFile, error) {
	details := []TransactionFile{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_id = ? ", t.TransactionID)
	if err != nil {
		return details, err
	}
	return details, nil
}
