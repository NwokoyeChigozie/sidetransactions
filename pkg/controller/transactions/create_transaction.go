package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) CreateTransaction(c *gin.Context) {
	var (
		req models.CreateTransactionRequest
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

	transaction, code, err := transactions.CreateTransactionService(base.ExtReq, base.Logger, base.Db, req, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "Created", transaction)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) CheckTransactionAmount(c *gin.Context) {
	var (
		req models.OnlyTransactionIDRequiredRequest
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

	status, code, err := transactions.CheckTransactionAmountService(base.ExtReq, base.Logger, base.Db, req.TransactionID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", status)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) UpdateTransactionAmountPaid(c *gin.Context) {
	var (
		req models.UpdateTransactionAmountPaid
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

	transaction, code, err := transactions.UpdateTransactionAmountPaidService(base.ExtReq, base.Logger, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transaction)
	c.JSON(http.StatusOK, rd)

}
