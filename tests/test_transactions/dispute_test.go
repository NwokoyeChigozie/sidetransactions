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

func TestCreateDispute(t *testing.T) {
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
		RequestBody  models.CreateDisputeRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create dispute",
			RequestBody: models.CreateDisputeRequest{
				TransactionID: transaction.TransactionID,
				Reason:        "I want to",
				DisputeStatus: "ongoing",
				Decision:      "nothing",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Transaction Disputed",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.CreateDisputeRequest{
				TransactionID: "not correct",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.CreateDisputeRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/dispute", trans.CreateDispute)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/dispute"}

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
func TestGetDisputeByTransactionID(t *testing.T) {
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

	dispute := models.TransactionDispute{
		DisputeID:     utility.RandomString(20),
		TransactionID: transaction.TransactionID,
		Reason:        "whatever",
	}
	err := dispute.CreateTransactionDispute(db.Transaction)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name          string
		RequestBody   interface{}
		ExpectedCode  int
		Headers       map[string]string
		Message       string
		TransactionID string
	}{
		{
			Name:         "OK get dispute by transaction id",
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
			TransactionID: transaction.TransactionID,
		},
		{
			Name:         "incorrect transaction_id",
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
			TransactionID: "not-correct",
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.GET("/dispute/fetch/:transaction_id", trans.GetDisputeByTransactionID)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/dispute/fetch/" + test.TransactionID}

			req, err := http.NewRequest(http.MethodGet, URI.String(), &b)
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

func TestUpdateDispute(t *testing.T) {
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

	dispute := models.TransactionDispute{
		DisputeID:     utility.RandomString(20),
		TransactionID: transaction.TransactionID,
		Reason:        "whatever",
	}
	err := dispute.CreateTransactionDispute(db.Transaction)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name         string
		RequestBody  models.CreateDisputeRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update dispute",
			RequestBody: models.CreateDisputeRequest{
				TransactionID: transaction.TransactionID,
				Reason:        "I want to",
				DisputeStatus: "ongoing",
				Decision:      "nothing",
			},
			ExpectedCode: http.StatusOK,
			Message:      "Transaction Dispute Modified",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.CreateDisputeRequest{
				TransactionID: "not correct",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.CreateDisputeRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.PATCH("/dispute/update", trans.UpdateDispute)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/dispute/update"}

			req, err := http.NewRequest(http.MethodPatch, URI.String(), &b)
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

func TestGetDisputeByUser(t *testing.T) {
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

	dispute := models.TransactionDispute{
		DisputeID:     utility.RandomString(20),
		TransactionID: transaction.TransactionID,
		Reason:        "whatever",
	}
	err := dispute.CreateTransactionDispute(db.Transaction)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name         string
		RequestBody  interface{}
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name:         "OK get dispute by user",
			ExpectedCode: http.StatusOK,
			Message:      "Data Retrieved",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.GET("/list/user_disputes", trans.GetDisputeByUser)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/list/user_disputes"}

			req, err := http.NewRequest(http.MethodGet, URI.String(), &b)
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
