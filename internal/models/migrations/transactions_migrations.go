package migrations

import "github.com/vesicash/transactions-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.ActivityLog{},
		models.ExchangeTransaction{},
		models.ProductTransaction{},
		models.Rate{},
		models.TransactionState{},
		models.TransactionBroker{},
		models.TransactionDispute{},
		models.TransactionDueDateExtensionRequest{},
		models.TransactionFile{},
		models.TransactionParty{},
		models.TransactionsRejected{},
		models.Transaction{},
	}
}
