package transactions

import (
	"fmt"
	"net/http"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func TransactionDeliveredService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.TransactionDeliveredRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID, MilestoneID: req.MilestoneID}
	)

	code, err := transaction.GetTransactionByTransactionIDAndMilestoneID(db.Transaction)
	if err != nil {
		return code, err
	}

	transaction.Status = GetTransactionStatus("d")
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, "d", req.TransactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendTransactionDeliveredNotification, external_models.TransactionIDRequestModel{
		TransactionId: req.TransactionID,
	})

	return http.StatusOK, nil
}

func SatisfiedService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: transactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	buyerParty := models.TransactionParty{TransactionID: transactionID, Role: "buyer"}
	code, err = buyerParty.GetTransactionPartyByTransactionIDAndRole(db.Transaction)
	if err != nil {
		return code, fmt.Errorf("buyer not found: %v", err.Error())
	}

	if user.AccountID != uint(buyerParty.AccountID) {
		return http.StatusBadRequest, fmt.Errorf("you cannot make this request as you are not the buyer")
	}

	transaction.Status = GetTransactionStatus("da")
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, "da", transactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendTransactionDeliveredAcceptedNotification, external_models.TransactionIDRequestModel{
		TransactionId: transactionID,
	})

	transaction.Status = GetTransactionStatus("cdp")
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, "cdp", transactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
