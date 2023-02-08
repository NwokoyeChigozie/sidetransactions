package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
	"gorm.io/gorm"
)

type Transaction struct {
	ID               uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	TransactionID    string    `gorm:"column:transaction_id; type:varchar(255); not null; comment: 12 characters long string" json:"transaction_id"`
	PartiesID        string    `gorm:"column:parties_id; type:varchar(255); not null; comment: " json:"parties_id"`
	MilestoneID      string    `gorm:"column:milestone_id; type:varchar(255); comment: " json:"milestone_id"`
	BrokerID         string    `gorm:"column:broker_id; type:varchar(255); comment: " json:"broker_id"`
	Title            string    `gorm:"column:title; type:varchar(255); not null; comment: " json:"title"`
	Type             string    `gorm:"column:type; type:varchar(255); not null; comment: Transaction Type: product, service[oneoff], service[milestone]" json:"type"`
	Description      string    `gorm:"column:description; type:text; not null; comment: " json:"description"`
	Amount           float64   `gorm:"column:amount; type:decimal(20,2); not null; comment:" json:"amount"`
	Status           string    `gorm:"column:status; type:varchar(255); default: draft; comment: Transaction Status" json:"status"`
	Quantity         int       `gorm:"column:quantity; type:int" json:"quantity"`
	InspectionPeriod string    `gorm:"column:inspection_period; type:varchar(255); comment: " json:"inspection_period"`
	DueDate          string    `gorm:"column:due_date; type:varchar(255); comment: " json:"due_date"`
	ShippingFee      float64   `gorm:"column:shipping_fee; type:decimal(20,2); not null; comment:" json:"shipping_fee"`
	GracePeriod      string    `gorm:"column:grace_period; type:varchar(255); comment: Grace Period 48 hours" json:"grace_period"`
	Currency         string    `gorm:"column:currency; type:varchar(255); comment: Currency transaction made in" json:"currency"`
	DeletedAt        time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt        time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
	BusinessID       int       `gorm:"column:business_id; type:int" json:"business_id"`
	IsPaylinked      bool      `gorm:"column:is_paylinked; type:bool; default:false" json:"is_paylinked"`
	Country          string    `gorm:"column:country; type:varchar(255)" json:"country"`
	Source           string    `gorm:"column:source; type:varchar(255); default: api" json:"source"`
	TransUssdCode    int       `gorm:"column:trans_ussd_code; type:int" json:"trans_ussd_code"`
	Recipients       string    `gorm:"column:recipients; type:varchar(255)" json:"recipients"`
	DisputeHandler   string    `gorm:"column:dispute_handler; type:varchar(255)" json:"dispute_handler"`
	AmountPaid       float64   `gorm:"column:amount_paid; type:decimal(20,2); not null; comment:" json:"amount_paid"`
	EscrowCharge     float64   `gorm:"column:escrow_charge; type:decimal(20,2); not null; comment:" json:"escrow_charge"`
	EscrowWallet     string    `gorm:"column:escrow_wallet; type:varchar(255); default: no" json:"escrow_wallet"`
}

type CreateTransactionRequest struct {
	BusinessID       int         `json:"business_id"  pgvalidate:"exists=auth$business_profiles$account_id"`
	Parties          []Party     `json:"parties"  validate:"required"`
	Title            string      `json:"title"  validate:"required"`
	Type             string      `json:"type"  validate:"required,oneof=oneoff milestone"`
	EscrowWallet     string      `json:"escrow_wallet"  validate:"required,oneof=yes no"`
	Description      string      `json:"description"`
	Files            []File      `json:"files"`
	Milestones       []MileStone `json:"milestones"`
	Quantity         int         `json:"quantity"`
	Amount           float64     `json:"amount"`
	InspectionPeriod int         `json:"inspection_period"`
	GracePeriod      string      `json:"grace_period"`
	DueDate          string      `json:"due_date" validate:"required"`
	ShippingFee      float64     `json:"shipping_fee"`
	Currency         string      `json:"currency"  validate:"required"`
	Source           string      `json:"source" validate:"oneof=api instantescrow trizact transfer"`
	DisputeHandler   string      `json:"dispute_handler"`
	Paylinked        bool        `json:"paylinked"`
}
type EditTransactionRequest struct {
	TransactionID    string  `json:"transaction_id" validate:"required" pgvalidate:"exists=transaction$transactions$transaction_id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Quantity         int     `json:"quantity"`
	InspectionPeriod int     `json:"inspection_period"`
	DueDate          string  `json:"due_date"`
	ShippingFee      float64 `json:"shipping_fee"`
	Currency         string  `json:"currency"  validate:"required"`
	GracePeriod      string  `json:"grace_period"`
}

type Party struct {
	AccountID    int              `json:"account_id"`
	EmailAddress string           `json:"email_address"`
	PhoneNumber  string           `json:"phone_number"`
	Role         string           `json:"role"`
	Status       string           `json:"status"`
	AccessLevel  PartyAccessLevel `json:"access_level"`
}

type PartyAccessLevel struct {
	CanView    bool `json:"can_view"`
	CanReceive bool `json:"can_receive"`
	MarkAsDone bool `json:"mark_as_done"`
	Approve    bool `json:"approve"`
}

type PartyResponse struct {
	PartyID     int              `json:"party_id"`
	AccountID   int              `json:"account_id"`
	AccountName string           `json:"account_name"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Role        string           `json:"role"`
	Status      string           `json:"status"`
	AccessLevel PartyAccessLevel `json:"access_level"`
}
type MilestonesResponse struct {
	Index            int                           `json:"index"`
	MilestoneID      string                        `json:"milestone_id"`
	Title            string                        `json:"title"`
	Amount           float64                       `json:"amount"`
	Status           string                        `json:"status"`
	InspectionPeriod string                        `json:"inspection_period"`
	DueDate          string                        `json:"due_date"`
	Recipients       []MilestonesRecipientResponse `json:"recipients"`
}

type MilestonesRecipientResponse struct {
	AccountID   int     `json:"title"`
	AccountName string  `json:"amount"`
	Email       string  `json:"status"`
	PhoneNumber string  `json:"inspection_period"`
	Amount      float64 `json:"due_date"`
}

type File struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MileStone struct {
	Title            string               `json:"title"`
	Amount           float64              `json:"amount"`
	InspectionPeriod int                  `json:"inspection_period"`
	DueDate          string               `json:"due_date"`
	Status           string               `json:"status"`
	Description      string               `json:"description"`
	Quantity         int                  `json:"quantity"`
	ShippingFee      float64              `json:"shipping_fee"`
	GracePeriod      string               `json:"grace_period"`
	Recipients       []MileStoneRecipient `json:"recipients"`
}
type MileStoneRecipient struct {
	AccountID    int     `json:"account_id"`
	Amount       float64 `json:"amount"`
	EmailAddress string  `json:"email_address"`
	PhoneNumber  string  `json:"phone_number"`
}

type ResolveTransactionObj struct {
	TransactionID        string
	TransactionPartiesID string
	Title                string
	Type                 string
	Description          string
	Amount               float64
	Quantity             int
	ShippingFee          float64
	GracePeriod          string
	Currency             string
	Country              string
	BusinessID           int
	DisputeHandler       string
	EscrowWallet         string
}

func (t *Transaction) CreateTransaction(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &t)
	if err != nil {
		return fmt.Errorf("transaction creation failed: %v", err.Error())
	}
	return nil
}

func (t *Transaction) UpdateAllFields(db *gorm.DB) error {
	_, err := postgresql.SaveAllFields(db, &t)
	return err
}

func (t *Transaction) GetTransactionByTransactionID(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "transaction_id = ?", t.TransactionID)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
func (t *Transaction) GetTransactionByUssdCode(db *gorm.DB) (int, error) {
	err, nilErr := postgresql.SelectOneFromDb(db, &t, "trans_ussd_code = ?", t.TransUssdCode)
	if nilErr != nil {
		return http.StatusBadRequest, nilErr
	}

	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (t *Transaction) GetAllByTransactionID(db *gorm.DB) ([]Transaction, error) {
	details := []Transaction{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "transaction_id = ? ", t.TransactionID)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (t *Transaction) GetAllOthersByIDAndPartiesID(db *gorm.DB) ([]Transaction, error) {
	details := []Transaction{}
	err := postgresql.SelectAllFromDb(db, "asc", &details, "id != ? and parties_id = ?", t.ID, t.PartiesID)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (t *Transaction) GetAllOthersByAndQueries(db *gorm.DB, usePaylinked bool, CreatedAtInterval string, orderBy, order string) ([]Transaction, error) {
	var (
		details = []Transaction{}
		query   = ``
	)

	if t.TransactionID != "" {
		query = addQuery(query, fmt.Sprintf("transaction_id = '%v'", t.TransactionID), "AND")
	}
	if t.ID != 0 {
		query = addQuery(query, fmt.Sprintf("id = %v", t.ID), "AND")
	}
	if t.PartiesID != "" {
		query = addQuery(query, fmt.Sprintf("parties_id = '%v'", t.PartiesID), "AND")
	}
	if t.BusinessID != 0 {
		query = addQuery(query, fmt.Sprintf("business_id = %v", t.BusinessID), "AND")
	}
	if t.Status != "" {
		query = addQuery(query, fmt.Sprintf("status = '%v'", t.Status), "AND")
	}

	if usePaylinked {
		query = addQuery(query, fmt.Sprintf("is_paylinked = %v", t.IsPaylinked), "AND")
	}

	if CreatedAtInterval != "" {
		start, end := utility.GetStartAndEnd(CreatedAtInterval)
		query = addQuery(query, fmt.Sprintf("(created_at BETWEEN '%s' AND '%s')", start, end), "AND")
	}

	err := postgresql.SelectAllFromDbOrderBy(db, orderBy, order, &details, query)
	if err != nil {
		return details, err
	}
	return details, nil
}

func (t *Transaction) Delete(db *gorm.DB) error {
	err := postgresql.DeleteRecordFromDb(db, &t)
	if err != nil {
		return fmt.Errorf("transaction delete failed: %v", err.Error())
	}
	return nil
}
