package handler

import (
	"encoding/json"
	"fmt"
	"loan-api/errors"
	"loan-api/internal/dto"
	"loan-api/internal/service/decision"
	"net/http"
)

type decisionHandler struct {
	service decision.DecisionService
}

func NewDecisionHandler(s decision.DecisionService) *decisionHandler {
	return &decisionHandler{
		service: s,
	}
}

//Calculate Loan
func (h decisionHandler) CalculateLoan(writer http.ResponseWriter, request *http.Request) {
	//check if user logged in
	_, err := VerifyAuth(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnauthorized)
		return
	}

	//validate params
	balanceSheetReq, err := processLoanParams(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//get BalanceSheet
	details, err := h.service.CalculateLoan(balanceSheetReq)
	fmt.Println(err)
	fmt.Println(details)
	if err != nil {
		fmt.Println("err getting providers")
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//prepare output
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(details)
}

func processLoanParams(request *http.Request) (dto.LoanDetails, error) {
	var params dto.LoanDetails

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return params, errors.ErrDecodingRequest
	}

	if params.BusinessName == "" {
		return params, errors.ErrBusinessNameRequired
	}

	if params.EstYear == 0 {
		return params, errors.ErrEstYearRequired
	}

	if params.LoanAmount == 0 {
		return params, errors.ErrAmountRequired
	}

	if params.LoanAmount < 0 {
		return params, errors.ErrAmountRequired
	}

	var amount interface{} = params.LoanAmount
	if _, ok := amount.(float64); ok {
		return params, errors.ErrAmountRequired
	}

	if params.TotalProfit == 0 {
		return params, errors.ErrTotalProfitRequired
	}

	if params.AvgAssets == 0 {
		return params, errors.ErrAvgAssetsRequired
	}

	return params, nil
}
