package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) ListRates(c *gin.Context) {
	rate := models.Rate{}
	rates, err := rate.GetAll(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", rates)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListExchangeTransactionByAccountID(c *gin.Context) {
	var (
		accountID = c.Param("account_id")
		paginator = postgresql.GetPagination(c)
	)

	exchangeTransaction := models.ExchangeTransaction{AccountID: accountID}
	exchangeTransactions, pagination, err := exchangeTransaction.GetAllResolvedByAccountID(base.Db.Transaction, paginator)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", exchangeTransactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetExchangeTransactionByID(c *gin.Context) {
	var (
		idString = c.Param("exchange_id")
	)

	id, err := strconv.Atoi(idString)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid exchange id", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	exchangeTransaction := models.ExchangeTransaction{ID: uint(id)}
	rExchangeTransaction, code, err := exchangeTransaction.GetResolvedExchangeTransactionByID(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", rExchangeTransaction)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) CreateExchangeTransaction(c *gin.Context) {
	var (
		req models.CreateExchangeTransactionRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	vr := postgresql.ValidateRequestM{Logger: base.Logger, Test: base.ExtReq.Test}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	exchangeTransaction := models.ExchangeTransaction{
		AccountID:     strconv.Itoa(req.AccountID),
		InitialAmount: req.InitialAmount,
		FinalAmount:   req.FinalAmount,
		RateID:        req.RateID,
		Status:        req.Status,
	}

	err = exchangeTransaction.CreateExchangeTransaction(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "Exchange Transaction Created", nil)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) GetRateByID(c *gin.Context) {
	var (
		idString = c.Param("id")
	)

	id, err := strconv.Atoi(idString)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid rate id", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rate := models.Rate{ID: int64(id)}
	code, err := rate.GetRateByID(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", rate)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) GetRateByFromAndToCurrencies(c *gin.Context) {
	var (
		fromCurrency = c.Param("from")
		toCurrency   = c.Param("to")
	)

	rate := models.Rate{FromCurrency: fromCurrency, ToCurrency: toCurrency}
	code, err := rate.GetRateByFromAndToCurrencies(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", rate)
	c.JSON(http.StatusOK, rd)

}
