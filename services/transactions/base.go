package transactions

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func GetBusinessProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.BusinessProfile, error) {
	businessProfileInterface, err := extReq.SendExternalRequest(request.GetBusinessProfile, external_models.GetBusinessProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Error(err.Error())
		return external_models.BusinessProfile{}, fmt.Errorf("Business lacks a profile.")
	}

	businessProfile, ok := businessProfileInterface.(external_models.BusinessProfile)
	if !ok {
		return external_models.BusinessProfile{}, fmt.Errorf("response data format error")
	}

	if businessProfile.ID == 0 {
		return external_models.BusinessProfile{}, fmt.Errorf("Business lacks a profile.")
	}
	return businessProfile, nil
}

func CreatePayment(extReq request.ExternalRequest, data external_models.CreatePaymentRequestWithToken) (external_models.Payment, error) {
	paymentInterface, err := extReq.SendExternalRequest(request.CreatePayment, data)
	if err != nil {
		return external_models.Payment{}, err
	}

	payment, ok := paymentInterface.(external_models.Payment)
	if !ok {
		return external_models.Payment{}, fmt.Errorf("response data format error")
	}

	if payment.ID == 0 {
		return external_models.Payment{}, fmt.Errorf("payment creation failed")
	}
	return payment, nil
}

func GetUserWithAccountID(extReq request.ExternalRequest, accountID int) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{AccountID: uint(accountID)})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}
func GetUserWithEmail(extReq request.ExternalRequest, email string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{EmailAddress: email})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}
func GetUserWithPhone(extReq request.ExternalRequest, phone string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.GetUserReq, external_models.GetUserRequestModel{PhoneNumber: phone})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}

func getBusinessChargeWithBusinessIDAndCurrency(extReq request.ExternalRequest, businessID int, currency string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.GetBusinessCharge, external_models.GetBusinessChargeModel{
		BusinessID: uint(businessID),
		Currency:   strings.ToUpper(currency),
	})

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.BusinessCharge{}, err
	}

	businessCharge, ok := dataInterface.(external_models.BusinessCharge)
	if !ok {
		return external_models.BusinessCharge{}, fmt.Errorf("response data format error")
	}

	if businessCharge.ID == 0 {
		return external_models.BusinessCharge{}, fmt.Errorf("business charge not found")
	}

	return businessCharge, nil
}
func initBusinessCharge(extReq request.ExternalRequest, businessID int, currency string) (external_models.BusinessCharge, error) {
	dataInterface, err := extReq.SendExternalRequest(request.InitBusinessCharge, external_models.InitBusinessChargeModel{
		BusinessID: uint(businessID),
		Currency:   strings.ToUpper(currency),
	})

	if err != nil {
		extReq.Logger.Error(err.Error())
		return external_models.BusinessCharge{}, err
	}

	businessCharge, ok := dataInterface.(external_models.BusinessCharge)
	if !ok {
		return external_models.BusinessCharge{}, fmt.Errorf("response data format error")
	}

	if businessCharge.ID == 0 {
		return external_models.BusinessCharge{}, fmt.Errorf("business charge init failed")
	}

	return businessCharge, nil
}
func GetCountryByNameOrCode(extReq request.ExternalRequest, logger *utility.Logger, NameOrCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		Name: NameOrCode,
	})

	if err != nil {
		logger.Error(err.Error())
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile.")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile")
	}

	return country, nil
}
func getCountryByCurrency(extReq request.ExternalRequest, logger *utility.Logger, currencyCode string) (external_models.Country, error) {

	countryInterface, err := extReq.SendExternalRequest(request.GetCountry, external_models.GetCountryModel{
		CurrencyCode: currencyCode,
	})

	if err != nil {
		logger.Error(err.Error())
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile.")
	}
	country, ok := countryInterface.(external_models.Country)
	if !ok {
		return external_models.Country{}, fmt.Errorf("response data format error")
	}
	if country.ID == 0 {
		return external_models.Country{}, fmt.Errorf("Your country could not be resolved, please update your profile")
	}

	return country, nil
}

func ListPayment(extReq request.ExternalRequest, transactionID string) (external_models.ListPayment, error) {
	paymentInterface, err := extReq.SendExternalRequest(request.ListPayment, transactionID)
	if err != nil {
		return external_models.ListPayment{}, err
	}

	payment, ok := paymentInterface.(external_models.ListPayment)
	if !ok {
		return external_models.ListPayment{}, fmt.Errorf("response data format error")
	}

	if payment.ID == 0 {
		return external_models.ListPayment{}, fmt.Errorf("payment listing failed")
	}
	return payment, nil
}

func SignupUserWithPhone(extReq request.ExternalRequest, phone, accountType string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.SignupUser, external_models.CreateUserRequestModel{PhoneNumber: phone, AccountType: accountType})
	if err != nil {
		return external_models.User{}, err
	}

	us, ok := usItf.(external_models.User)
	if !ok {
		return external_models.User{}, fmt.Errorf("response data format error")
	}
	fmt.Println("user", us)
	fmt.Println("user.id", us.ID)

	if us.ID == 0 {
		return external_models.User{}, fmt.Errorf("user not found")
	}
	return us, nil
}

func GetTransactionStatus(index string) string {
	dataMap := map[string]string{
		"":        "Draft",
		"sac":     "Sent - Awaiting Confirmation",
		"sr":      "Sent - Rejected",
		"af":      "Accepted - Funded",
		"anf":     "Accepted - Not Funded",
		"fr":      "Funded - Rejected",
		"ip":      "In Progress",
		"d":       "Delivered",
		"da":      "Delivered - Accepted",
		"dr":      "Delivered - Rejected",
		"cdp":     "Closed - Disbursement Pending",
		"cmdp":    "Closed - Manual Disbursement Pending",
		"cdc":     "Closed - Disbursement Complete",
		"cd":      "Closed - Disputed",
		"cnf":     "Closed - Not Funded",
		"closed":  "Closed",
		"draft":   "Draft",
		"active":  "Active",
		"cr":      "Closed - Refunded",
		"deleted": "Deleted",
	}
	status := dataMap[strings.ToLower(index)]
	if status == "" {
		status = dataMap[""]
	}
	return status
}

func CheckTransactionStatus(index string) bool {
	dataMap := map[string]string{
		"sac":     "Sent - Awaiting Confirmation",
		"sr":      "Sent - Rejected",
		"af":      "Accepted - Funded",
		"anf":     "Accepted - Not Funded",
		"fr":      "Funded - Rejected",
		"ip":      "In Progress",
		"d":       "Delivered",
		"da":      "Delivered - Accepted",
		"dr":      "Delivered - Rejected",
		"cdp":     "Closed - Disbursement Pending",
		"cmdp":    "Closed - Manual Disbursement Pending",
		"cdc":     "Closed - Disbursement Complete",
		"cd":      "Closed - Disputed",
		"cnf":     "Closed - Not Funded",
		"closed":  "Closed",
		"draft":   "Draft",
		"active":  "Active",
		"cr":      "Closed - Refunded",
		"deleted": "Deleted",
	}
	return dataMap[strings.ToLower(index)] != ""
}

func GetAccessTokenByKeyFromRequest(extReq request.ExternalRequest, c *gin.Context) (external_models.AccessToken, error) {
	privateKey := utility.GetHeader(c, "v-private-key")
	publicKey := utility.GetHeader(c, "v-public-key")
	key := privateKey
	if key == "" {
		key = publicKey
	}
	acItf, err := extReq.SendExternalRequest(request.GetAccessTokenByKey, key)
	if err != nil {
		return external_models.AccessToken{}, err
	}

	accessToken, ok := acItf.(external_models.AccessToken)
	if !ok {
		return external_models.AccessToken{}, fmt.Errorf("response data format error")
	}

	return accessToken, nil
}

func DebitWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, creditEscrow string, creditMor string, transactionID string) (external_models.WalletBalance, error) {
	walletItf, err := extReq.SendExternalRequest(request.DebitWallet, external_models.DebitWalletRequest{
		Amount:        amount,
		Currency:      currency,
		BusinessID:    businessID,
		EscrowWallet:  creditEscrow,
		MorWallet:     creditMor,
		TransactionID: transactionID,
	})
	if err != nil {
		return external_models.WalletBalance{}, err
	}

	walletBalance, ok := walletItf.(external_models.WalletBalance)
	if !ok {
		return external_models.WalletBalance{}, fmt.Errorf("response data format error")
	}

	return walletBalance, nil
}

func CreditWallet(extReq request.ExternalRequest, db postgresql.Databases, amount float64, currency string, businessID int, isRefund bool, creditEscrow string, creditMor string, transactionID string) (external_models.WalletBalance, error) {
	walletItf, err := extReq.SendExternalRequest(request.CreditWallet, external_models.CreditWalletRequest{
		Amount:        amount,
		Currency:      currency,
		BusinessID:    businessID,
		EscrowWallet:  creditEscrow,
		MorWallet:     creditMor,
		TransactionID: transactionID,
		IsRefund:      isRefund,
	})
	if err != nil {
		return external_models.WalletBalance{}, err
	}

	walletBalance, ok := walletItf.(external_models.WalletBalance)
	if !ok {
		return external_models.WalletBalance{}, fmt.Errorf("response data format error")
	}

	return walletBalance, nil
}
