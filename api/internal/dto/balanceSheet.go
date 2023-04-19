package dto

import "time"

type BalanceSheetRequestParams struct {
	BusinessName       string    `json:"businessName"` //Business registration no/gst no
	EstYear            int       `json:"estYear"`      //Year established
	LoanAmount         int       `json:"loanAmount"`   //total principal
	Currency           string    `json:"currency"`     //default SGD
	AccountingProvider string    `json:"provider"`
	LoanTerm           int       `json:"loanTerm"` //To be paid over (in months)
	InterestRate       float64   `json:"interestRate"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
}

type BalanceSheet struct {
	Sheet       []Sheet `json:"sheet"`
	TotalProfit int     `json:"totalProfit"`
	AvgAssets   int     `json:"avgAssets"`
	LoanAmount  int     `json:"loanAmount"`
}

type Sheet struct {
	Year         int `json:"year"`
	Month        int `json:"month"`
	ProfitOrLoss int `json:"profitOrLoss"`
	AssetsValue  int `json:"assetsValue"`
}
