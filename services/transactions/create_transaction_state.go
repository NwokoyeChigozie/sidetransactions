package transactions

import (
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
)

func CreateTransactionState(db postgresql.Databases, status, transactionID, mileStoneID string, AccountID int) (models.TransactionState, error) {
	var (
		transactionState = models.TransactionState{
			AccountID:     int64(AccountID),
			TransactionID: transactionID,
			MilestoneID:   mileStoneID,
			Status:        GetTransactionStatus(status),
		}
	)

	err := transactionState.CreateTransactionState(db.Transaction)
	if err != nil {
		return transactionState, err
	}

	return transactionState, nil
}
