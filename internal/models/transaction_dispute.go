package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionDispute struct {
	ID            int64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	DisputeID     string    `gorm:"column:dispute_id" json:"dispute_id"`
	TransactionID string    `gorm:"column:transaction_id" json:"transaction_id"`
	Reason        string    `gorm:"column:reason" json:"reason"`
	DisputeStatus string    `gorm:"column:dispute_status" json:"dispute_status"`
	Decision      string    `gorm:"column:decision" json:"decision"`
	MediatorID    string    `gorm:"column:mediator_id" json:"mediator_id"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type CreateDisputeRequest struct {
	TransactionID string `json:"transaction_id" validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	Reason        string `json:"reason"`
	DisputeStatus string `json:"dispute_status"`
	Decision      string `json:"decision"`
}

func (t *TransactionDispute) CreateTransactionDispute(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction dispute creation failed: %v", err.Error())
	}
	return nil
}

func (t *TransactionDispute) IsDisputed(db *gorm.DB) (bool, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ?", t.TransactionID)
	if nilErr != nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TransactionDispute) GetTransactionDisputeByTransactionID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ?", t.TransactionID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (t *TransactionDispute) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &t)
	return err
}
