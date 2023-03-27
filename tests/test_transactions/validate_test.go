package test_transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/config"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/controller/transactions"
	"github.com/vesicash/transactions-ms/pkg/middleware"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/transactions-ms/tests"
	"github.com/vesicash/transactions-ms/utility"
)

func TestValidateOnDb(t *testing.T) {
	logger := tst.Setup()
	app := config.GetConfig().App
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()

	transaction := models.Transaction{
		TransactionID: utility.RandomString(20),
		PartiesID:     utility.RandomString(20),
		MilestoneID:   utility.RandomString(20),
		Title:         "test transaction;milestone title;2000;1",
		Type:          "milestone",
		Description:   "description",
		Amount:        2000,
		ShippingFee:   10,
		AmountPaid:    0,
		EscrowCharge:  10,
		EscrowWallet:  "yes",
	}
	err := transaction.CreateTransaction(db.Transaction)
	if err != nil {
		panic("error creating transaction: " + err.Error())
	}

	tests := []struct {
		Name              string
		RequestBody       models.ValidateOnDBReq
		ExpectedCode      int
		Headers           map[string]string
		Message           string
		CheckResponseData bool
		Response          bool
	}{
		{
			Name: "OK exists validate on db with value",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Type:  "exists",
				Query: "transaction_id = ?",
				Value: transaction.TransactionID,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK exists validate on db without value",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Type:  "exists",
				Query: "transaction_id = '" + transaction.TransactionID + "'",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK notexists validate on db with value",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Type:  "notexists",
				Query: "transaction_id = ?",
				Value: utility.RandomString(20) + "not",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          true,
		}, {
			Name: "OK exists validate on db without value",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Type:  "notexists",
				Query: "transaction_id = '" + transaction.TransactionID + "'",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
			CheckResponseData: true,
			Response:          false,
		}, {
			Name: "table omitted",
			RequestBody: models.ValidateOnDBReq{
				Type:  "notexists",
				Query: "transaction_id = '" + transaction.TransactionID + "'",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "type omitted",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Query: "transaction_id = '" + transaction.TransactionID + "'",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "query omitted",
			RequestBody: models.ValidateOnDBReq{
				Table: "transactions",
				Type:  "notexists",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name:         "no input",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	r := gin.Default()

	trans := transactions.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}

	transactionsAppUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AppType))
	{
		transactionsAppUrl.POST("/validate_on_db", trans.ValidateOnDB)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/validate_on_db"}

			req, err := http.NewRequest(http.MethodPost, URI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			for i, v := range test.Headers {
				req.Header.Set(i, v)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tst.ParseResponse(rr)

			code := int(data["code"].(float64))
			tst.AssertStatusCode(t, code, test.ExpectedCode)

			if test.CheckResponseData {
				resData := data["data"].(bool)
				tst.AssertBool(t, resData, test.Response)
			}

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tst.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tst.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
