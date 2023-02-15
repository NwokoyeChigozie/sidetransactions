package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
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

type UpdateTransactionPartiesRequest struct {
	TransactionID string           `json:"transaction_id" validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	Parties       map[string]Party `json:"parties"  validate:"required"`
}
type UpdateTransactionPartyStatusRequest struct {
	TransactionID string `json:"transaction_id" validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	AccountID     int    `json:"account_id" validate:"required"  pgvalidate:"exists=auth$users$account_id"`
	Status        string `json:"status" validate:"required"`
}
type AssignTransactionBuyerRequest struct {
	TransactionID string `json:"transaction_id" pgvalidate:"exists=transaction$transactions$transaction_id"`
	UssdCode      int    `json:"ussd_code" pgvalidate:"exists=transaction$transactions$trans_ussd_code"`
	PhoneNumber   string `json:"phone_number" validate:"required"`
}

func (t *TransactionParty) CreateTransactionParty(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction party creation failed: %v", err.Error())
	}
	return nil
}
func (t *TransactionParty) CreateTransactionsParties(db *gorm.DB, parties []TransactionParty) ([]TransactionParty, error) {
	err := postgresql.CreateOneRecord(db, &parties)
	if err != nil {
		return parties, fmt.Errorf("transaction creation failed: %v", err.Error())
	}
	return parties, nil
}

func (t *TransactionParty) GetTransactionPartyByTransactionPartiesIDAndRole(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_parties_id = ? and role = ?", t.TransactionPartiesID, t.Role)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (t *TransactionParty) GetTransactionPartyByTransactionIDAndAccountID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ? and account_id = ?", t.TransactionID, t.AccountID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
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

func (t *TransactionParty) GetAllByAndQueriesForUniqueValue(db *gorm.DB, CreatedAtInterval string, orderBy, order string, groupColumn string, paginator postgresql.Pagination) ([]TransactionParty, postgresql.PaginationResponse, error) {
	var (
		details = []TransactionParty{}
		query   = ``
	)

	if t.TransactionID != "" {
		query = addQuery(query, fmt.Sprintf("transaction_id = '%v'", t.TransactionID), "AND")
	}
	if t.TransactionPartiesID != "" {
		query = addQuery(query, fmt.Sprintf("transaction_parties_id = '%v'", t.TransactionPartiesID), "AND")
	}
	if t.ID != 0 {
		query = addQuery(query, fmt.Sprintf("id = %v", t.ID), "AND")
	}
	if t.Role != "" {
		query = addQuery(query, fmt.Sprintf("role = '%v'", t.Role), "AND")
	}
	if t.Status != "" {
		query = addQuery(query, fmt.Sprintf("status = '%v'", t.Status), "AND")
	}

	if CreatedAtInterval != "" {
		start, end := utility.GetStartAndEnd(CreatedAtInterval)
		query = addQuery(query, fmt.Sprintf("(created_at BETWEEN '%s' AND '%s')", start.Format(time.RFC3339), end.Format(time.RFC3339)), "AND")
	}

	totalPages, err := postgresql.SelectAllFromByGroup(db, orderBy, order, &paginator, &details, query, groupColumn)
	if err != nil {
		return details, totalPages, err
	}
	return details, totalPages, nil
}

func (t *TransactionParty) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &t)
	return err
}
