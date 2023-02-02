package migrations

import "github.com/vesicash/transactions-ms/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.ActivityLog{},
		models.ProductTransaction{},
		models.TransactionState{},
		models.TransactionBroker{},
		models.TransactionDispute{},
		models.TransactionFile{},
		models.TransactionParty{},
		models.Transaction{},
	}
}
