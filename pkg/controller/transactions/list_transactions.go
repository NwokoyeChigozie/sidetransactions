package transactions

import (
	"net/http"

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
