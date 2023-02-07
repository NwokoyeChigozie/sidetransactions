package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionState struct {
	ID            int64     `gorm:"primary_key;column:id" json:"id"`
	AccountID     int64     `gorm:"column:account_id; comment: User account id of the user that resulted in the action" json:"account_id"`
	TransactionID string    `gorm:"column:transaction_id" json:"transaction_id"`
	MilestoneID   string    `gorm:"column:milestone_id; comment: Milestone ID of the transaction if serviced based" json:"milestone_id"`
	Status        string    `gorm:"column:status" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (t *TransactionState) GetTransactionStateByTransactionIDAndStatus(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ? and status=?", t.TransactionID, t.Status)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (t *TransactionState) CreateTransactionState(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction state creation failed: %v", err.Error())
	}
	return nil
}
