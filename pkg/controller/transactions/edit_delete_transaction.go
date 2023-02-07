package transactions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) EditTransaction(c *gin.Context) {
	var (
		req models.EditTransactionRequest
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

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	transaction, code, err := transactions.EditTransactionService(base.ExtReq, base.Logger, base.Db, req, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Transaction Details Updated", transaction)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) DeleteTransaction(c *gin.Context) {
	var (
		transactionID = c.Param("id")
	)

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", fmt.Errorf("error retrieving authenticated user"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := transactions.DeleteTransactionService(base.ExtReq, base.Logger, base.Db, transactionID, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Transaction Deleted", nil)
	c.JSON(http.StatusOK, rd)

}
