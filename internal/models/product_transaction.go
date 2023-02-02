package models

import (
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type ProductTransaction struct {
	ID                   int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	TransactionID        string    `gorm:"column:transaction_id;not null" json:"transaction_id"`
	ProductTransactionID string    `gorm:"column:product_transaction_id;not null" json:"product_transaction_id"`
	Title                string    `gorm:"column:title" json:"title"`
	Quantity             int64     `gorm:"column:quantity" json:"quantity"`
	Photo                string    `gorm:"column:photo" json:"photo"`
	Amount               float64   `gorm:"column:amount" json:"amount"`
	DeletedAt            time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt            time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (p *ProductTransaction) GetAllByTransactionID(db *gorm.DB) ([]ProductTransaction, error) {
	details := []ProductTransaction{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_id = ? ", p.TransactionID)
	if err != nil {
		return details, err
	}
	return details, nil
}
