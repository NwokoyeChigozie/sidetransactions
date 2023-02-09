package models

import (
	"time"

	"github.com/vesicash/transactions-ms/external/external_models"
)

type TransactionByIDResponse struct {
	ID                  uint                        `json:"id"`
	TransactionID       string                      `json:"transaction_id"`
	PartiesID           string                      `json:"parties_id"`
	MilestoneID         string                      `json:"milestone_id"`
	BrokerID            string                      `json:"broker_id"`
	Title               string                      `json:"title"`
	Type                string                      `json:"type"`
	Description         string                      `json:"description"`
	Amount              float64                     `json:"amount"`
	Status              string                      `json:"status"`
	Quantity            int                         `json:"quantity"`
	InspectionPeriod    string                      `json:"inspection_period"`
	DueDate             string                      `json:"due_date"`
	ShippingFee         float64                     `json:"shipping_fee"`
	GracePeriod         string                      `json:"grace_period"`
	Currency            string                      `json:"currency"`
	DeletedAt           time.Time                   `json:"deleted_at"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
	BusinessID          int                         `json:"business_id"`
	IsPaylinked         bool                        `json:"is_paylinked"`
	Source              string                      `json:"source"`
	TransUssdCode       int                         `json:"trans_ussd_code"`
	Recipients          []MileStoneRecipient        `json:"recipients"`
	DisputeHandler      string                      `json:"dispute_handler"`
	AmountPaid          float64                     `json:"amount_paid"`
	EscrowCharge        float64                     `json:"escrow_charge"`
	EscrowWallet        string                      `json:"escrow_wallet"`
	Products            []ProductTransaction        `json:"products"`
	Parties             map[string]TransactionParty `json:"parties"`
	Members             []PartyResponse             `json:"members"`
	Files               []TransactionFile           `json:"files"`
	TotalAmount         float64                     `json:"total_amount"`
	Milestones          []MilestonesResponse        `json:"milestones"`
	Broker              TransactionBroker           `json:"broker"`
	Activities          []ActivityLog               `json:"activities"`
	Country             external_models.Country     `json:"country"`
	DueDateFormatted    string                      `json:"due_date_formatted"`
	TransactionClosedAt time.Time                   `json:"transaction_closed_at"`
	IsDisputed          bool                        `json:"is_disputed"`
}

type TransactionCreateResponse struct {
	ID               uint                 `json:"id"`
	TransactionID    string               `json:"transaction_id"`
	PartiesID        string               `json:"parties_id"`
	MilestoneID      string               `json:"milestone_id"`
	BrokerID         string               `json:"broker_id"`
	Title            string               `json:"title"`
	Type             string               `json:"type"`
	Description      string               `json:"description"`
	Amount           float64              `json:"amount"`
	Status           string               `json:"status"`
	Quantity         int                  `json:"quantity"`
	InspectionPeriod string               `json:"inspection_period"`
	DueDate          string               `json:"due_date"`
	ShippingFee      float64              `json:"shipping_fee"`
	Currency         string               `json:"currency"`
	DeletedAt        time.Time            `json:"deleted_at"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
	IsPaylinked      bool                 `json:"is_paylinked"`
	Source           string               `json:"source"`
	TransUssdCode    int                  `json:"trans_ussd_code"`
	Recipients       []MileStoneRecipient `json:"recipients"`
	DisputeHandler   string               `json:"dispute_handler"`
	AmountPaid       float64              `json:"amount_paid"`
	EscrowCharge     float64              `json:"escrow_charge"`
	EscrowWallet     string               `json:"escrow_wallet"`
	Country          string               `json:"country"`
	Products         []ProductTransaction `json:"products"`
	Parties          []PartyResponse      `json:"parties"`
	Files            []TransactionFile    `json:"files"`
	Milestones       []MilestonesResponse `json:"milestones"`
	Broker           TransactionBroker    `json:"broker"`
}

type ListTransactionsRequest struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Filter     string `json:"filter" validate:"oneof=day week month"`
}
type ListTransactionByBusinessRequest struct {
	BusinessID int    `json:"business_id" validate:"required" pgvalidate:"exists=auth$business_profiles$account_id"`
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Paylinked  bool   `json:"paylinked"`
	Filter     string `json:"filter" validate:"oneof=day week month"`
}
type ListByBusinessFromMondayToThursdayRequest struct {
	BusinessID int    `json:"business_id" validate:"required" pgvalidate:"exists=auth$business_profiles$account_id"`
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	Paylinked  bool   `json:"paylinked"`
}

type ListTransactionByUserRequest struct {
	Role       string `json:"role"`
	Paylinked  bool   `json:"paylinked"`
	StatusCode string `json:"status_code"`
}
