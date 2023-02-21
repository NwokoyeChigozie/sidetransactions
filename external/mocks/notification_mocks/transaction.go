package notification_mocks

import (
	"fmt"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/mocks/auth_mocks"
	"github.com/vesicash/transactions-ms/utility"
)

func SendNewTransactionNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("new transaction notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("new transaction notification", nil, err.Error())
		return nil, err
	}

	logger.Info("new transaction notification", nil)

	return nil, nil
}

func SendTransactionAcceptedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction accepted notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("transaction accepted notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("transaction accepted notification", outBoundResponse)

	return nil, nil
}

func SendTransactionRejectedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction rejected notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("transaction rejected notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("transaction rejected notification", outBoundResponse)

	return nil, nil
}
func SendTransactionDeliveredRejectedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered rejected notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("transaction delivered rejected notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("transaction delivered rejected notification", outBoundResponse)

	return nil, nil
}
func SendDisputeOpenedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDAccountIDRequestModel)
	if !ok {
		logger.Info("dispute opened notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("dispute opened notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("dispute opened notification", outBoundResponse)

	return nil, nil
}
func SendTransactionDeliveredNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("transaction delivered notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("transaction delivered notification", outBoundResponse)

	return nil, nil
}

func SendDueDateProposalNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.DueDateExtensionProposalRequestModel)
	if !ok {
		logger.Info("due date proposal notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("due date proposal notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("due date proposal notification", outBoundResponse)

	return nil, nil
}
func SendDueDateExtendedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("due date extended notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("due date extended notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("due date extended notification", outBoundResponse)

	return nil, nil
}
func SendTransactionDeliveredAcceptedNotification(logger *utility.Logger, idata interface{}) (interface{}, error) {
	var (
		outBoundResponse map[string]interface{}
	)
	_, ok := idata.(external_models.TransactionIDRequestModel)
	if !ok {
		logger.Info("transaction delivered accepted notification", idata, "request data format error")
		return nil, fmt.Errorf("request data format error")
	}
	_, err := auth_mocks.GetAccessToken(logger)
	if err != nil {
		logger.Info("transaction delivered accepted notification", outBoundResponse, err.Error())
		return nil, err
	}

	logger.Info("transaction delivered accepted notification", outBoundResponse)

	return nil, nil
}
