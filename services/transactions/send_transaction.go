package transactions

import (
	"net/http"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func SendTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: transactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	transaction.Status = GetTransactionStatus("sac")
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, "sac", transactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendNewTransactionNotification, external_models.TransactionIDRequestModel{
		TransactionId: transactionID,
	})

	return http.StatusOK, nil
}
