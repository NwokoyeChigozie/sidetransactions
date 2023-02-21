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

func CreateDisputeService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.CreateDisputeRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}
	transactionDispute := models.TransactionDispute{
		DisputeID:     utility.RandomString(16),
		TransactionID: transaction.TransactionID,
		Reason:        req.Reason,
		DisputeStatus: req.DisputeStatus,
		Decision:      req.Decision,
	}
	err = transactionDispute.CreateTransactionDispute(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	transaction.Status = GetTransactionStatus("cd")
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendDisputeOpenedNotification, external_models.TransactionIDAccountIDRequestModel{
		TransactionId: req.TransactionID,
		AccountId:     user.AccountID,
	})
	return http.StatusOK, nil
}

func UpdateDisputeService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.CreateDisputeRequest, user external_models.User) (int, error) {
	var (
		transactionDispute = models.TransactionDispute{TransactionID: req.TransactionID}
	)

	code, err := transactionDispute.GetTransactionDisputeByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	if req.Reason != "" {
		transactionDispute.Reason = req.Reason
	}
	if req.DisputeStatus != "" {
		transactionDispute.DisputeStatus = req.DisputeStatus
	}
	if req.Decision != "" {
		transactionDispute.Decision = req.Decision
	}
	err = transactionDispute.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func GetDisputeByTransactionIDService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, user external_models.User) (*models.TransactionDispute, int, error) {
	var (
		transaction = models.Transaction{TransactionID: transactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return nil, code, err
	}

	transactionDispute := models.TransactionDispute{TransactionID: transactionID}
	code, err = transactionDispute.GetTransactionDisputeByTransactionID(db.Transaction)
	if err != nil {
		if code == http.StatusInternalServerError {
			return nil, code, err
		}
		return nil, http.StatusOK, nil
	}

	return &transactionDispute, http.StatusOK, nil
}

func GetDisputeByUserService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, user external_models.User, paginator postgresql.Pagination) ([]models.TransactionDispute, postgresql.PaginationResponse, int, error) {
	var (
		disputes         = []models.TransactionDispute{}
		transactionParty = models.TransactionParty{AccountID: int(user.AccountID)}
	)

	transactionParties, pagination, err := transactionParty.GetAllByAndQueriesForUniqueValueForDispute(db.Transaction, "", "id", "desc", "transaction_id", paginator)
	if err != nil {
		return []models.TransactionDispute{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	for _, t := range transactionParties {
		transactionDispute := models.TransactionDispute{TransactionID: t.TransactionID}
		code, err := transactionDispute.GetTransactionDisputeByTransactionID(db.Transaction)
		if err != nil {
			if code == http.StatusInternalServerError {
				return []models.TransactionDispute{}, pagination, code, err
			}
			fmt.Println(err)
		} else {
			disputes = append(disputes, transactionDispute)
		}
	}

	return disputes, pagination, http.StatusOK, nil
}
