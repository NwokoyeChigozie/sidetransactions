package transactions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) ImportTransactions(c *gin.Context) {
	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", fmt.Errorf("error retrieving authenticated user"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	transactions, code, err := transactions.ImportTransactions(c, base.ExtReq, base.Logger, base.Db, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Imported", transactions)
	c.JSON(http.StatusOK, rd)
}
