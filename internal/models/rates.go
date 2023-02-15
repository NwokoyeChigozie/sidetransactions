package models

import (
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type Rate struct {
	ID           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	FromCurrency string    `gorm:"column:from_currency;type:varchar(255); not null" json:"from_currency"`
	ToCurrency   string    `gorm:"column:to_currency;type:varchar(255); not null" json:"to_currency"`
	From_symbol  string    `gorm:"column:from_symbol;type:varchar(255); not null" json:"from_symbol"`
	ToSymbol     string    `gorm:"column:to_symbol;type:varchar(255); not null" json:"to_symbol"`
	Amount       float64   `gorm:"column:amount; type:decimal(8,5); not null" json:"amount"`
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"-"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"-"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"-"`
}

func (r Rate) GetAll(db *gorm.DB) ([]Rate, error) {
	details := []Rate{}
	err := postgresql.SelectAllFromDb(db, "desc", &details, "")
	if err != nil {
		return details, err
	}
	return details, nil
}

func (r *Rate) GetRateByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &r, "id = ?", r.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
