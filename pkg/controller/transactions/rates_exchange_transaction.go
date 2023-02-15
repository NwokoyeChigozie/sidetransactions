package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
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
	)

	exchangeTransaction := models.ExchangeTransaction{AccountID: accountID}
	exchangeTransactions, err := exchangeTransaction.GetAllResolvedByAccountID(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", exchangeTransactions)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetExchangeTransactionByID(c *gin.Context) {
	var (
		idString = c.Param("exchange_id")
	)

	id, err := strconv.Atoi(idString)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "invalid echange id", err, nil)
		c.JSON(http.StatusInternalServerError, rd)
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
