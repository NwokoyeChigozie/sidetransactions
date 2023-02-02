package transactions

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}
