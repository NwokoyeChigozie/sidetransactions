package test_transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/mocks/auth_mocks"
	"github.com/vesicash/transactions-ms/external/mocks/payment_mocks"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/controller/transactions"
	"github.com/vesicash/transactions-ms/pkg/middleware"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/transactions-ms/tests"
	"github.com/vesicash/transactions-ms/utility"
)

func TestRequestDueDateExtension(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		token, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
		transactionID = utility.RandomString(20)
	)

	auth_mocks.User = &testUser
	auth_mocks.ValidateAuthorizationRes = &external_models.ValidateAuthorizationDataModel{
		Status:  true,
		Message: "authorized",
		Data:    testUser,
	}
	auth_mocks.UserProfile = &external_models.UserProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}

	auth_mocks.Country = &external_models.Country{
		ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		Name:         "nigeria",
		CountryCode:  "NG",
		CurrencyCode: "NGN",
	}

	auth_mocks.BusinessCharge = &external_models.BusinessCharge{
		ID:                  uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		BusinessId:          int(testUser.AccountID),
		Country:             "NG",
		Currency:            "NGN",
		BusinessCharge:      "0",
		VesicashCharge:      "2.5",
		ProcessingFee:       "0",
		PaymentGateway:      "rave",
		DisbursementGateway: "rave_momo",
		ProcessingFeeMode:   "fixed",
	}

	payment_mocks.Payment = &external_models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		PaymentID:        utility.RandomString(20),
		TransactionID:    transactionID,
		TotalAmount:      3000,
		EscrowCharge:     10,
		IsPaid:           false,
		AccountID:        int64(testUser.AccountID),
		Currency:         "NGN",
		ShippingFee:      20,
		DisburseCurrency: "NGN",
	}

	payment_mocks.ListPaymentObj = &external_models.ListPayment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		PaymentID:        payment_mocks.Payment.PaymentID,
		TransactionID:    transactionID,
		TotalAmount:      3000,
		EscrowCharge:     10,
		IsPaid:           false,
		AccountID:        int64(testUser.AccountID),
		Currency:         "NGN",
		ShippingFee:      20,
		DisburseCurrency: "NGN",
		SummedAmount:     4000,
	}

	trans := transactions.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	transaction := tst.CreateTransactionUser(t, db, validatorRef, trans.ExtReq, int(testUser.AccountID), false)
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.DueDateExtensionRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK request due date extension",
			RequestBody: models.DueDateExtensionRequest{
				TransactionID: transaction.TransactionID,
				Note:          "i feel like",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Transaction Due Date Extension Requested",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.DueDateExtensionRequest{
				TransactionID: "not correct",
				Note:          "i feel like",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.DueDateExtensionRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/request/due_date_extension", trans.RequestDueDateExtension)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/transactions/request/due_date_extension"}

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
func TestApproveDueDateExtension(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
		token, _  = uuid.NewV4()
		accountID = uint(utility.GetRandomNumbersInRange(1000000000, 9999999999))
		testUser  = external_models.User{
			ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
			AccountID:    accountID,
			EmailAddress: fmt.Sprintf("testuser%v@qa.team", muuid.String()),
			PhoneNumber:  fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			AccountType:  "individual",
			Firstname:    "test",
			Lastname:     "user",
			Username:     fmt.Sprintf("test_username%v", muuid.String()),
		}
		transactionID = utility.RandomString(20)
	)

	auth_mocks.User = &testUser
	auth_mocks.ValidateAuthorizationRes = &external_models.ValidateAuthorizationDataModel{
		Status:  true,
		Message: "authorized",
		Data:    testUser,
	}
	auth_mocks.UserProfile = &external_models.UserProfile{
		ID:        uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		AccountID: int(testUser.AccountID),
		Country:   "NG",
		Currency:  "NGN",
	}

	auth_mocks.Country = &external_models.Country{
		ID:           uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		Name:         "nigeria",
		CountryCode:  "NG",
		CurrencyCode: "NGN",
	}

	auth_mocks.BusinessCharge = &external_models.BusinessCharge{
		ID:                  uint(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		BusinessId:          int(testUser.AccountID),
		Country:             "NG",
		Currency:            "NGN",
		BusinessCharge:      "0",
		VesicashCharge:      "2.5",
		ProcessingFee:       "0",
		PaymentGateway:      "rave",
		DisbursementGateway: "rave_momo",
		ProcessingFeeMode:   "fixed",
	}

	payment_mocks.Payment = &external_models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		PaymentID:        utility.RandomString(20),
		TransactionID:    transactionID,
		TotalAmount:      3000,
		EscrowCharge:     10,
		IsPaid:           false,
		AccountID:        int64(testUser.AccountID),
		Currency:         "NGN",
		ShippingFee:      20,
		DisburseCurrency: "NGN",
	}

	payment_mocks.ListPaymentObj = &external_models.ListPayment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		PaymentID:        payment_mocks.Payment.PaymentID,
		TransactionID:    transactionID,
		TotalAmount:      3000,
		EscrowCharge:     10,
		IsPaid:           false,
		AccountID:        int64(testUser.AccountID),
		Currency:         "NGN",
		ShippingFee:      20,
		DisburseCurrency: "NGN",
		SummedAmount:     4000,
	}

	trans := transactions.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	transaction := tst.CreateTransactionUser(t, db, validatorRef, trans.ExtReq, int(testUser.AccountID), false)
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.ApproveDueDateExtensionRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK approve due date extension",
			RequestBody: models.ApproveDueDateExtensionRequest{
				TransactionID:    transaction.TransactionID,
				MilestoneID:      transaction.MilestoneID,
				DueDate:          "2023-11-15",
				InspectionPeriod: int(time.Now().Unix()),
			},
			ExpectedCode: http.StatusOK,
			Message:      "Transaction Due Date Extension Approved",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.ApproveDueDateExtensionRequest{
				TransactionID:    "not correct",
				MilestoneID:      transaction.MilestoneID,
				DueDate:          "2023-11-15",
				InspectionPeriod: int(time.Now().Unix()),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.ApproveDueDateExtensionRequest{
				TransactionID:    transaction.TransactionID,
				MilestoneID:      "not correct",
				DueDate:          "2023-11-15",
				InspectionPeriod: int(time.Now().Unix()),
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.ApproveDueDateExtensionRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/approve/due_date_extension", trans.ApproveDueDateExtension)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/transactions/approve/due_date_extension"}

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
