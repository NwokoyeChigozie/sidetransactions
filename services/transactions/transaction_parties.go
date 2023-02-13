package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func UpdateTransactionPartiesService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.UpdateTransactionPartiesRequest) (int, error) {

	var (
		transactionID = models.Transaction{TransactionID: req.TransactionID}
		users         = map[string]external_models.User{}
	)

	if len(req.Parties) < 1 {
		return http.StatusBadRequest, fmt.Errorf("parties not specified")
	}
	for name, party := range req.Parties {
		if party.AccountID == 0 {
			return http.StatusBadRequest, fmt.Errorf("%v has no accountID", name)
		}

		user, err := GetUserWithAccountID(extReq, party.AccountID)
		if err != nil {
			return http.StatusBadRequest, err
		}

		users[name] = user
	}

	code, err := transactionID.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return code, err
	}

	for name, party := range req.Parties {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactionID.PartiesID, Role: name}
		code, err := transactionParty.GetTransactionPartyByTransactionPartiesIDAndRole(db.Transaction)
		if err != nil {
			if code == http.StatusInternalServerError {
				return http.StatusInternalServerError, err
			}
		} else {
			var roleCapabilities map[string]interface{}
			inrec, err := json.Marshal(&party.AccessLevel)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			err = json.Unmarshal(inrec, &roleCapabilities)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			transactionParty.AccountID = party.AccountID
			transactionParty.RoleCapabilities = roleCapabilities
			err = transactionParty.UpdateAllFields(db.Transaction)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

	}
	return http.StatusOK, nil
}

func UpdateTransactionPartyStatusService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.UpdateTransactionPartyStatusRequest) (int, error) {

	var (
		transactionParty = models.TransactionParty{TransactionID: req.TransactionID, AccountID: req.AccountID}
	)
	user, err := GetUserWithAccountID(extReq, req.AccountID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	code, err := transactionParty.GetTransactionPartyByTransactionIDAndAccountID(db.Transaction)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}

		return code, fmt.Errorf("this transaction has no party with this account id")
	}

	transactionParty.Status = req.Status
	err = transactionParty.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	activityLog := models.ActivityLog{
		TransactionID: req.TransactionID,
		Description:   fmt.Sprintf("%v has %v transaction invitation", user.EmailAddress, req.Status),
	}
	err = activityLog.CreateActivityLog(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func AssignTransactionBuyerService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.AssignTransactionBuyerRequest) (int, error) {
	var (
		phoneNumber, _ = utility.PhoneValid(req.PhoneNumber)
		transaction    = models.Transaction{TransactionID: req.TransactionID, TransUssdCode: req.UssdCode}
		roles          = []string{"buyer", "recipient", "charge_bearer"}
	)

	user, err := GetUserWithPhone(extReq, phoneNumber)
	if err != nil {
		us, er := SignupUserWithPhone(extReq, phoneNumber)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("Oops, No User exist with that phone number. %v. %v", err, er)
		}
		user = us
	}

	if req.TransactionID != "" {
		code, err := transaction.GetTransactionByTransactionID(db.Transaction)
		if err != nil {
			return code, err
		}
	} else if req.UssdCode != 0 {
		code, err := transaction.GetTransactionByUssdCode(db.Transaction)
		if err != nil {
			return code, err
		}
	} else {
		return http.StatusBadRequest, fmt.Errorf("provide either transaction id or ussd code")
	}

	for _, role := range roles {
		party := models.TransactionParty{TransactionPartiesID: transaction.PartiesID, Role: role}
		code, err := party.GetTransactionPartyByTransactionPartiesIDAndRole(db.Transaction)
		if err != nil {
			if code == http.StatusInternalServerError {
				return code, err
			}
			party.AccountID = int(user.AccountID)
			party.TransactionID = party.TransactionID
			party.TransactionPartiesID = transaction.PartiesID
			party.Role = "buyer"
			err = party.CreateTransactionParty(db.Transaction)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			party.AccountID = int(user.AccountID)
			err = party.UpdateAllFields(db.Transaction)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}

	return http.StatusOK, nil
}

func UpdateTransactionBrokerService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.UpdateTransactionBrokerRequest) (int, error) {
	var (
		broker = models.TransactionBroker{TransactionID: req.TransactionID}
	)

	code, err := broker.GetTransactionBrokerByTransactionID(db.Transaction)
	if err != nil {
		if code == http.StatusInternalServerError {
			return code, err
		}
		return http.StatusBadRequest, fmt.Errorf("there is broker for this transaction")

	}

	if req.BrokerCharge != "" {
		broker.BrokerCharge = req.BrokerCharge
	}

	if req.BrokerChargeBearer != "" {
		broker.BrokerChargeBearer = req.BrokerChargeBearer
	}

	if req.BrokerChargeType != "" {
		broker.BrokerChargeType = req.BrokerChargeType
	}

	if req.IsBuyerAccepted != nil {
		broker.IsBuyerAccepted = *req.IsBuyerAccepted
	}

	if req.IsSellerAccepted != nil {
		broker.IsSellerAccepted = *req.IsSellerAccepted
	}

	err = broker.UpdateAllFields(db.Transaction)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// $broker->broker_charge = (!empty($request->broker_charge) ? $request->broker_charge : $broker->broker_charge);

//             $broker->broker_charge_bearer = (!empty($request->broker_charge_bearer) ? $request->broker_charge_bearer : $broker->broker_charge_bearer);

//             $broker->broker_charge_type = (!empty($request->broker_charge_type) ? $request->broker_charge_type : $broker->broker_charge_type);

//             $broker->is_seller_accepted = (!empty($request->is_seller_accepted) ? $request->is_seller_accepted : $broker->is_seller_accepted);

//             $broker->is_buyer_accepted = (!empty($request->is_buyer_accepted) ? $request->is_buyer_accepted : $broker->is_buyer_accepted);
