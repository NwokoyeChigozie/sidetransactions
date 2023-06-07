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

func AcceptTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: transactionID}
		statusCode  = ""
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	if transaction.Status == "Accepted - Not Funded" || transaction.Status == "Accepted - Funded" {
		return http.StatusOK, nil
	}

	payment, err := ListPayment(extReq, transactionID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !payment.IsPaid {
		statusCode = "anf"
	} else {
		statusCode = "af"
	}

	transaction.Status = GetTransactionStatus(statusCode)
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, statusCode, transactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendTransactionAcceptedNotification, external_models.TransactionIDRequestModel{
		TransactionId: transactionID,
	})
	return http.StatusOK, nil
}
func RejectTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.RejectTransactionRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID}
		statusCode  = ""
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	payment, err := ListPayment(extReq, req.TransactionID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !payment.IsPaid {
		statusCode = "sr"
	} else {
		statusCode = "fr"
	}

	transaction.Status = GetTransactionStatus(statusCode)
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, statusCode, req.TransactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.Reason != "" && transaction.BusinessID != 0 {
		transactionRejected := models.TransactionsRejected{
			AccountID:     int64(transaction.BusinessID),
			TransactionID: req.TransactionID,
			Reason:        req.Reason,
		}

		err = transactionRejected.CreateTransactionsRejected(db.Transaction)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// Send Refund If The Transaction Has Been PAid For
	if payment.IsPaid {
		// Send Refund
		_, err := extReq.SendExternalRequest(request.RequestManualRefund, req.TransactionID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("refund failed")
		}
		transaction.Status = GetTransactionStatus("cr")
	} else {
		transaction.Status = GetTransactionStatus("closed")
	}

	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, transaction.Status, req.TransactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendTransactionRejectedNotification, external_models.TransactionIDRequestModel{
		TransactionId: req.TransactionID,
	})

	return http.StatusOK, nil
}

func RejectTransactionDeliveryService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.RejectTransactionRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID}
		statusCode  = "dr"
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	transaction.Status = GetTransactionStatus(statusCode)
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, statusCode, req.TransactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	statusCode = "closed"

	transaction.Status = GetTransactionStatus(statusCode)
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = CreateTransactionState(db, statusCode, req.TransactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendTransactionDeliveredRejectedNotification, external_models.TransactionIDRequestModel{
		TransactionId: req.TransactionID,
	})

	return http.StatusOK, nil
}
