package test_upload

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/controller/health"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	tst "github.com/vesicash/transactions-ms/tests"
)

func TestGetPing(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	// getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()
	requestURI := url.URL{Path: "/v2/transactions/health"}

	tests := []struct {
		Name         string
		ExpectedCode int
		RequestBody  *string
		Message      string
	}{
		{
			Name:         "OK",
			RequestBody:  nil,
			ExpectedCode: http.StatusOK,
			Message:      "ping successful",
		},
	}

	auth := health.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.GET("/v2/transactions/health", auth.Get)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodGet, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

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
func TestPostPing(t *testing.T) {
	logger := tst.Setup()
	gin.SetMode(gin.TestMode)
	// getConfig := config.GetConfig()
	validatorRef := validator.New()
	db := postgresql.Connection()
	requestURI := url.URL{Path: "/v2/transactions/health"}

	tests := []struct {
		Name         string
		ExpectedCode int
		RequestBody  models.Ping
		Message      string
	}{
		{
			Name:         "OK",
			RequestBody:  models.Ping{Message: "test"},
			ExpectedCode: http.StatusOK,
			Message:      "ping successful",
		},
	}

	auth := health.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range tests {
		r := gin.Default()

		r.POST("/v2/transactions/health", auth.Get)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			data := tst.ParseResponse(rr)

			tst.AssertStatusCode(t, rr.Code, test.ExpectedCode)

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
