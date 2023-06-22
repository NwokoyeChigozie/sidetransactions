package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/vesicash/transactions-ms/external/external_models"
	"github.com/vesicash/transactions-ms/external/request"
	"github.com/vesicash/transactions-ms/internal/models"
	"github.com/vesicash/transactions-ms/pkg/repository/storage/postgresql"
	"github.com/vesicash/transactions-ms/utility"
)

func ListTransactionsByIDService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string) (models.TransactionByIDResponse, int, error) {
	var (
		transaction        = models.Transaction{TransactionID: transactionID}
		productTransaction = models.ProductTransaction{TransactionID: transactionID}
		parties            = []models.TransactionParty{}
		transactionFile    = models.TransactionFile{TransactionID: transactionID}
	)

	transactions, err := transaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, code, fmt.Errorf("transaction not found: %v", err.Error())
	}

	productTransactions, err := productTransaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if transaction.Type == "milestone" {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		parties, err = transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	} else {
		transactionParty := models.TransactionParty{TransactionID: transaction.TransactionID}
		parties, err = transactionParty.GetAllByTransactionID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	}

	transactionFiles, err := transactionFile.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionReponse := resolveTransactionAndListTransactionResponse(transaction)
	transactionReponse.Products = productTransactions

	partiess, members, err := getPartiesAndMembersFromParties(extReq, parties)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if len(parties) > 0 {
		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	} else {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		tParties, err := transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
		partiess, members, err = getPartiesAndMembersFromParties(extReq, tParties)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}

		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	}

	transactionReponse.Files = transactionFiles
	milestones := []models.MilestonesResponse{}

	type chanData struct {
		MileStoneSlice models.MilestonesResponse
		TotalAmount    float64
		Err            error
	}

	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactions))

	for i, t := range transactions {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, index int, transaction models.Transaction, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()

			totalAmount, mileSresponse := resolveTransactionForAmountAndMilestoneResponse(extReq, index, transaction)
			response := chanData{MileStoneSlice: mileSresponse, TotalAmount: totalAmount}

			ch <- response
		}(extReq, db, i, t, ch, &wg)

	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("error getting transactions", data.Err.Error())
		} else {
			if len(milestones) == 0 {
				transactionReponse.TotalAmount = data.TotalAmount
			}
			milestones = append(milestones, data.MileStoneSlice)
		}

	}

	sort.SliceStable(milestones, func(i, j int) bool {
		return milestones[i].Index < milestones[j].Index
	})

	transactionBroker := models.TransactionBroker{TransactionID: transactionID}
	activity := models.ActivityLog{TransactionID: transactionID}

	code, err = transactionBroker.GetTransactionBrokerByTransactionID(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		logger.Error("error getting transaction broker", err.Error())
	}

	activities, err := activity.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	country, err := GetCountryByNameOrCode(extReq, logger, transaction.Country)
	if err != nil && code == http.StatusInternalServerError {
		logger.Error("error getting country", err.Error())
	}

	dDateFormatted, err := utility.FormatDate(transaction.DueDate, "2006-01-02", "2006-01-02 15:04:05")
	if err != nil {
		dDateFormatted = ""
	}

	transactionState := models.TransactionState{TransactionID: transactionID, Status: "Closed"}
	code, err = transactionState.GetTransactionStateByTransactionIDAndStatus(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionDispute := models.TransactionDispute{TransactionID: transactionID}
	isDisputed, err := transactionDispute.IsDisputed(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	titleSlice := strings.Split(transactionReponse.Title, ";")
	transactionReponse.Title = titleSlice[0]
	transactionReponse.Milestones = milestones
	transactionReponse.Broker = transactionBroker
	transactionReponse.Activities = activities
	transactionReponse.Country = country
	transactionReponse.DueDateFormatted = dDateFormatted
	transactionReponse.TransactionClosedAt = transactionState.CreatedAt
	transactionReponse.IsDisputed = isDisputed
	return transactionReponse, http.StatusOK, nil
}
func ListTransactionsByIDLegacyService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string) (models.TransactionByIDResponse, int, error) {
	var (
		transaction        = models.Transaction{TransactionID: transactionID}
		productTransaction = models.ProductTransaction{TransactionID: transactionID}
		parties            = []models.TransactionParty{}
		transactionFile    = models.TransactionFile{TransactionID: transactionID}
	)

	transactions, err := transaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	code, err := transaction.GetTransactionByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, code, fmt.Errorf("transaction not found: %v", err.Error())
	}

	productTransactions, err := productTransaction.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if transaction.Type == "milestone" {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		parties, err = transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	} else {
		transactionParty := models.TransactionParty{TransactionID: transaction.TransactionID}
		parties, err = transactionParty.GetAllByTransactionID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
	}

	transactionFiles, err := transactionFile.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionReponse := resolveTransactionAndListTransactionResponse(transaction)
	transactionReponse.Products = productTransactions

	partiess, members, err := getPartiesAndMembersFromParties(extReq, parties)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	if len(parties) > 0 {
		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	} else {
		transactionParty := models.TransactionParty{TransactionPartiesID: transactions[0].PartiesID}
		tParties, err := transactionParty.GetAllByTransactionPartiesID(db.Transaction)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}
		partiess, members, err = getPartiesAndMembersFromParties(extReq, tParties)
		if err != nil {
			return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
		}

		transactionReponse.Parties = partiess
		transactionReponse.Members = members
	}

	transactionReponse.Files = transactionFiles
	milestones := map[int][]models.MilestonesResponse{}

	type chanData struct {
		MileStoneSlice []models.MilestonesResponse
		Index          int
		Err            error
	}

	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactions))
	wg.Add(len(transactions))
	for i, t := range transactions {
		go func(extReq request.ExternalRequest, db postgresql.Databases, index int, transaction models.Transaction, currentMileStoneSlice []models.MilestonesResponse, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			response := chanData{MileStoneSlice: currentMileStoneSlice, Index: index}
			totalAmount, mileSresponse := resolveTransactionForAmountAndMilestoneResponse(extReq, index, transaction)
			transactionReponse.TotalAmount = totalAmount
			currentMileStoneSlice = append(currentMileStoneSlice, mileSresponse)

			otherTransactions, err := transaction.GetAllOthersByIDAndPartiesID(db.Transaction)
			if err != nil {
				response.Err = err
			}

			for oi, ot := range otherTransactions {
				_, otherMileSresponse := resolveTransactionForAmountAndMilestoneResponse(extReq, oi, ot)
				currentMileStoneSlice = append(currentMileStoneSlice, otherMileSresponse)
			}
			response.MileStoneSlice = currentMileStoneSlice

			ch <- response
		}(extReq, db, i, t, milestones[i], ch, &wg)

	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("error getting other transactions", data.Err.Error())
		}
		milestones[data.Index] = data.MileStoneSlice
	}

	transactionBroker := models.TransactionBroker{TransactionID: transactionID}
	activity := models.ActivityLog{TransactionID: transactionID}

	code, err = transactionBroker.GetTransactionBrokerByTransactionID(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		logger.Error("error getting transaction broker", err.Error())
	}

	activities, err := activity.GetAllByTransactionID(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	country, err := GetCountryByNameOrCode(extReq, logger, transaction.Country)
	if err != nil && code == http.StatusInternalServerError {
		logger.Error("error getting country", err.Error())
	}

	dDateFormatted, err := utility.FormatDate(transaction.DueDate, "2006-01-02", "2006-01-02 15:04:05")
	if err != nil {
		dDateFormatted = ""
	}

	transactionState := models.TransactionState{TransactionID: transactionID, Status: "Closed"}
	code, err = transactionState.GetTransactionStateByTransactionIDAndStatus(db.Transaction)
	if err != nil && code == http.StatusInternalServerError {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	transactionDispute := models.TransactionDispute{TransactionID: transactionID}
	isDisputed, err := transactionDispute.IsDisputed(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, http.StatusInternalServerError, err
	}

	titleSlice := strings.Split(transactionReponse.Title, ";")
	transactionReponse.Title = titleSlice[0]
	transactionReponse.Milestones = milestones[0]
	transactionReponse.Broker = transactionBroker
	transactionReponse.Activities = activities
	transactionReponse.Country = country
	transactionReponse.DueDateFormatted = dDateFormatted
	transactionReponse.TransactionClosedAt = transactionState.CreatedAt
	transactionReponse.IsDisputed = isDisputed
	return transactionReponse, http.StatusOK, nil
}

func ListTransactionsByUssdCodeService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, ussdCode int) (models.TransactionByIDResponse, int, error) {
	var (
		transaction = models.Transaction{TransUssdCode: ussdCode}
	)

	code, err := transaction.GetTransactionByUssdCode(db.Transaction)
	if err != nil {
		return models.TransactionByIDResponse{}, code, err
	}

	return ListTransactionsByIDService(extReq, logger, db, transaction.TransactionID)
}

func ListTransactionsService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionsRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else if req.StatusCode != "" {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, false, req.Filter, "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	type chanData struct {
		TransactionResponse models.TransactionByIDResponse
		Err                 error
	}

	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactions))

	for _, t := range transactions {
		wg.Add(1)
		go func(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			response := chanData{}
			transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, transactionID)
			if err != nil {
				response.Err = err
			} else {
				response.TransactionResponse = transactionResponse
			}
			ch <- response
		}(extReq, logger, db, t.TransactionID, ch, &wg)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("list transaction by id error", data.Err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, data.TransactionResponse)
		}
	}
	sort.SliceStable(transactionsResponses, func(i, j int) bool {
		return transactionsResponses[i].ID > transactionsResponses[j].ID
	})

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListTransactionsByBusinessService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionByBusinessRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{BusinessID: req.BusinessID, IsPaylinked: req.Paylinked}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else if req.StatusCode != "" {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, req.Paylinked, req.Filter, "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	type chanData struct {
		TransactionResponse models.TransactionByIDResponse
		Err                 error
	}
	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactions))

	for _, t := range transactions {
		wg.Add(1)
		go func(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			response := chanData{}
			transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, transactionID)
			if err != nil {
				response.Err = err
			} else {
				payment, err := ListPayment(extReq, t.TransactionID)
				if err != nil {
					transactionResponse.TotalAmount = 0
					transactionResponse.EscrowCharge = 0
					logger.Error("list payment by transaction id error", err.Error())
				} else {
					transactionResponse.TotalAmount = payment.TotalAmount
					transactionResponse.EscrowCharge = payment.EscrowCharge
				}
				response.TransactionResponse = transactionResponse
			}
			ch <- response
		}(extReq, logger, db, t.TransactionID, ch, &wg)

	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("list transaction by id error or list payment by transaction id error", data.Err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, data.TransactionResponse)
		}
	}

	sort.SliceStable(transactionsResponses, func(i, j int) bool {
		return transactionsResponses[i].ID > transactionsResponses[j].ID
	})

	return transactionsResponses, pagination, http.StatusOK, nil

}
func ListByBusinessFromMondayToThursdayService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListByBusinessFromMondayToThursdayRequest, paginator postgresql.Pagination) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		transaction  = models.Transaction{BusinessID: req.BusinessID, IsPaylinked: req.Paylinked}
		transactions = []models.Transaction{}
	)

	if req.Status != "" {
		transaction.Status = req.Status
	} else if req.StatusCode != "" {
		transaction.Status = GetTransactionStatus(req.StatusCode)
	}

	transactions, pagination, err := transaction.GetAllByAndQueries(db.Transaction, req.Paylinked, "monday_to_thursday", "id", "asc", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	var transactionsResponses []models.TransactionByIDResponse

	type chanData struct {
		TransactionResponse models.TransactionByIDResponse
		Err                 error
	}
	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactions))
	for _, t := range transactions {
		wg.Add(1)
		go func(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			response := chanData{}
			transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, transactionID)
			if err != nil {
				response.Err = err
			} else {
				payment, err := ListPayment(extReq, t.TransactionID)
				if err != nil {
					transactionResponse.TotalAmount = 0
					transactionResponse.EscrowCharge = 0
					logger.Error("list payment by transaction id error", err.Error())
				} else {
					transactionResponse.TotalAmount = payment.TotalAmount
					transactionResponse.EscrowCharge = payment.EscrowCharge
				}
				response.TransactionResponse = transactionResponse
			}
			ch <- response
		}(extReq, logger, db, t.TransactionID, ch, &wg)

	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("list transaction by id error or list payment by transaction id error", data.Err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, data.TransactionResponse)
		}
	}

	sort.SliceStable(transactionsResponses, func(i, j int) bool {
		return transactionsResponses[i].ID > transactionsResponses[j].ID
	})

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListTransactionsByUserService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, req models.ListTransactionByUserRequest, paginator postgresql.Pagination, user external_models.User) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		// transactions          = []models.Transaction{}
		transactionsResponses = []models.TransactionByIDResponse{}
		transactionParty      = models.TransactionParty{AccountID: int(user.AccountID), Role: req.Role}
	)

	statusC := ""
	if req.StatusCode != "" {
		statusC = GetTransactionStatus(req.StatusCode)
	}

	transactionParties, pagination, err := transactionParty.GetAllByAndQueriesForUniqueValueForTransactionStatus(db.Transaction, "", "id", "desc", "transaction_id", statusC, paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	// for _, tp := range transactionParties {
	// 	lTransaction := models.Transaction{TransactionID: tp.TransactionID, IsPaylinked: req.Paylinked, Status: statusC}
	// 	code, err := lTransaction.GetLatestByAndQueries(db.Transaction, req.Paylinked, "")
	// 	if err != nil {
	// 		if code == http.StatusInternalServerError {
	// 			return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, code, err
	// 		}
	// 		logger.Error("error getting latest transaction by transaction id", err.Error())
	// 	} else {
	// 		transactions = append(transactions, lTransaction)
	// 	}
	// }

	type chanData struct {
		TransactionResponse models.TransactionByIDResponse
		Err                 error
	}

	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactionParties))

	for _, t := range transactionParties {
		wg.Add(1)
		go func(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, transactionID)
			response := chanData{Err: err, TransactionResponse: transactionResponse}
			ch <- response
		}(extReq, logger, db, t.TransactionID, ch, &wg)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("list transaction by id error", data.Err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, data.TransactionResponse)
		}
	}

	sort.SliceStable(transactionsResponses, func(i, j int) bool {
		return transactionsResponses[i].ID > transactionsResponses[j].ID
	})

	return transactionsResponses, pagination, http.StatusOK, nil

}

func ListArchivedTransactionsService(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, paginator postgresql.Pagination, user external_models.User) ([]models.TransactionByIDResponse, postgresql.PaginationResponse, int, error) {
	var (
		// transactions          = []models.Transaction{}
		transactionsResponses = []models.TransactionByIDResponse{}
		transactionParty      = models.TransactionParty{AccountID: int(user.AccountID), Role: "sender"}
	)

	transactionParties, pagination, err := transactionParty.GetAllByAndQueriesForUniqueValueForTransactionStatus(db.Transaction, "", "id", "desc", "transaction_id", "Deleted", paginator)
	if err != nil {
		return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	// for _, tp := range transactionParties {
	// 	lTransaction := models.Transaction{TransactionID: tp.TransactionID, Status: "Deleted"}
	// 	code, err := lTransaction.GetLatestByAndQueries(db.Transaction, false, "")
	// 	if err != nil {
	// 		if code == http.StatusInternalServerError {
	// 			return []models.TransactionByIDResponse{}, postgresql.PaginationResponse{}, code, err
	// 		}
	// 		logger.Error("error getting lastest transaction byy transaction id", err.Error())
	// 	} else {
	// 		transactions = append(transactions, lTransaction)
	// 	}
	// }

	type chanData struct {
		TransactionResponse models.TransactionByIDResponse
		Err                 error
	}

	var wg sync.WaitGroup
	ch := make(chan chanData, len(transactionParties))
	for _, t := range transactionParties {
		wg.Add(1)
		go func(extReq request.ExternalRequest, logger *utility.Logger, db postgresql.Databases, transactionID string, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			response := chanData{}
			transactionResponse, _, err := ListTransactionsByIDService(extReq, logger, db, transactionID)
			if err != nil {
				response.Err = err
			} else {
				response.TransactionResponse = transactionResponse
			}
			ch <- response
		}(extReq, logger, db, t.TransactionID, ch, &wg)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			logger.Error("list transaction by id error", data.Err.Error())
		} else {
			transactionsResponses = append(transactionsResponses, data.TransactionResponse)
		}
	}

	sort.SliceStable(transactionsResponses, func(i, j int) bool {
		return transactionsResponses[i].ID > transactionsResponses[j].ID
	})

	return transactionsResponses, pagination, http.StatusOK, nil

}

func resolveTransactionForAmountAndMilestoneResponse(extReq request.ExternalRequest, i int, t models.Transaction) (float64, models.MilestonesResponse) {
	var (
		totalAmount float64 = 0
	)
	var (
		titleSlice = strings.Split(t.Title, ";")
		title      = ""
		index      = 1
	)

	if len(titleSlice) >= 3 {
		totalAmount, _ = strconv.ParseFloat(titleSlice[2], 64)
	}

	if len(titleSlice) > 1 {
		title = titleSlice[1]
	}

	if len(titleSlice) >= 4 {
		index, _ = strconv.Atoi(titleSlice[3])
		if index == 0 {
			index = 1
		}
	}

	var recipients []models.MileStoneRecipient
	var recipientsResponse []models.MilestonesRecipientResponse

	fmt.Println("recipients string", t.Recipients)
	err := json.Unmarshal([]byte(t.Recipients), &recipients)
	if err != nil {
		extReq.Logger.Error("error unmarshaling recipient json string to struct", t.Recipients, err.Error())
	}

	for _, r := range recipients {
		user, _ := GetUserWithAccountID(extReq, r.AccountID)
		accountName := ""
		if user.ID != 0 {
			accountName = user.Lastname + " " + user.Firstname
		}

		recipientsResponse = append(recipientsResponse, models.MilestonesRecipientResponse{
			AccountID:   r.AccountID,
			AccountName: accountName,
			Email:       user.EmailAddress,
			PhoneNumber: user.PhoneNumber,
			Amount:      r.Amount,
		})

	}

	return totalAmount, models.MilestonesResponse{
		Index:            index,
		MilestoneID:      t.MilestoneID,
		Title:            title,
		Amount:           t.Amount,
		Status:           t.Status,
		InspectionPeriod: t.InspectionPeriod,
		DueDate:          t.DueDate,
		Recipients:       recipientsResponse,
	}
}

func getPartiesAndMembersFromParties(extReq request.ExternalRequest, parties []models.TransactionParty) (map[string]models.TransactionParty, []models.PartyResponse, error) {
	var (
		partiess = map[string]models.TransactionParty{}
		members  = []models.PartyResponse{}
		wg       sync.WaitGroup
	)
	type chanData struct {
		Member models.PartyResponse
		Err    error
	}

	ch := make(chan chanData, len(parties))
	wg.Add(len(parties))
	for _, p := range parties {
		partiess[p.Role] = p
		go func(party models.TransactionParty, ch chan chanData, wg *sync.WaitGroup) {
			defer wg.Done()
			user, _ := GetUserWithAccountID(extReq, party.AccountID)
			accountName := ""
			if user.ID != 0 {
				accountName = user.Lastname + " " + user.Firstname
			}
			var roleCapabilities models.PartyAccessLevel

			inrec, err := json.Marshal(party.RoleCapabilities)
			if err != nil {
				ch <- chanData{
					Member: models.PartyResponse{},
					Err:    err,
				}
				return
			}

			err = json.Unmarshal(inrec, &roleCapabilities)
			if err != nil {
				ch <- chanData{
					Member: models.PartyResponse{},
					Err:    err,
				}
				return
			}

			ch <- chanData{
				Member: models.PartyResponse{
					PartyID:     int(party.ID),
					AccountID:   party.AccountID,
					AccountName: accountName,
					PhoneNumber: user.PhoneNumber,
					Email:       user.EmailAddress,
					Role:        party.Role,
					Status:      party.Status,
					AccessLevel: roleCapabilities,
				},
				Err: nil,
			}

		}(p, ch, &wg)

	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for data := range ch {
		if data.Err != nil {
			return partiess, members, data.Err
		} else {
			members = append(members, data.Member)
		}
	}
	return partiess, members, nil
}

func resolveTransactionAndListTransactionResponse(transaction models.Transaction) models.TransactionByIDResponse {
	var rRrecipients []models.MileStoneRecipient
	json.Unmarshal([]byte(transaction.Recipients), &rRrecipients)
	return models.TransactionByIDResponse{
		ID:               transaction.ID,
		TransactionID:    transaction.TransactionID,
		PartiesID:        transaction.PartiesID,
		MilestoneID:      transaction.MilestoneID,
		BrokerID:         transaction.BrokerID,
		Title:            transaction.Title,
		Type:             transaction.Type,
		Description:      transaction.Description,
		Amount:           transaction.Amount,
		Status:           transaction.Status,
		Quantity:         transaction.Quantity,
		InspectionPeriod: transaction.InspectionPeriod,
		DueDate:          transaction.DueDate,
		ShippingFee:      transaction.ShippingFee,
		GracePeriod:      transaction.GracePeriod,
		Currency:         transaction.Currency,
		DeletedAt:        transaction.DeletedAt,
		CreatedAt:        transaction.CreatedAt,
		UpdatedAt:        transaction.UpdatedAt,
		BusinessID:       transaction.BusinessID,
		IsPaylinked:      transaction.IsPaylinked,
		Source:           transaction.Source,
		TransUssdCode:    transaction.TransUssdCode,
		Recipients:       rRrecipients,
		DisputeHandler:   transaction.DisputeHandler,
		AmountPaid:       transaction.AmountPaid,
		EscrowCharge:     transaction.EscrowCharge,
		EscrowWallet:     transaction.EscrowWallet,
	}
}
func resolveTransactionAndListTransactionResponse2(transaction models.Transaction) models.TransactionByIDResponse {
	v := models.TransactionByIDResponse{}
	utility.CopyStruct(&transaction, &v)
	return v
}
