package transactions

import (
	"fmt"
	"strings"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/utility"
)

func GetBusinessProfileByAccountID(extReq request.ExternalRequest, logger *utility.Logger, accountID int) (external_models.BusinessProfile, error) {
	businessProfileInterface, err := extReq.SendExternalRequest(request.GetBusinessProfile, external_models.GetBusinessProfileModel{
		AccountID: uint(accountID),
	})
	if err != nil {
		logger.Info(err.Error())
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
		extReq.Logger.Info(err.Error())
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
		extReq.Logger.Info(err.Error())
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
		logger.Info(err.Error())
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
		return external_models.ListPayment{}, fmt.Errorf("payment creation failed")
	}
	return payment, nil
}

func SignupUserWithPhone(extReq request.ExternalRequest, phone string) (external_models.User, error) {
	usItf, err := extReq.SendExternalRequest(request.SignupUser, external_models.CreateUserRequestModel{PhoneNumber: phone})
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
