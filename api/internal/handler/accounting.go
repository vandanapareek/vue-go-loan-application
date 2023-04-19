package handler

import (
	"encoding/json"
	"fmt"
	"loan-api/errors"
	"loan-api/internal/dto"
	"loan-api/internal/service/accounting"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type accountingHandler struct {
	service accounting.AccountingService
}

func NewProviderHandler(s accounting.AccountingService) *accountingHandler {
	return &accountingHandler{
		service: s,
	}
}

type ApiSuccess struct {
	Code    int
	Message string
}

//get list of all Providers
func (h accountingHandler) GetProviders(writer http.ResponseWriter, request *http.Request) {
	//check if user logged in
	_, err := VerifyAuth(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnauthorized)
		return
	}

	//get all Providers

	details, err := h.service.GetProviders()
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

//get BalanceSheet
func (h accountingHandler) GetBalanceSheet(writer http.ResponseWriter, request *http.Request) {
	//check if user logged in
	_, err := VerifyAuth(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnauthorized)
		return
	}

	//validate params
	balanceSheetReq, err := processParams(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//get BalanceSheet
	details, err := h.service.GetBalanceSheet(balanceSheetReq)
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

func VerifyAuth(request *http.Request) (*Claims, error) {
	//check if user is authorized and logged in
	auth := request.Header.Get("Authorization")
	if auth == "" {
		return nil, errors.ErrUnauthorisedRequest
	}

	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	token, err := jwt.ParseWithClaims(auth, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	//wrong token or expired
	if err != nil || !token.Valid {
		fmt.Println("invalid token")
		return nil, errors.ErrTokenExpired
	}

	claims, _ := token.Claims.(*Claims)
	return claims, nil
}

func processParams(request *http.Request) (dto.BalanceSheetRequestParams, error) {
	var params dto.BalanceSheetRequestParams

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

	if params.AccountingProvider == "" {
		return params, errors.ErrProviderRequired
	}

	return params, nil
}
