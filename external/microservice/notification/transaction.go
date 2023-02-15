package notification

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
)

func (r *RequestObj) SendNewTransactionNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("new transaction notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("new transaction notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("new transaction notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("new transaction notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("new transaction notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) SendTransactionAcceptedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction accepted notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("transaction accepted notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction accepted notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("transaction accepted notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction accepted notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) SendTransactionRejectedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction rejected notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("transaction rejected notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction rejected notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("transaction rejected notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction rejected notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SendTransactionDeliveredRejectedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered rejected notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("transaction delivered rejected notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction delivered rejected notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("transaction delivered rejected notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction delivered rejected notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SendDisputeOpenedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDAccountIDRequestModel)
	if !ok {
		logger.Info("dispute opened notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("dispute opened notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("dispute opened notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("dispute opened notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("dispute opened notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SendTransactionDeliveredNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("transaction delivered notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction delivered notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("transaction delivered notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction delivered notification", outBoundResponse)

	return nil, nil
}

func (r *RequestObj) SendDueDateProposalNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.DueDateExtensionProposalRequestModel)
	if !ok {
		logger.Info("due date proposal notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("due date proposal notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("due date proposal notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("due date proposal notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("due date proposal notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SendDueDateExtendedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("due date extended notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("due date extended notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("due date extended notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("due date extended notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("due date extended notification", outBoundResponse)

	return nil, nil
}
func (r *RequestObj) SendTransactionDeliveredAcceptedNotification() (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
		logger           = r.Logger
		idata            = r.RequestData
	)
	data, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered accepted notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	accessToken, err := r.getAccessTokenObject().GetAccessToken()
	if err != nil {
		logger.Info("transaction delivered accepted notification", outBoundResponse, err.Error())
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"v-private-key": accessToken.PrivateKey,
		"v-public-key":  accessToken.PublicKey,
	}

	logger.Info("transaction delivered accepted notification", data)
	err = r.getNewSendRequestObject(data, headers, "").SendRequest(&outBoundResponse)
	if err != nil {
		logger.Info("transaction delivered accepted notification", outBoundResponse, err.Error())
		return nil, err
	}
	logger.Info("transaction delivered accepted notification", outBoundResponse)

	return nil, nil
}
