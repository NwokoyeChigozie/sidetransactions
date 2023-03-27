package test_transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestUpdateTransactionParties(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
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
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  models.UpdateTransactionPartiesRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update transaction parties",
			RequestBody: models.UpdateTransactionPartiesRequest{
				TransactionID: transaction.TransactionID,
				Parties: map[string]models.Party{
					"buyer": {
						AccountID:    int(testUser.AccountID),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "buyer",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
					"seller": {
						AccountID:    int(testUser.AccountID),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "buyer",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
				},
			},
			ExpectedCode: http.StatusOK,
			Message:      "updated",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.UpdateTransactionPartiesRequest{
				TransactionID: "not correct",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.UpdateTransactionPartiesRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.ApiType))
	{
		transactionsApiUrl.PATCH("/parties/update", trans.UpdateTransactionParties)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/parties/update"}

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
func TestUpdateTransactionPartyStatus(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
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
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  models.UpdateTransactionPartyStatusRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update transaction party status",
			RequestBody: models.UpdateTransactionPartyStatusRequest{
				TransactionID: transaction.TransactionID,
				AccountID:     int(testUser.AccountID),
				Status:        "draft",
			},
			ExpectedCode: http.StatusOK,
			Message:      "updated",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.UpdateTransactionPartyStatusRequest{
				TransactionID: "not correct",
				AccountID:     int(testUser.AccountID),
				Status:        "draft",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect account_id",
			RequestBody: models.UpdateTransactionPartyStatusRequest{
				TransactionID: transaction.TransactionID,
				AccountID:     0011,
				Status:        "draft",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.UpdateTransactionPartyStatusRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.ApiType))
	{
		transactionsApiUrl.PATCH("/parties/update-status", trans.UpdateTransactionPartyStatus)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/parties/update-status"}

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
func TestAssignTransactionBuyer(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
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
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	tests := []struct {
		Name         string
		RequestBody  models.AssignTransactionBuyerRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK assign buyer by transaction_id",
			RequestBody: models.AssignTransactionBuyerRequest{
				TransactionID: transaction.TransactionID,
				PhoneNumber:   testUser.PhoneNumber,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "OK assign buyer by ussd code",
			RequestBody: models.AssignTransactionBuyerRequest{
				UssdCode:    transaction.TransUssdCode,
				PhoneNumber: fmt.Sprintf("+234%v", utility.GetRandomNumbersInRange(7000000000, 9099999999)),
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.AssignTransactionBuyerRequest{
				TransactionID: "not correct",
				PhoneNumber:   testUser.PhoneNumber,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect ussd code",
			RequestBody: models.AssignTransactionBuyerRequest{
				UssdCode:    00022,
				PhoneNumber: testUser.PhoneNumber,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.AssignTransactionBuyerRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.ApiType))
	{
		transactionsApiUrl.POST("/assign/buyer", trans.AssignTransactionBuyer)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/assign/buyer"}

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
func TestUpdateTransactionBroker(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := postgresql.Connection()
	var (
		muuid, _  = uuid.NewV4()
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
	pvKey := utility.RandomString(20)
	pbKey := utility.RandomString(20)

	broker := models.TransactionBroker{TransactionID: transaction.TransactionID}

	code, err := broker.GetTransactionBrokerByTransactionID(db.Transaction)
	if err != nil {
		if code == http.StatusInternalServerError {
			t.Fatal(err)
		}
		broker = models.TransactionBroker{
			TransactionBrokerID: utility.RandomString(20),
			TransactionID:       transaction.TransactionID,
			BrokerCharge:        "11",
			BrokerChargeBearer:  strconv.Itoa(int(testUser.AccountID)),
			BrokerChargeType:    "fixed",
			IsSellerAccepted:    false,
			IsBuyerAccepted:     false,
		}
		err = broker.CreateTransactionBroker(db.Transaction)
		if err != nil {
			t.Fatal(err)
		}

	}

	falseV := false

	tests := []struct {
		Name         string
		RequestBody  models.UpdateTransactionBrokerRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update broker",
			RequestBody: models.UpdateTransactionBrokerRequest{
				TransactionID:      transaction.TransactionID,
				BrokerCharge:       "11",
				BrokerChargeBearer: strconv.Itoa(int(testUser.AccountID)),
				BrokerChargeType:   "fixed",
				IsSellerAccepted:   &falseV,
				IsBuyerAccepted:    &falseV,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.UpdateTransactionBrokerRequest{
				TransactionID:      "not correct",
				BrokerCharge:       "11",
				BrokerChargeBearer: strconv.Itoa(int(testUser.AccountID)),
				BrokerChargeType:   "fixed",
				IsSellerAccepted:   &falseV,
				IsBuyerAccepted:    &falseV,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.UpdateTransactionBrokerRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"v-private-key": pvKey,
				"v-public-key":  pbKey,
			},
		},
	}

	transactionsApiUrl := r.Group(fmt.Sprintf("%v", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.ApiType))
	{
		transactionsApiUrl.PATCH("/broker/update", trans.UpdateTransactionBroker)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/broker/update"}

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
