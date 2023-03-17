package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
	"gorm.io/gorm"
)

type ExchangeTransaction struct {
	ID            uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	AccountID     string    `gorm:"column:account_id; type:varchar(255); not null" json:"account_id"`
	InitialAmount float64   `gorm:"column:initial_amount; type:decimal(8,5)" json:"initial_amount"`
	FinalAmount   float64   `gorm:"column:final_amount; type:decimal(8,2)" json:"final_amount"`
	RateID        int       `gorm:"column:rate_id; type:int; not null" json:"rate_id"`
	Status        string    `gorm:"column:status; type:varchar(255); not null; default: pending; comment: failed,pending,completed" json:"status"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"-"`
	CreatedAt     time.Time `gorm:"column:created_at; autoCreateTime" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at; autoUpdateTime" json:"-"`
}
type ExchangeTransactionWithRate struct {
	ID              uint    `json:"id"`
	AccountID       string  `json:"account_id"`
	InitialAmount   float64 `json:"initial_amount"`
	FinalAmount     float64 `json:"final_amount"`
	Rate            Rate    `json:"rate"`
	Status          string  `json:"status"`
	TransactionName string  `json:"transaction_name"`
	Date            string  `json:"date"`
}
type CreateExchangeTransactionRequest struct {
	AccountID     int     `json:"account_id" validate:"required" pgvalidate:"exists=auth$users$account_id"`
	InitialAmount float64 `json:"initial_amount" validate:"required"`
	FinalAmount   float64 `json:"final_amount" validate:"required"`
	RateID        int     `json:"rate_id" validate:"required" pgvalidate:"exists=transaction$rates$id"`
	Status        string  `json:"status" validate:"required,oneof=failed pending completed"`
}

func (t *ExchangeTransaction) CreateExchangeTransaction(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("exchange transaction creation failed: %v", err.Error())
	}
	return nil
}

func (e *ExchangeTransaction) GetAllByAccountID(db *gorm.DB) ([]ExchangeTransaction, error) {
	details := []ExchangeTransaction{}
	err := postgresql.SelectAllFromDb(db, "desc", &details, "account_id = ?", e.AccountID)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (e *ExchangeTransaction) GetExchangeTransactionByID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &e, "id = ?", e.ID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (e *ExchangeTransaction) GetAllResolvedByAccountID(db *gorm.DB, paginator postgresql.Pagination) ([]ExchangeTransactionWithRate, postgresql.PaginationResponse, error) {
	details := []ExchangeTransaction{}
	resolved := []ExchangeTransactionWithRate{}
	err := postgresql.SelectAllFromDb(db, "desc", &details, "account_id = ?", e.AccountID)
	if err != nil {
		return resolved, postgresql.PaginationResponse{}, err
	}

	pagination, err := postgresql.SelectAllFromDbOrderByPaginated(db, "id", "desc", paginator, &details, "account_id = ?", e.AccountID)
	if err != nil {
		return resolved, pagination, err
	}

	for _, e := range details {
		r, _, err := e.ResolveExchangeTransaction(db)
		if err != nil {
			return resolved, pagination, err
		} else {
			resolved = append(resolved, r)
		}
	}
	return resolved, pagination, nil
}

func (e *ExchangeTransaction) GetResolvedExchangeTransactionByID(db *gorm.DB) (ExchangeTransactionWithRate, int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &e, "id = ?", e.ID)
	if nilErr != nil {
		return ExchangeTransactionWithRate{}, http.StatusBadRequest, nilErr
	}

	if err != nil {
		return ExchangeTransactionWithRate{}, http.StatusInternalServerError, err
	}

	resolved, code, err := e.ResolveExchangeTransaction(db)
	if err != nil {
		return resolved, code, err
	}
	return resolved, http.StatusOK, nil
}

func (e *ExchangeTransaction) ResolveExchangeTransaction(db *gorm.DB) (ExchangeTransactionWithRate, int, error) {
	var rate = Rate{ID: int64(e.RateID)}
	date := utility.FormatDateSpecialCase(e.CreatedAt)
	if e.RateID != 0 {
		code, err := rate.GetRateByID(db)
		if err != nil {
			return ExchangeTransactionWithRate{}, code, err
		}

	}
	return ExchangeTransactionWithRate{
		ID:              e.ID,
		AccountID:       e.AccountID,
		InitialAmount:   e.InitialAmount,
		FinalAmount:     e.FinalAmount,
		Status:          e.Status,
		TransactionName: fmt.Sprintf("%s to %s", rate.FromCurrency, rate.ToCurrency),
		Rate:            rate,
		Date:            date,
	}, http.StatusOK, nil
}
