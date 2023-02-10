package transactions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func EditTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.EditTransactionRequest, user external_models.User) (models.Transaction, int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return transaction, code, err
	}

	if req.Title != "" {
		titleSlice := strings.Split(transaction.Title, ";")
		if len(titleSlice) > 0 {
			titleSlice[0] = req.Title
		} else {
			titleSlice = []string{req.Title}
		}
		newTitle := strings.Join(titleSlice, ";")
		transaction.Title = newTitle
	}

	if req.Description != "" {
		transaction.Description = req.Description
	}

	if req.Quantity != 0 {
		transaction.Quantity = req.Quantity
	}

	if req.InspectionPeriod != 0 {
		transaction.InspectionPeriod = strconv.Itoa(req.InspectionPeriod)
	}

	if req.DueDate != "" {
		dueDate, err := validateDueDate(req.DueDate)
		if err != nil {
			return transaction, http.StatusBadRequest, fmt.Errorf("wrong due date format, try 2006-01-23")
		}
		transactionDueDate, _ := utility.GetUnixString(dueDate, "2006-01-02", "2006-01-02")
		transaction.DueDate = transactionDueDate
	}

	if req.ShippingFee != 0 {
		transaction.ShippingFee = req.ShippingFee
	}

	if req.Currency != "" {
		transaction.Currency = strings.ToUpper(req.Currency)
	}

	if req.GracePeriod != "" {
		gracePeriod, err := validateDueDate(req.GracePeriod)
		if err != nil {
			return transaction, http.StatusBadRequest, fmt.Errorf("wrong grace period format, try 2006-01-23")
		}
		transactionGracePeriod, _ := utility.GetUnixString(gracePeriod, "2006-01-02", "2006-01-02")
		transaction.GracePeriod = transactionGracePeriod
	}

	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return transaction, http.StatusInternalServerError, err
	}

	return transaction, http.StatusOK, nil
}

func DeleteTransactionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: transactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		if code == http.StatusBadRequest {
			return http.StatusOK, fmt.Errorf("Transaction has been deleted")
		}
		return code, err
	}

	err = transaction.Delete(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	status := GetTransactionStatus("deleted")

	_, err = CreateTransactionState(db, status, transactionID, transaction.MilestoneID, int(user.AccountID))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
