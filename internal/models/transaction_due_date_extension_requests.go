package models

import (
	"fmt"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionDueDateExtensionRequest struct {
	ID            int64     `gorm:"primary_key;column:id" json:"id"`
	AccountID     int64     `gorm:"column:account_id; not null" json:"account_id"`
	TransactionID string    `gorm:"column:transaction_id; type:varchar(255); not null" json:"transaction_id"`
	Note          string    `gorm:"column:note; type:text; not null" json:"reason"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (t *TransactionDueDateExtensionRequest) CreateTransactionDueDateExtensionRequest(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction due date extension request creation failed: %v", err.Error())
	}
	return nil
}
