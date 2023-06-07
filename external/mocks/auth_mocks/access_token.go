package auth_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/utility"
)

var (
	AccessToken external_models.AccessToken
)

func GetAccessToken(logger *utility.Logger) (external_models.AccessToken, error) {
	logger.Info("get access tokens", "get access tokens called")
	return AccessToken, nil
}

func GetAccessTokenByKey(logger *utility.Logger, idata interface{}) (external_models.AccessToken, error) {
	var (
		outBoundResponse external_models.GetAccessTokenModel
	)
	_, ok := idata.(string)
	if !ok {
		logger.Error("get access token by key", idata, "request data format error")
		return outBoundResponse.Data, fmt.Errorf("request data format error")
	}

	return AccessToken, nil
}
