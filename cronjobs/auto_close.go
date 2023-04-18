package cronjobs

import (
	"fmt"
	"strings"

	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
)

func HandleTransactionAutoClose(extReq request.ExternalRequest, db postgresql.Databases) {
	query := fmt.Sprintf(`LOWER(status)='%v' or LOWER(status)='%v'`, strings.ToLower(transactions.GetTransactionStatus("cdc")), strings.ToLower("Closed - Disbursement Complete"))
	txb := models.Transaction{}
	transactionsSlice, err := txb.GetAllByQuery(db.Transaction, query)
	if err != nil {
		extReq.Logger.Error("error getting transactions: ", err.Error())
		return
	}

	for _, tx := range transactionsSlice {
		_, err = transactions.CreateTransactionState(db, transactions.GetTransactionStatus("closed"), tx.TransactionID, tx.MilestoneID, tx.BusinessID)
		if err != nil {
			extReq.Logger.Error("error creating transactions state: ", err.Error())
		} else {
			extReq.Logger.Info(fmt.Sprintf("transaction %v, automatically updated to DELIVERED", tx.TransactionID))
		}
	}

}
