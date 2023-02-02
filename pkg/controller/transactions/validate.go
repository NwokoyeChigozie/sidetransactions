package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/services/transactions"
	"github.com/vesicash/transactions-ms/utility"
)

func (base *Controller) ValidateOnDB(c *gin.Context) {
	var (
		req models.ValidateOnDBReq
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

	status, code, err := transactions.ValidateOnDbService(req, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", status)
	c.JSON(http.StatusOK, rd)

}
