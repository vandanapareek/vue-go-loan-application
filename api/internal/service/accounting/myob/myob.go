package myob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loan-api/internal/dto"
	"os"
	"strconv"
)

func RequestBalanceSheet(reqData dto.BalanceSheetRequestParams) (dto.BalanceSheet, error) {

	//call myob API to fetch balance sheet
	/*
		URL, _ := config.GetConf("MyobApiURL")
		form := url.Values{}
		form.Add("Name", reqData.BusinessName)
		form.Add("YearEstablished", strconv.Itoa(reqData.EstYear))
		form.Add("LoanAmount", strconv.Itoa(reqData.LoanAmount))
		form.Add("Currency", reqData.Currency)
		form.Add("StartDate", reqData.StartDate)
		form.Add("EndDate", reqData.EndDate)
		body := strings.NewReader(form.Encode())

		headerParams := externalapi.RequestHeader{
			ContentType: "application/x-www-form-urlencoded",
		}
		headerParamsJson, _ := json.Marshal(headerParams)
		apiResponse, err := externalapi.RestCall("POST", URL, body, string(headerParamsJson))
		if err != nil {
			return reqData, err
		}
	*/

	var details dto.BalanceSheet
	// Open jsonFile
	jsonFile, err := os.Open("balanceSheet/" + reqData.BusinessName + ".json")
	if err != nil {
		//try another dummy file
		jsonFile, err = os.Open("balanceSheet/business1.json")
		if err != nil {
			return details, err
		}
	}
	fmt.Println("File Successfully Opened")

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &details)
	totalProfit := 0
	avgAssets := 0
	for i := 0; i < len(details.Sheet); i++ {
		totalProfit = totalProfit + details.Sheet[i].ProfitOrLoss
		avgAssets = avgAssets + details.Sheet[i].AssetsValue
		fmt.Println("User Year: " + strconv.Itoa(details.Sheet[i].Year))
		fmt.Println("User Month: " + strconv.Itoa(details.Sheet[i].Month))
		fmt.Println("User ProfitOrLoss: " + strconv.Itoa(details.Sheet[i].ProfitOrLoss))
	}
	avgAssets = avgAssets / len(details.Sheet)
	fmt.Println("totalProfit: " + strconv.Itoa(totalProfit))
	fmt.Println("avgAssets: " + strconv.Itoa(avgAssets))
	details.TotalProfit = totalProfit
	details.AvgAssets = avgAssets
	details.LoanAmount = int(reqData.LoanAmount)

	// defer the closing of jsonFile
	defer jsonFile.Close()
	return details, nil
}
