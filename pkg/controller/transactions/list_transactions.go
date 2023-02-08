package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
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

func (base *Controller) ListTransactionsByBusiness(c *gin.Context) {
	var (
		req models.ListTransactionByBusinessRequest
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

	transactions, code, err := transactions.ListTransactionsByBusinessService(base.ExtReq, base.Logger, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions)
	c.JSON(http.StatusOK, rd)

}
