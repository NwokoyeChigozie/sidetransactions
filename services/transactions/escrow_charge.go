package transactions

import (
	"math"
	"net/http"

	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func GetEscrowChargeService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.GetEscrowChargeRequest) (models.GetEscrowChargeResponse, int, error) {
	businessProfile, err := GetBusinessProfileByAccountID(extReq, logger, req.BusinessID)
	if err != nil {
		return models.GetEscrowChargeResponse{}, http.StatusInternalServerError, err
	}
	currency := businessProfile.Currency

	businessCharge, err := getBusinessChargeWithBusinessIDAndCurrency(extReq, req.BusinessID, currency)
	if err != nil {
		businessCharge, err = initBusinessCharge(extReq, req.BusinessID, currency)
		if err != nil {
			return models.GetEscrowChargeResponse{}, http.StatusInternalServerError, err
		}
	}

	charge := getEscrowCharge(businessCharge, req.Amount)
	return models.GetEscrowChargeResponse{
		Amount: req.Amount,
		Charge: math.Round(charge*10000) / 10000,
	}, http.StatusOK, nil
}
