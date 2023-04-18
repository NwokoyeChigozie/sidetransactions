package cronjobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func HandleTransactionAutoMark(extReq request.ExternalRequest, db postgresql.Databases) {
	query := fmt.Sprintf(`LOWER(status) <> '%v'`, strings.ToLower(transactions.GetTransactionStatus("d")))
	txb := models.Transaction{}
	transactionsSlice, err := txb.GetAllByQuery(db.Transaction, query)
	if err != nil {
		extReq.Logger.Error("error getting transactions: ", err.Error())
		return
	}

	for _, tx := range transactionsSlice {
		business, _ := transactions.GetBusinessProfileByAccountID(extReq, extReq.Logger, tx.BusinessID)
		if business.AutoTransactionStatusSettings {
			dueDate, err := utility.UnFormatDueDate(tx.DueDate)
			if err != nil {
				extReq.Logger.Error(fmt.Sprintf("error parsing due date %v for transaction %v", tx.DueDate, tx.TransactionID))
			} else {
				if dueDate.After(time.Now()) {
					user, _ := transactions.GetUserWithAccountID(extReq, tx.BusinessID)
					_, err := transactions.TransactionDeliveredService(extReq, extReq.Logger, db, models.TransactionDeliveredRequest{
						TransactionID: tx.TransactionID,
						MilestoneID:   tx.MilestoneID,
					}, user)
					if err != nil {
						extReq.Logger.Error("error updating transaction to delivered: ", err.Error())
					} else {
						extReq.Logger.Info(fmt.Sprintf("transaction %v automatically updated to delivered", tx.TransactionID))
					}
				}
			}
		}
	}
}
