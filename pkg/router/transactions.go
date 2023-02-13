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
		transactionsAuthUrl.PATCH("/edit", transaction.EditTransaction)
		transactionsAuthUrl.DELETE("/delete/:id", transaction.DeleteTransaction)
		transactionsAuthUrl.POST("/listByUser", transaction.ListTransactionsByUser)
		transactionsAuthUrl.GET("list/archived", transaction.ListArchivedTransactions)
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion), middleware.Authorize(db, extReq, middleware.ApiType))
	{
		transactionsApiUrl.POST("/list", transaction.ListTransactions)
		transactionsApiUrl.GET("/listById/:id", transaction.ListTransactionsByID)
		transactionsApiUrl.GET("/list-transactions-by-ussd-code/:code", transaction.ListTransactionsByUSSDCode)
		transactionsApiUrl.POST("/listByBusiness", transaction.ListTransactionsByBusiness)
		transactionsApiUrl.POST("/listByBusinessFromMondayToThursday", transaction.ListByBusinessFromMondayToThursday)
		transactionsApiUrl.PATCH("/parties/update", transaction.UpdateTransactionParties)
		transactionsApiUrl.PATCH("/parties/update-status", transaction.UpdateTransactionPartyStatus)
		transactionsApiUrl.POST("/assign/buyer", transaction.AssignTransactionBuyer)
		transactionsApiUrl.POST("/broker/update", transaction.UpdateTransactionBroker)
	}

	transactionsAppUrl := r.Group(fmt.Sprintf("%v/transactions", ApiVersion), middleware.Authorize(db, extReq, middleware.AppType))
	{
		transactionsAppUrl.POST("/validate_on_db", transaction.ValidateOnDB)
	}
	return r
}
