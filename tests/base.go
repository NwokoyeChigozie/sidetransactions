package tests

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
	"github.com/vesicash/transactions-ms/external/mocks/payment_mocks"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/config"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/internal/models/migrations"
	"github.com/vesicash/transactions-ms/pkg/controller/transactions"
	"github.com/vesicash/transactions-ms/pkg/middleware"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func Setup() *utility.Logger {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../../app")
	db := postgresql.ConnectToDatabases(logger, config.TestDatabases)
	if config.TestDatabases.Migrate {
		migrations.RunAllMigrations(db)
	}
	return logger
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
func AssertBool(t *testing.T, got, expected bool) {
	if got != expected {
		t.Errorf("handler returned wrong boolean: got %v expected %v", got, expected)
	}
}

type ListTransactionsByIDResponse struct {
	Status  string                         `json:"status"`
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Data    models.TransactionByIDResponse `json:"data"`
}

func CreateTransactionUser(t *testing.T, db postgresql.Databases, validator *validator.Validate, extReq request.ExternalRequest, accountID int, isArchived bool) models.TransactionByIDResponse {
	var (
		trans                 = transactions.Controller{Db: db, Validator: validator, Logger: extReq.Logger, ExtReq: extReq}
		createTransactionPath = "/v2/transactions/create"
		createTransactionURI  = url.URL{Path: createTransactionPath}
		token, _              = uuid.NewV4()
		headers               = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token.String(),
		}
		r = gin.Default()
	)
	transactionsAuthUrl := r.Group(fmt.Sprintf("%v/transactions", "v2"), middleware.Authorize(db, trans.ExtReq, middleware.AuthType))
	{
		transactionsAuthUrl.POST("/create", trans.CreateTransaction)
	}

	payment_mocks.Payment = &external_models.Payment{
		ID:               int64(utility.GetRandomNumbersInRange(1000000000, 9999999999)),
		PaymentID:        utility.RandomString(20),
		TransactionID:    utility.RandomString(20),
		TotalAmount:      3000,
		EscrowCharge:     10,
		IsPaid:           false,
		AccountID:        int64(accountID),
		Currency:         "NGN",
		ShippingFee:      20,
		DisburseCurrency: "NGN",
	}

	createTransactionReq := models.CreateTransactionRequest{
		BusinessID: accountID,
		Parties: []models.Party{
			{
				AccountID:    accountID,
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
	}

	if isArchived {
		cc := createTransactionReq.Parties
		cc = append(cc, models.Party{
			AccountID:    accountID,
			EmailAddress: "sus@gmail.com",
			PhoneNumber:  "+2349876473847",
			Role:         "sender",
			Status:       "draft",
			AccessLevel: models.PartyAccessLevel{
				Approve:    true,
				CanReceive: false,
				CanView:    true,
				MarkAsDone: false,
			}})
		createTransactionReq.Parties = cc
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(createTransactionReq)
	req, err := http.NewRequest(http.MethodPost, createTransactionURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range headers {
		req.Header.Set(i, v)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	res := ListTransactionsByIDResponse{}
	json.NewDecoder(rr.Body).Decode(&res)
	return res.Data
}
