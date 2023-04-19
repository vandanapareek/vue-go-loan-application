package accounting

import (
	"loan-api/internal/config"
	"loan-api/internal/dto"
	"loan-api/internal/models"
	xero "loan-api/internal/service/accounting/xero"
)

type AccountingService interface {
	GetProviders() ([]*models.Provider, error)
	GetBalanceSheet(req dto.BalanceSheetRequestParams) (dto.BalanceSheet, error)
}

type accountingService struct {
	repo models.ProviderRepo
}

func NewAccountingService(r models.ProviderRepo) AccountingService {
	return &accountingService{
		repo: r,
	}
}

func (s *accountingService) GetProviders() ([]*models.Provider, error) {
	result, err := s.repo.GetAllProviders()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *accountingService) GetBalanceSheet(reqData dto.BalanceSheetRequestParams) (dto.BalanceSheet, error) {
	//check which provider is used
	var response dto.BalanceSheet
	var err error
	switch reqData.AccountingProvider {
	case config.ACCOUNTING_PROVIDER_XERO:
		response, err = xero.RequestBalanceSheet(reqData)
	case config.ACCOUNTING_PROVIDER_MYOB:
		response, err = xero.RequestBalanceSheet(reqData)
	default:
		response, err = xero.RequestBalanceSheet(reqData)
	}

	if err != nil {
		return response, err
	}

	return response, nil
}
