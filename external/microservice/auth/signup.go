package auth

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
)

func (r *RequestObj) SignupUser() (external_models.User, error) {

	var (
		outBoundResponse external_models.SignupModel
		logger           = r.Logger
		idata            = r.RequestData
	)

	data, ok := idata.(external_models.CreateUserRequestModel)
	if !ok {
		logger.Info("signup", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	logger.Info("signup", data)
	err := r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("signup", outBoundResponse, err.Error())
		return outBoundResponse.Data, err
	}
	logger.Info("signup", outBoundResponse, outBoundResponse.Data)

	return outBoundResponse.Data, nil
}
