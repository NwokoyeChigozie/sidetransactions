package transactions

import (
	"net/http"

	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
)

func ValidateOnDbService(req models.ValidateOnDBReq, db postgresql.Databases) (bool, int, error) {
	if req.Type == "notexists" {
		if req.Value != nil {
			return !postgresql.CheckExistsInTable(db.Transaction, req.Table, req.Query, req.Value), http.StatusOK, nil
		} else {
			return !postgresql.CheckExistsInTable(db.Transaction, req.Table, req.Query), http.StatusOK, nil
		}

	} else if req.Type == "exists" {
		if req.Value != nil {
			return postgresql.CheckExistsInTable(db.Transaction, req.Table, req.Query, req.Value), http.StatusOK, nil
		} else {
			return postgresql.CheckExistsInTable(db.Transaction, req.Table, req.Query), http.StatusOK, nil
		}

	} else {
		return false, http.StatusOK, nil
	}
}
