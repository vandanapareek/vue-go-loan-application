package dto

type LoanDetails struct {
	BusinessName  string `json:"businessName"` //Business registration no/gst no
	EstYear       int    `json:"estYear"`      //Year established
	LoanAmount    int    `json:"loanAmount"`   //total principal
	TotalProfit   int    `json:"totalProfit"`
	AvgAssets     int    `json:"avgAssets"`
	PreAssessment int    `json:"preAssessment"`
	LoanPermitVal int    `json:"loanPermitVal"`
}
