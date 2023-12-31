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

	// transactionsUrl := r.Group(fmt.Sprintf("%v", ApiVersion))
	// {

	// }

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/create", transaction.CreateTransaction)
		transactionsAuthUrl.PATCH("/edit", transaction.EditTransaction)
		transactionsAuthUrl.DELETE("/delete/:id", transaction.DeleteTransaction)
		transactionsAuthUrl.POST("/listByUser", transaction.ListTransactionsByUser)
		transactionsAuthUrl.GET("/list/archived", transaction.ListArchivedTransactions)
		transactionsAuthUrl.POST("/send", transaction.SendTransaction)
		transactionsAuthUrl.POST("/dispute", transaction.CreateDispute)
		transactionsAuthUrl.GET("/dispute/fetch/:transaction_id", transaction.GetDisputeByTransactionID)
		transactionsAuthUrl.PATCH("/dispute/update", transaction.UpdateDispute)
		transactionsAuthUrl.GET("/list/user_disputes", transaction.GetDisputeByUser)
		transactionsAuthUrl.POST("/accept", transaction.AcceptTransaction)
		transactionsAuthUrl.POST("/delivered", transaction.TransactionDelivered)
		transactionsAuthUrl.POST("/reject", transaction.RejectTransaction)
		transactionsAuthUrl.POST("/reject_delivery", transaction.RejectTransactionDelivery)
		transactionsAuthUrl.POST("/request/due_date_extension", transaction.RequestDueDateExtension)
		transactionsAuthUrl.POST("/approve/due_date_extension", transaction.ApproveDueDateExtension)
		transactionsAuthUrl.POST("/satisfied", transaction.Satisfied)
		transactionsAuthUrl.PATCH("/updateStatus", transaction.UpdateTransactionStatus)
		transactionsAuthUrl.POST("/import", transaction.ImportTransactions)

	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.ApiType))
	{
		transactionsApiUrl.POST("/list", transaction.ListTransactions)
		transactionsApiUrl.GET("/listById/:id", transaction.ListTransactionsByID)
		transactionsApiUrl.GET("/list-transactions-by-ussd-code/:code", transaction.ListTransactionsByUSSDCode)
		transactionsApiUrl.POST("/listByBusiness", transaction.ListTransactionsByBusiness)
		transactionsApiUrl.POST("/listByBusinessFromMondayToThursday", transaction.ListByBusinessFromMondayToThursday)
		transactionsApiUrl.PATCH("/parties/update", transaction.UpdateTransactionParties)
		transactionsApiUrl.PATCH("/parties/update-status", transaction.UpdateTransactionPartyStatus)
		transactionsApiUrl.POST("/assign/buyer", transaction.AssignTransactionBuyer)
		transactionsApiUrl.PATCH("/broker/update", transaction.UpdateTransactionBroker)
		transactionsApiUrl.POST("/check-amount", transaction.CheckTransactionAmount)
		transactionsApiUrl.POST("/escrowcharge", transaction.GetEscrowCharge)
		transactionsApiUrl.GET("/rates", transaction.ListRates)
		transactionsApiUrl.GET("/exchange-transaction/:account_id", transaction.ListExchangeTransactionByAccountID)
		transactionsApiUrl.GET("exchange-transaction/show/:exchange_id", transaction.GetExchangeTransactionByID)
		transactionsApiUrl.PATCH("/api/updateStatus", transaction.UpdateTransactionStatusApi)
		transactionsApiUrl.POST("/api/satisfied", transaction.SatisfiedApi)

	}

	transactionsAppUrl := r.Group(fmt.Sprintf("%v", ApiVersion), middleware.Authorize(db, extReq, middleware.AppType))
	{
		transactionsAppUrl.POST("/validate_on_db", transaction.ValidateOnDB)
		transactionsAppUrl.PATCH("/update_transaction_amount_paid", transaction.UpdateTransactionAmountPaid)
		transactionsAppUrl.POST("/create_activity_log", transaction.CreateActivityLog)
		transactionsAppUrl.POST("/create_exchange_transaction", transaction.CreateExchangeTransaction)
		transactionsAppUrl.GET("/get_rate_by_currency/:from/:to", transaction.GetRateByFromAndToCurrencies)
		transactionsAppUrl.GET("/get_rate/:id", transaction.GetRateByID)
	}

	transactionsjobsUrl := r.Group(fmt.Sprintf("%v/jobs", ApiVersion))
	{
		transactionsjobsUrl.POST("/start", transaction.StartCronJob)
		transactionsjobsUrl.POST("/start-bulk", transaction.StartCronJobsBulk)
		transactionsjobsUrl.POST("/stop", transaction.StopCronJob)
		transactionsjobsUrl.PATCH("/update_interval", transaction.UpdateCronJobInterval)
	}
	return r
}
