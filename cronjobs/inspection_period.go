package cronjobs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
)

func HandleTransactionInspectionPeriod(extReq request.ExternalRequest, db postgresql.Databases) {
	query := fmt.Sprintf(`LOWER(status) = '%v'`, strings.ToLower(transactions.GetTransactionStatus("d")))
	txb := models.Transaction{}
	transactionsSlice, err := txb.GetAllByQueryWithLimit(db.Transaction, query, 100)
	if err != nil {
		extReq.Logger.Error("error getting transactions: ", err.Error())
		return
	}

	for _, tx := range transactionsSlice {
		inspectionPeriodUnix, err := strconv.Atoi(tx.InspectionPeriod)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error parsing inspectionperiod %v for transaction %v", tx.InspectionPeriod, tx.TransactionID))
		} else {
			inspectionPeriod := time.Unix(int64(inspectionPeriodUnix), 0)
			if inspectionPeriod.After(time.Now()) {
				payment, err := transactions.ListPayment(extReq, tx.TransactionID)
				if err != nil {
					extReq.Logger.Error(fmt.Sprintf("error getting payment record for transaction %v", tx.TransactionID))
				} else {
					user, _ := transactions.GetUserWithAccountID(extReq, tx.BusinessID)
					if payment.IsPaid {
						_, err := transactions.SatisfiedService(extReq, extReq.Logger, db, tx.TransactionID, user)
						if err != nil {
							extReq.Logger.Error(fmt.Sprintf("error making transaction %v as satisfied", tx.TransactionID))
						} else {
							extReq.Logger.Info(fmt.Sprintf("Transaction %v marked as satisfied", tx.TransactionID))
						}
					}

				}
			}
		}

	}
}
