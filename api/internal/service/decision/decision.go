package decision

import (
	"loan-api/internal/dto"
)

type DecisionService interface {
	CalculateLoan(req dto.LoanDetails) (dto.LoanDetails, error)
}

type decisionService struct {
}

func NewDecisionService() DecisionService {
	return &decisionService{}
}

func (s *decisionService) CalculateLoan(reqData dto.LoanDetails) (dto.LoanDetails, error) {
	//set preAssessment value
	preAssessment := 20
	if reqData.TotalProfit > 0 && reqData.AvgAssets > reqData.LoanAmount {
		preAssessment = 100
	} else if reqData.TotalProfit > 0 {
		preAssessment = 60
	} else {
		//consider as loss or terminate
		reqData.PreAssessment = 0
		reqData.LoanPermitVal = 0
		return reqData, nil
	}
	reqData.PreAssessment = preAssessment
	reqData.LoanPermitVal = (reqData.LoanAmount * preAssessment) / 100

	//call decision engine API
	/*
		URL, _ := config.GetConf("DecisionEngineApiURL")
		form := url.Values{}
		form.Add("Name", reqData.BusinessName)
		form.Add("YearEstablished", strconv.Itoa(reqData.EstYear))
		form.Add("Summary", strconv.Itoa(reqData.TotalProfit))
		form.Add("PreAssessment", strconv.Itoa(preAssessment))
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
	return reqData, nil
}
