package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) CreateActivityLog(c *gin.Context) {
	var (
		req models.CreateActivityLogRequest
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

	activityLog := models.ActivityLog{
		TransactionID: req.TransactionID,
		Description:   req.Description,
	}
	err = activityLog.CreateActivityLog(base.Db.Transaction)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", nil)
	c.JSON(http.StatusOK, rd)

}
