package models

import "github.com/vesicash/transactions-ms/external/external_models"

var (
	MyIdentity *external_models.User
	Token      string
)

type GetEscrowChargeRequest struct {
	BusinessID int     `json:"business_id" validate:"required" pgvalidate:"exists=auth$business_profiles$account_id"`
	Amount     float64 `json:"amount" validate:"required"`
}
type GetEscrowChargeResponse struct {
	Amount float64 `json:"amount"`
	Charge float64 `json:"charge"`
}
