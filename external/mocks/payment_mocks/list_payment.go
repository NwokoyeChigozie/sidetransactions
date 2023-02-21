package payment_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/mocks/auth_mocks"
	"github.com/vesicash/transactions-ms/utility"
)

var (
	ListPaymentObj *external_models.ListPayment
)

func ListPayment(logger *utility.Logger, idata interface{}) (external_models.ListPayment, error) {
	var (
		outBoundResponse external_models.ListPaymentsResponse
	)
	_, ok := idata.(string)
	if !ok {
		logger.Info("list payment by transaction id", idata, "request data format error")
		return outBoundResponse.Data.Payment, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("list payment by transaction id", outBoundResponse, err.Error())
		return outBoundResponse.Data.Payment, err
	}

	if ListPaymentObj == nil {
		logger.Info("list payment by tansaction id", ListPaymentObj, "ListPayment not provided")
		return outBoundResponse.Data.Payment, fmt.Errorf("ListPayment not provided")
	}

	return *ListPaymentObj, nil
}
