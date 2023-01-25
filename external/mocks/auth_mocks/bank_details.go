package auth_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/utility"
)

var (
	BankDetail *external_models.BankDetail
)

func GetBankDetails(logger *utility.Logger, idata interface{}) (external_models.BankDetail, error) {

	var (
		outBoundResponse external_models.GetBankDetailResponse
	)

	_, ok := idata.(external_models.GetBankDetailModel)
	if !ok {
		logger.Info("get bank detail", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	if BankDetail == nil {
		logger.Info("get BankDetail", BankDetail, "BankDetail not provided")
		return external_models.BankDetail{}, fmt.Errorf("BankDetail not provided")
	}

	logger.Info("get BankDetail", BankDetail, "BankDetail found")
	return *BankDetail, nil
}
