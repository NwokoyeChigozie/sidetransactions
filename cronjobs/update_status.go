package cronjobs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
)

func HandleUpdateStatus(extReq request.ExternalRequest, db postgresql.Databases) {
	query := fmt.Sprintf(`LOWER(status) = '%v'`, strings.ToLower(transactions.GetTransactionStatus("da")))
	txb := models.Transaction{}
	transactionsSlice, err := txb.GetAllByQueryWithLimit(db.Transaction, query, 20)
	if err != nil {
		extReq.Logger.Error("error getting transactions: ", err.Error())
		return
	}

	for _, tx := range transactionsSlice {
		extReq.Logger.Info(fmt.Sprintf("processing update status job for transaction with id: %v", tx.ID))
		createActivityLog(extReq, db, tx, "cdp")
		_, err := transactions.ListPayment(extReq, tx.TransactionID)
		if err != nil {
			extReq.Logger.Error("error getting payment record for transaction %v", tx.TransactionID)
		} else {
			sendTransactionConfirmed(extReq, db, tx)
			createActivityLog(extReq, db, tx, "cdc")
		}
	}
}

func createActivityLog(extReq request.ExternalRequest, db postgresql.Databases, tx models.Transaction, statusCode string) {
	transactionType := tx.Type
	tx.Status = transactions.GetTransactionStatus(statusCode)
	tx.UpdateAllFields(db.Transaction)

	var transactionTitle string
	transactionTitleSlice := strings.Split(tx.Title, ";")
	if len(transactionTitleSlice) > 1 {
		transactionTitle = transactionTitleSlice[1]
	}
	cType := transactionTitle
	if strings.EqualFold(transactionType, "oneoff") {
		cType = "current transaction"
	}

	var description string
	if strings.EqualFold(statusCode, "cdp") {
		description = fmt.Sprintf("Payment for %v is currently being processed", cType)
	} else if strings.EqualFold(statusCode, "cdc") {
		description = fmt.Sprintf("Payment for %v is disbursed successfully", cType)
	}

	activityLog := models.ActivityLog{
		TransactionID: tx.TransactionID,
		Description:   description,
	}
	activityLog.CreateActivityLog(db.Transaction)
}

func sendTransactionConfirmed(extReq request.ExternalRequest, db postgresql.Databases, transaction models.Transaction) {
	buyer := models.TransactionParty{TransactionID: transaction.TransactionID, Role: "buyer"}
	_, err := buyer.GetTransactionPartyByTransactionIDAndRole(db.Transaction)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting buyer party for transaction %v", transaction.TransactionID))
		return
	}

	var milestoneRecipients []models.MileStoneRecipient
	transactionCurrency := strings.ToUpper(transaction.Currency)
	err = json.Unmarshal([]byte(transaction.Recipients), &milestoneRecipients)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error unmarshaling recipients for transaction %v", transaction.TransactionID))
		return
	}

	for _, recipient := range milestoneRecipients {
		_, err := extReq.SendExternalRequest(request.WalletTransfer, external_models.WalletTransferRequest{
			SenderAccountID:    buyer.AccountID,
			RecipientAccountID: recipient.AccountID,
			FinalAmount:        recipient.Amount,
			SenderCurrency:     "ESCROW_" + transactionCurrency,
			RecipientCurrency:  transactionCurrency,
		})
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("wallet transfer error for recipient %v, transaction %v, error: %v", recipient.AccountID, transaction.TransactionID, err.Error()))
		}
	}

}
