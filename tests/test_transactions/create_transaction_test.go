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
	"github.com/vesicash/transactions-ms/internal/config"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/controller/transactions"
	"github.com/vesicash/transactions-ms/pkg/middleware"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/transactions-ms/tests"
	"github.com/vesicash/transactions-ms/utility"
)

func TestCreateTransaction(t *testing.T) {
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

	trans := transactions.Controller{Db: db, Validator: validatorRef, Logger: logger, ExtReq: request.ExternalRequest{
		Logger: logger,
		Test:   true,
	}}
	r := gin.Default()

	tests := []struct {
		Name         string
		RequestBody  models.CreateTransactionRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK create transaction",
			RequestBody: models.CreateTransactionRequest{
				BusinessID: int(testUser.AccountID),
				Parties: []models.Party{
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
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
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "seller",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
				},
				Title:        "test title",
				Type:         "milestone",
				EscrowWallet: "yes",
				Description:  "transaction description",
				Files: []models.File{
					{
						Name: "file name",
						URL:  "https://linktofile.com",
					},
				},
				Milestones: []models.MileStone{
					{
						Title:            "milestone title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
					{
						Title:            "milestone 2 title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
				},
				Quantity:         1,
				Amount:           2000,
				InspectionPeriod: 2,
				GracePeriod:      "2023-04-18",
				DueDate:          "2023-04-15",
				ShippingFee:      0,
				Currency:         "NGN",
				Source:           "transfer",
				Paylinked:        false,
			},
			ExpectedCode: http.StatusCreated,
			Message:      "Created",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no party",
			RequestBody: models.CreateTransactionRequest{
				BusinessID:   int(testUser.AccountID),
				Title:        "test title",
				Type:         "milestone",
				EscrowWallet: "yes",
				Description:  "transaction description",
				Files: []models.File{
					{
						Name: "file name",
						URL:  "https://linktofile.com",
					},
				},
				Milestones: []models.MileStone{
					{
						Title:            "milestone title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
					{
						Title:            "milestone 2 title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
				},
				Quantity:         1,
				Amount:           2000,
				InspectionPeriod: 2,
				GracePeriod:      "2023-04-18",
				DueDate:          "2023-04-15",
				ShippingFee:      0,
				Currency:         "NGN",
				Source:           "transfer",
				Paylinked:        false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no title",
			RequestBody: models.CreateTransactionRequest{
				BusinessID: int(testUser.AccountID),
				Parties: []models.Party{
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
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
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "seller",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
				},
				Type:         "milestone",
				EscrowWallet: "yes",
				Description:  "transaction description",
				Files: []models.File{
					{
						Name: "file name",
						URL:  "https://linktofile.com",
					},
				},
				Milestones: []models.MileStone{
					{
						Title:            "milestone title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
					{
						Title:            "milestone 2 title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
				},
				Quantity:         1,
				Amount:           2000,
				InspectionPeriod: 2,
				GracePeriod:      "2023-04-18",
				DueDate:          "2023-04-15",
				ShippingFee:      0,
				Currency:         "NGN",
				Source:           "transfer",
				Paylinked:        false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no type",
			RequestBody: models.CreateTransactionRequest{
				BusinessID: int(testUser.AccountID),
				Parties: []models.Party{
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
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
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "seller",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
				},
				Title:        "test title",
				EscrowWallet: "yes",
				Description:  "transaction description",
				Files: []models.File{
					{
						Name: "file name",
						URL:  "https://linktofile.com",
					},
				},
				Milestones: []models.MileStone{
					{
						Title:            "milestone title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
					{
						Title:            "milestone 2 title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
				},
				Quantity:         1,
				Amount:           2000,
				InspectionPeriod: 2,
				GracePeriod:      "2023-04-18",
				DueDate:          "2023-04-15",
				ShippingFee:      0,
				Currency:         "NGN",
				Source:           "transfer",
				Paylinked:        false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "no currency",
			RequestBody: models.CreateTransactionRequest{
				BusinessID: int(testUser.AccountID),
				Parties: []models.Party{
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
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
					{
						AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
						EmailAddress: "sus@gmail.com",
						PhoneNumber:  "+2349876473847",
						Role:         "seller",
						Status:       "draft",
						AccessLevel: models.PartyAccessLevel{
							Approve:    true,
							CanReceive: false,
							CanView:    true,
							MarkAsDone: false,
						},
					},
				},
				Title:        "test title",
				Type:         "milestone",
				EscrowWallet: "yes",
				Description:  "transaction description",
				Files: []models.File{
					{
						Name: "file name",
						URL:  "https://linktofile.com",
					},
				},
				Milestones: []models.MileStone{
					{
						Title:            "milestone title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
					{
						Title:            "milestone 2 title",
						Amount:           1000,
						InspectionPeriod: 4,
						DueDate:          "2023-03-16",
						Status:           "draft",
						Description:      "milestone description",
						Quantity:         1,
						ShippingFee:      0,
						GracePeriod:      "2023-03-18",
						Recipients: []models.MileStoneRecipient{
							{
								AccountID:    utility.GetRandomNumbersInRange(1000000000, 9999999999),
								Amount:       500,
								EmailAddress: "sus@gmail.com",
								PhoneNumber:  "+23456789776789",
							},
						},
					},
				},
				Quantity:         1,
				Amount:           2000,
				InspectionPeriod: 2,
				GracePeriod:      "2023-04-18",
				DueDate:          "2023-04-15",
				ShippingFee:      0,
				Source:           "transfer",
				Paylinked:        false,
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.CreateTransactionRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/create", trans.CreateTransaction)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/transactions/create"}

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

func TestCheckTransactionAmount(t *testing.T) {
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
		RequestBody  models.OnlyTransactionIDRequiredRequest
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK check transaction amount",
			RequestBody: models.OnlyTransactionIDRequiredRequest{
				TransactionID: transaction.TransactionID,
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.OnlyTransactionIDRequiredRequest{
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
			RequestBody:  models.OnlyTransactionIDRequiredRequest{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer " + token.String(),
			},
		},
	}

	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/check-amount", trans.CheckTransactionAmount)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/transactions/check-amount"}

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
func TestUpdateTransactionAmountPaid(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	app := config.GetConfig().App
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

	tests := []struct {
		Name         string
		RequestBody  models.UpdateTransactionAmountPaid
		ExpectedCode int
		Headers      map[string]string
		Message      string
	}{
		{
			Name: "OK update transaction amount +",
			RequestBody: models.UpdateTransactionAmountPaid{
				TransactionID: transaction.TransactionID,
				Amount:        200,
				Action:        "+",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "OK update transaction amount -",
			RequestBody: models.UpdateTransactionAmountPaid{
				TransactionID: transaction.TransactionID,
				Amount:        200,
				Action:        "-",
			},
			ExpectedCode: http.StatusOK,
			Message:      "successful",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		}, {
			Name: "wrong action *",
			RequestBody: models.UpdateTransactionAmountPaid{
				TransactionID: transaction.TransactionID,
				Amount:        200,
				Action:        "*",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name: "incorrect transaction_id",
			RequestBody: models.UpdateTransactionAmountPaid{
				TransactionID: "not correct",
			},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
		{
			Name:         "empty request",
			RequestBody:  models.UpdateTransactionAmountPaid{},
			ExpectedCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"v-app":        app.Key,
			},
		},
	}

	transactionsAppUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AppType))
	{
		transactionsAppUrl.PATCH("/update_transaction_amount_paid", trans.UpdateTransactionAmountPaid)
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)
			URI := url.URL{Path: "/v2/transactions/update_transaction_amount_paid"}

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
