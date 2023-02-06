package auth_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/utility"
)

var (
	BusinessCharge *external_models.BusinessCharge
)

func GetBusinessCharge(logger *utility.Logger, idata interface{}) (external_models.BusinessCharge, error) {

	_, ok := idata.(external_models.GetBusinessChargeModel)
	if !ok {
		logger.Info("get business charge", idata, "request data format error")
		return external_models.BusinessCharge{}, fmt.Errorf("request data format error")
	}

	if BusinessCharge == nil {
		logger.Info("get business charge", BusinessCharge, "businessCharge not provided")
		return external_models.BusinessCharge{}, fmt.Errorf("businessCharge not provided")
	}

	return *BusinessCharge, nil
}

func InitBusinessCharge(logger *utility.Logger, idata interface{}) (external_models.BusinessCharge, error) {

	_, ok := idata.(external_models.InitBusinessChargeModel)
	if !ok {
		logger.Info("init business charge", idata, "request data format error")
		return external_models.BusinessCharge{}, fmt.Errorf("request data format error")
	}

	if BusinessCharge == nil {
		logger.Info("init business charge", UserProfile, "businessCharge not provided")
		return external_models.BusinessCharge{}, fmt.Errorf("businessCharge not provided")
	}

	return *BusinessCharge, nil
}
