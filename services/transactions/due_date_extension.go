package transactions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func RequestDueDateExtensionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.DueDateExtensionRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID}
	)

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	sellerParty := models.TransactionParty{TransactionID: req.TransactionID, Role: "seller"}
	code, err = sellerParty.GetTransactionPartyByTransactionIDAndRole(db.Transaction)
	if err != nil {
		return code, fmt.Errorf("seller not found: %v", err.Error())
	}

	if user.AccountID != uint(sellerParty.AccountID) {
		return http.StatusBadRequest, fmt.Errorf("you cannot make this request as you are not the seller")
	}

	dueDateExtension := models.TransactionDueDateExtensionRequest{
		AccountID:     int64(sellerParty.AccountID),
		TransactionID: req.TransactionID,
		Note:          req.Note,
	}
	err = dueDateExtension.CreateTransactionDueDateExtensionRequest(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendDueDateProposalNotification, external_models.DueDateExtensionProposalRequestModel{
		TransactionId: req.TransactionID,
		Note:          req.Note,
	})

	return http.StatusOK, nil
}

func ApproveDueDateExtensionService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ApproveDueDateExtensionRequest, user external_models.User) (int, error) {
	var (
		transaction = models.Transaction{TransactionID: req.TransactionID, MilestoneID: req.MilestoneID}
	)

	code, err := transaction.GetTransactionByTransactionIDAndMilestoneID(db.Transaction)
	if err != nil {
		return code, err
	}

	buyerParty := models.TransactionParty{TransactionID: req.TransactionID, Role: "buyer"}
	code, err = buyerParty.GetTransactionPartyByTransactionIDAndRole(db.Transaction)
	if err != nil {
		return code, fmt.Errorf("buyer not found: %v", err.Error())
	}

	if user.AccountID != uint(buyerParty.AccountID) {
		return http.StatusBadRequest, fmt.Errorf("you cannot make this request as you are not the buyer")
	}

	dueDate, err := validateDueDate(req.DueDate)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("incorrect due date format, try 2006-01-15")
	}
	newDueDate, _ := utility.GetUnixString(dueDate, "2006-01-02", "2006-01-02")

	transaction.DueDate = newDueDate
	transaction.InspectionPeriod = strconv.Itoa(req.InspectionPeriod)
	err = transaction.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extReq.SendExternalRequest(request.SendDueDateExtendedNotification, external_models.TransactionIDRequestModel{
		TransactionId: req.TransactionID,
	})

	return http.StatusOK, nil
}
