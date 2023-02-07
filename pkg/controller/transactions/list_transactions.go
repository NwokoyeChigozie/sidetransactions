package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) ListTransactionsByID(c *gin.Context) {
	transactionID := c.Param("id")

	transactions, code, err := transactions.ListTransactionsByIDService(base.ExtReq, base.Logger, base.Db, transactionID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListTransactionsByUSSDCode(c *gin.Context) {
	uSSDCodeStr := c.Param("code")
	ussdCode, err := strconv.Atoi(uSSDCodeStr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "invalid ussd code", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	transactions, code, err := transactions.ListTransactionsByUssdCodeService(base.ExtReq, base.Logger, base.Db, ussdCode)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions)
	c.JSON(http.StatusOK, rd)

}
