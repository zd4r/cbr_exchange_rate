package usecase

import "github.com/zd4r/cbr_exchange_rate/internal/entity"

type (
	DynamicQuoteWebAPI interface {
		GetCurrencies(entity.DynamicQuotes) (entity.DynamicQuotes, error)
		GetCurrencyDynamicQuote(entity.DynamicQuote) (entity.DynamicQuote, error)
	}

	DynamicQuote interface {
		Calculate(entity.DynamicQuotes) (entity.DynamicQuotes, error)
	}
)
