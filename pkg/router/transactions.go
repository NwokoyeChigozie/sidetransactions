package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/pkg/controller/transactions"
	"github.com/vesicash/transactions-ms/pkg/middleware"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func Transaction(r *gin.Engine, ApiVersion string, validator *validator.Validate, db postgresql.Databases, logger *utility.Logger) *gin.Engine {
	extReq := request.ExternalRequest{Logger: logger, Test: false}
	transaction := transactions.Controller{Db: db, Validator: validator, Logger: logger, ExtReq: extReq}

	// transactionsUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion))
	// {

	// }

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/create", transaction.CreateTransaction)
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion), middleware.Authorize(db, extReq, middleware.ApiType))
	{
		transactionsApiUrl.GET("/listById/:id", transaction.ListTransactionsByID)
	}

	transactionsAppUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion), middleware.Authorize(db, extReq, middleware.AppType))
	{
		transactionsAppUrl.POST("/validate_on_db", transaction.ValidateOnDB)
	}
	return r
}
