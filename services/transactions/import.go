package transactions

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func ImportTransactions(c *gin.Context, extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, user external_models.User) ([]models.Transaction, int, error) {
	var (
		columnLength = 11
		transactions = []models.Transaction{}
		mg512        = 512
	)
	code, err := ValidateUploadRequest(c, logger)
	if err != nil {
		return []models.Transaction{}, code, err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return []models.Transaction{}, http.StatusBadRequest, err
	}

	// Open the file
	src, err := fileHeader.Open()
	if err != nil {
		return []models.Transaction{}, http.StatusBadRequest, err
	}

	defer src.Close()

	buff := make([]byte, mg512)

	_, err = src.Read(buff)
	if err != nil {
		return []models.Transaction{}, http.StatusInternalServerError, err
	}

	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return []models.Transaction{}, http.StatusInternalServerError, err
	}

	// Parse the CSV file
	reader := csv.NewReader(src)
	rows, err := reader.ReadAll()
	if err != nil {
		return []models.Transaction{}, http.StatusInternalServerError, err
	}

	for index, row := range rows {
		var (
			duedate = row[8]
			amount  = row[10]
		)
		if len(row) != columnLength {
			return []models.Transaction{}, http.StatusBadRequest, fmt.Errorf("row %v has length either less than or greater than %v", index+1, columnLength)
		}
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return transactions, http.StatusBadRequest, fmt.Errorf("amount on row %v is not a valid number", index+1)
		}

		_, err = validateDueDate(duedate)
		if err != nil {
			return transactions, http.StatusBadRequest, fmt.Errorf("incorrect due date format for row %v, try 2006-01-15", index+1)
		}
	}

	for _, row := range rows {
		var (
			transactionID        = utility.RandomString(20)
			transactionPartiesID = utility.RandomString(20)
			title                = row[0]
			tType                = row[1]
			desc                 = row[2]
			buyer                = row[3]
			seller               = row[4]
			chargeBearer         = row[5]
			sender               = row[6]
			recipient            = row[7]
			duedate              = row[8]
			currency             = row[9]
			amount               = row[10]
		)

		if title != "" {
			var (
				errs    []error
				parties []models.TransactionParty
			)

			buyerAccount, err := GetUserWithEmail(extReq, buyer)
			if err != nil {
				errs = append(errs, err)
			}
			sellerAccount, err := GetUserWithEmail(extReq, seller)
			if err != nil {
				errs = append(errs, err)
			}
			chargeBearerAccount, _ := GetUserWithEmail(extReq, chargeBearer)
			senderAccount, _ := GetUserWithEmail(extReq, sender)
			recipientAccount, _ := GetUserWithEmail(extReq, recipient)

			if len(errs) < 1 {
				if buyerAccount.ID != 0 {
					parties = append(parties, models.TransactionParty{
						TransactionID:        transactionID,
						TransactionPartiesID: transactionPartiesID,
						AccountID:            int(buyerAccount.AccountID),
						Role:                 "buyer",
					})
				}
				if sellerAccount.ID != 0 {
					parties = append(parties, models.TransactionParty{
						TransactionID:        transactionID,
						TransactionPartiesID: transactionPartiesID,
						AccountID:            int(sellerAccount.AccountID),
						Role:                 "seller",
					})
				}
				if chargeBearerAccount.ID != 0 {
					parties = append(parties, models.TransactionParty{
						TransactionID:        transactionID,
						TransactionPartiesID: transactionPartiesID,
						AccountID:            int(chargeBearerAccount.AccountID),
						Role:                 "charge_bearer",
					})
				}
				if senderAccount.ID != 0 {
					parties = append(parties, models.TransactionParty{
						TransactionID:        transactionID,
						TransactionPartiesID: transactionPartiesID,
						AccountID:            int(senderAccount.AccountID),
						Role:                 "sender",
					})
				}
				if recipientAccount.ID != 0 {
					parties = append(parties, models.TransactionParty{
						TransactionID:        transactionID,
						TransactionPartiesID: transactionPartiesID,
						AccountID:            int(recipientAccount.AccountID),
						Role:                 "recipient",
					})
				}

				transactionParty := models.TransactionParty{}
				_, err = transactionParty.CreateTransactionsParties(db.Transaction, parties)
				if err != nil {
					logger.Info("error bulk creating transaction parties", err.Error())
				}

				amount, _ := strconv.ParseFloat(amount, 64)
				duedateUnix, _ := utility.GetUnixString(duedate, "2006-01-02", "2006-01-02")
				countryObj, _ := getCountryByCurrency(extReq, logger, currency)
				country := countryObj.CountryCode
				if country == "" {
					country = "NG"
				}
				transaction := models.Transaction{
					TransactionID:    transactionID,
					PartiesID:        transactionPartiesID,
					Title:            title,
					Type:             tType,
					Description:      desc,
					Amount:           amount,
					Status:           GetTransactionStatus("draft"),
					Quantity:         1,
					InspectionPeriod: duedateUnix,
					DueDate:          duedateUnix,
					ShippingFee:      0,
					GracePeriod:      duedateUnix,
					Currency:         strings.ToUpper(currency),
					Country:          strings.ToUpper(country),
					BusinessID:       user.BusinessId,
					TransUssdCode:    utility.GetRandomNumbersInRange(10000, 99999),
				}
				err := transaction.CreateTransaction(db.Transaction)
				if err != nil {
					return transactions, http.StatusInternalServerError, err
				}
				transactions = append(transactions, transaction)

			}

		}

	}

	return transactions, http.StatusOK, nil
}

func ValidateUploadRequest(c *gin.Context, logger *utility.Logger) (int, error) {
	var (
		maxSize          = 2097152
		maxSizeString    = fmt.Sprintf("%vMB", int(math.Floor(float64(maxSize)/1000000)))
		mg512            = 512
		allowedExtension = []string{".csv"}
	)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return http.StatusBadRequest, err
	}

	defer file.Close()

	buff := make([]byte, mg512)

	_, err = file.Read(buff)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	extension := filepath.Ext(fileHeader.Filename)
	if !contains(extension, allowedExtension) {
		return http.StatusBadRequest, fmt.Errorf("extension %v is not allowed", extension)
	}

	fileSize := fileHeader.Size
	if fileSize > int64(maxSize) {
		return http.StatusBadRequest, fmt.Errorf("File size is greater than %v.", maxSizeString)
	}

	return http.StatusOK, nil
}

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}

	return false
}
