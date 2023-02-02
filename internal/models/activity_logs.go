package models

import (
	"fmt"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type ActivityLog struct {
	ID            uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	TransactionID string    `gorm:"column:transaction_id; type:varchar(255); not null; comment: 12 characters long string" json:"transaction_id"`
	Description   string    `gorm:"column:description; type:varchar(255); not null" json:"description"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (a *ActivityLog) GetAllByTransactionID(db *gorm.DB) ([]ActivityLog, error) {
	details := []ActivityLog{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_id = ? ", a.TransactionID)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (a *ActivityLog) CreateActivityLog(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &a)
	if err != nil {
		return fmt.Errorf("ActivityLog creation failed: %v", err.Error())
	}
	return nil
}
