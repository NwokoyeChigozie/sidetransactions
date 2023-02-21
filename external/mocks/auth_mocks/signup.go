package auth_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/utility"
)

func SignupUser(logger *utility.Logger, idata interface{}) (external_models.User, error) {

	var (
		outBoundResponse external_models.SignupModel
	)

	data, ok := idata.(external_models.CreateUserRequestModel)
	if !ok {
		logger.Info("signup", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	logger.Info("signup", data)

	if User == nil {
		logger.Info("signup", User, "user not provided")
		return external_models.User{}, fmt.Errorf("user not provided")
	}

	logger.Info("signup", User, "user found")
	return *User, nil
}
