package transactions

import (
	"fmt"
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

	if transactionID == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "id not provided", c.Params, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

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

func (base *Controller) ListTransactions(c *gin.Context) {
	var (
		req       models.ListTransactionsRequest
		paginator = postgresql.GetPagination(c)
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

	transactions, pagination, code, err := transactions.ListTransactionsService(base.ExtReq, base.Logger, base.Db, req, paginator)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListTransactionsByBusiness(c *gin.Context) {
	var (
		req       models.ListTransactionByBusinessRequest
		paginator = postgresql.GetPagination(c)
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

	transactions, pagination, code, err := transactions.ListTransactionsByBusinessService(base.ExtReq, base.Logger, base.Db, req, paginator)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListByBusinessFromMondayToThursday(c *gin.Context) {
	var (
		req       models.ListByBusinessFromMondayToThursdayRequest
		paginator = postgresql.GetPagination(c)
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

	transactions, pagination, code, err := transactions.ListByBusinessFromMondayToThursdayService(base.ExtReq, base.Logger, base.Db, req, paginator)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListTransactionsByUser(c *gin.Context) {
	var (
		req       models.ListTransactionByUserRequest
		paginator = postgresql.GetPagination(c)
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

	transactions, pagination, code, err := transactions.ListTransactionsByUserService(base.ExtReq, base.Logger, base.Db, req, paginator, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	fmt.Println(pagination, pagination)
	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) ListArchivedTransactions(c *gin.Context) {
	var (
		paginator = postgresql.GetPagination(c)
	)

	user := models.MyIdentity
	if user == nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "error retrieving authenticated user", fmt.Errorf("error retrieving authenticated user"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	transactions, pagination, code, err := transactions.ListArchivedTransactionsService(base.ExtReq, base.Logger, base.Db, paginator, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}
