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

func HandleTransactionClose(extReq request.ExternalRequest, db postgresql.Databases) {
	nowTimeFormatted := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(`
	due_date < '%v' 
	and LOWER(status) = '%v' 
	and LOWER(status) = '%v' 
	and LOWER(status) = '%v' 
	and LOWER(status) = '%v'`,
		nowTimeFormatted,
		strings.ToLower(transactions.GetTransactionStatus("closed")),
		strings.ToLower(transactions.GetTransactionStatus("cr")),
		strings.ToLower(transactions.GetTransactionStatus("cnf")),
		strings.ToLower(transactions.GetTransactionStatus("cdc")),
	)
	txb := models.Transaction{}
	transactionsSlice, err := txb.GetAllByQueryWithLimit(db.Transaction, query, 20)
	if err != nil {
		extReq.Logger.Error("error getting transactions: ", err.Error())
		return
	}

	for _, tx := range transactionsSlice {
		pty := models.TransactionParty{TransactionID: tx.TransactionID}
		amountPaid := tx.AmountPaid
		transactionStatus := tx.Status
		transactionCurrency := tx.Currency
		parties, err := pty.GetAllByTransactionID(db.Transaction)
		if err != nil {
			extReq.Logger.Error(fmt.Errorf("error getting parties for transaction %v", tx.TransactionID))
		} else {
			continueProcess := true
			for _, party := range parties {
				if !strings.EqualFold(party.Status, "accepted") {
					// money has not been paid
					// close transaction by setting status to closed
					tx.Status = transactions.GetTransactionStatus("closed")
					tx.UpdateAllFields(db.Transaction)
					continueProcess = false
				}
			}

			if continueProcess {
				if amountPaid > 0 {
					// do refund
					refund(extReq, db, amountPaid, transactionCurrency, tx)
					tx.Status = transactions.GetTransactionStatus("cr")
					tx.UpdateAllFields(db.Transaction)
					continueProcess = false
				} else {
					tx.Status = transactions.GetTransactionStatus("cnf")
					tx.UpdateAllFields(db.Transaction)
				}
			}

			if continueProcess {
				if strings.EqualFold(transactionStatus, transactions.GetTransactionStatus("d")) {
					dueDate, err := utility.UnFormatDueDate(tx.DueDate)
					if err != nil {
						extReq.Logger.Error(fmt.Sprintf("error parsing due date %v for transaction %v", tx.DueDate, tx.TransactionID))
					} else {
						if dueDate.Before(time.Now()) {
							refund(extReq, db, amountPaid, transactionCurrency, tx)
							tx.Status = transactions.GetTransactionStatus("cr")
							tx.UpdateAllFields(db.Transaction)
						}
					}
					continueProcess = false
				}
			}

			if continueProcess {
				if statusInList(transactionStatus, []string{"dr", "ip", "af", "sr"}) {
					refund(extReq, db, amountPaid, transactionCurrency, tx)
					tx.Status = transactions.GetTransactionStatus("cr")
					tx.UpdateAllFields(db.Transaction)
					continueProcess = false
				}
			}

			if continueProcess {
				if statusInList(transactionStatus, []string{"anf", "draft"}) {
					tx.Status = transactions.GetTransactionStatus("closed")
					tx.UpdateAllFields(db.Transaction)
					continueProcess = false
				}
			}
		}

	}
}

func statusInList(txStatus string, statusCodes []string) bool {
	for _, statusCode := range statusCodes {
		if strings.EqualFold(txStatus, transactions.GetTransactionStatus(statusCode)) {
			return true
		}

	}
	return false
}

func refund(extReq request.ExternalRequest, db postgresql.Databases, amountPaid float64, transactionCurrency string, transaction models.Transaction) {
	buyer := models.TransactionParty{TransactionID: transaction.TransactionID, Role: "buyer"}
	_, err := buyer.GetTransactionPartyByTransactionIDAndRole(db.Transaction)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error getting buyer party for transaction %v", transaction.TransactionID))
		return
	}

	recipientCurrency := strings.ToUpper(transactionCurrency)
	senderCurrency := "ESCROW_" + strings.ToUpper(transactionCurrency)

	_, err = transactions.DebitWallet(extReq, db, amountPaid, senderCurrency, buyer.AccountID, "no", transaction.TransactionID)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error debiting buyer %v, walletcurrency:%v for transaction %v", buyer.AccountID, senderCurrency, transaction.TransactionID))
		return
	}

	_, err = transactions.CreditWallet(extReq, db, amountPaid, recipientCurrency, buyer.AccountID, true, "no", transaction.TransactionID)
	if err != nil {
		extReq.Logger.Error(fmt.Sprintf("error crediting buyer %v, walletcurrency:%v for transaction %v", buyer.AccountID, recipientCurrency, transaction.TransactionID))
		return
	}
}
