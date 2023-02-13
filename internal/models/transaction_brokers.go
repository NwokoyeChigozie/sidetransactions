package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type TransactionBroker struct {
	ID                  uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	TransactionBrokerID string    `gorm:"column:transaction_broker_id; type:varchar(255); not null" json:"transaction_broker_id"`
	TransactionID       string    `gorm:"column:transaction_id; type:varchar(255); not null; comment: 12 characters long string" json:"transaction_id"`
	BrokerCharge        string    `gorm:"column:broker_charge; type:varchar(255); not null" json:"broker_charge"`
	BrokerChargeBearer  string    `gorm:"column:broker_charge_bearer; type:varchar(255); not null" json:"broker_charge_bearer"`
	CreatedAt           time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	BrokerChargeType    string    `gorm:"column:broker_charge_type; type:varchar(255); not null;default:fixed; comment: fixed|percentage" json:"broker_charge_type"`
	IsSellerAccepted    bool      `gorm:"column:is_seller_accepted; type:bool; default:false" json:"is_seller_accepted"`
	IsBuyerAccepted     bool      `gorm:"column:is_buyer_accepted; type:bool; default:false" json:"is_buyer_accepted"`
}

type UpdateTransactionBrokerRequest struct {
	TransactionID      string `json:"transaction_id" validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	BrokerCharge       string `json:"broker_charge"`
	BrokerChargeBearer string `json:"broker_charge_bearer"`
	BrokerChargeType   string `json:"broker_charge_type"`
	IsSellerAccepted   *bool  `json:"is_seller_accepted"`
	IsBuyerAccepted    *bool  `json:"is_buyer_accepted"`
}

func (t *TransactionBroker) GetTransactionBrokerByTransactionID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ?", t.TransactionID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (t *TransactionBroker) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &t)
	return err
}

func (t *TransactionBroker) CreateTransactionBroker(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction broker creation failed: %v", err.Error())
	}
	return nil
}
