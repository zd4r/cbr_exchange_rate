package usecase

import (
	"fmt"

	"github.com/zd4r/cbr_exchange_rate/internal/entity"
)

type DynamicQuoteUseCase struct {
	webAPI DynamicQuoteWebAPI
}

func NewDynamicQuote(w DynamicQuoteWebAPI) *DynamicQuoteUseCase {
	return &DynamicQuoteUseCase{
		webAPI: w,
	}
}

func (uc *DynamicQuoteUseCase) Calculate(dqs entity.DynamicQuotes) (entity.DynamicQuotes, error) {
	currencies, err := uc.webAPI.GetCurrencies(dqs)
	if err != nil {
		return entity.DynamicQuotes{}, fmt.Errorf("DynamicQuoteUseCase - CalcDynamicQuote - uc.webAPI.GetCurrencies: %w", err)
	}

	for i, currency := range currencies {
		currencies[i], err = uc.webAPI.GetCurrencyDynamicQuote(currency)
		currencies[i].CountMin()
		currencies[i].CountMax()
		currencies[i].CountAvg()
		if err != nil {
			return entity.DynamicQuotes{}, fmt.Errorf("DynamicQuoteUseCase - CalcDynamicQuote - uc.webAPI.GetCurrencyDynamicQuote: %w", err)
		}
	}

	return currencies, nil
}
