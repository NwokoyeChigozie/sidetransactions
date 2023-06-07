package transactions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func UpdateTransactionStatusService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.UpdateTransactionStatusRequest, user external_models.User) (int, error) {
	var (
		transaction               = models.Transaction{TransactionID: req.TransactionID, MilestoneID: req.MilestoneID}
		localStatus               = ""
		transactionMessage        = ""
		closedTransactionMessage  = ""
		transactionPartiesMessage = ""
		tType                     = ""
		message                   = ""
	)

	code, err := transaction.GetTransactionByTransactionIDAndMilestoneID(db.Transaction)
	if err != nil {
		return code, err
	}

	if req.Status == "cr" {
		localStatus = GetTransactionStatus("closed")
	}

	if req.Status == "sr" {
		mainTransaction := models.Transaction{TransactionID: req.TransactionID}
		code, err := mainTransaction.GetTransactionByTransactionID(db.Transaction)
		if err != nil {
			return code, err
		}
		amountPaid := mainTransaction.AmountPaid
		transactionMessage = "has failed to mark ongoing transaction as done"
		if amountPaid == 0 {
			localStatus = GetTransactionStatus("closed")
			closedTransactionMessage = "Transaction has been closed."
		} else {
			_, err := extReq.SendExternalRequest(request.RequestManualRefund, req.TransactionID)
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf("refund failed")
			}
			localStatus = GetTransactionStatus("cr")
			closedTransactionMessage = "Transaction has closed and payment refunded back to buyer."
		}

	}

	if !CheckTransactionStatus(req.Status) {
		return http.StatusBadRequest, fmt.Errorf("transaction status does not exist")
	}
	status := GetTransactionStatus(req.Status)

	if status == GetTransactionStatus("dr") {
		transactionTitle := strings.Split(transaction.Title, ";")[0]
		milestoneName := transactionTitle
		if transaction.Type != "milestone" {
			milestoneName = ""
		}
		transactionMessage = fmt.Sprintf("has rejected delivered %s transaction", milestoneName)
	} else if status == GetTransactionStatus("d") {
		if transaction.Type == "oneoff" {
			tType = "current"
		} else {
			tType = "ongoing milestone"
		}
		transactionMessage = fmt.Sprintf("has marked %s transaction as done.", tType)
	} else if status == GetTransactionStatus("closed") {
		closedTransactionMessage = "Transaction has been closed."
	} else if status == GetTransactionStatus("anf") {
		transactionPartiesMessage = "All parties have accepted transaction"
	}

	if localStatus != "" {
		transaction.Status = localStatus
	} else {
		transaction.Status = status
	}
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if req.Status == "da" {
		transaction.Status = GetTransactionStatus("cdp")
		err = transaction.UpdateAllFields(db.Transaction)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		transactionTitleSlice := strings.Split(transaction.Title, ";")
		transactionTitle := ""
		if len(transactionTitleSlice) > 1 {
			transactionTitle = transactionTitleSlice[1]
		}

		if transaction.Type == "oneoff" {
			tType = "current transaction"
		} else {
			tType = transactionTitle
		}
		transactionMessage = fmt.Sprintf("has approved delivered %s for payment.", tType)
		_, err = CreateTransactionState(db, "cdp", req.TransactionID, transaction.MilestoneID, int(user.AccountID))
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		_, err = CreateTransactionState(db, req.Status, req.TransactionID, transaction.MilestoneID, int(user.AccountID))
		if err != nil {
			return http.StatusInternalServerError, err

		}
	}

	if closedTransactionMessage != "" {
		message = closedTransactionMessage
	} else if transactionPartiesMessage != "" {
		message = transactionPartiesMessage
	} else {
		message = fmt.Sprintf("%v %v", user.EmailAddress, transactionMessage)
	}

	activityLog := models.ActivityLog{
		TransactionID: req.TransactionID,
		Description:   message,
	}
	err = activityLog.CreateActivityLog(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
