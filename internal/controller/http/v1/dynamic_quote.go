package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zd4r/cbr_exchange_rate/internal/entity"
	"github.com/zd4r/cbr_exchange_rate/internal/usecase"
)

type dynamicQuoteRoutes struct {
	d usecase.DynamicQuote
	l *zerolog.Logger
}

func newDynamicQuoteRoutes(handler *echo.Group, d usecase.DynamicQuote, l *zerolog.Logger) {
	r := &dynamicQuoteRoutes{d, l}

	handler.GET("/dynamic_quotes", r.dynamicQuotes())
}

type DynamicQuotesResponse struct {
	Quotes []Quote `json:"quotes"`
}

type Quote struct {
	Currency string          `json:"currency"`
	MinQuote entity.MinPrice `json:"min_quote"`
	MaxQuote entity.MaxPrice `json:"max_quote"`
	AvgQuote entity.AvgPrice `json:"avg_quote"`
}

func (dqr *dynamicQuoteRoutes) dynamicQuotes() echo.HandlerFunc {
	return func(c echo.Context) error {
		dqs, err := dqr.d.Calculate(entity.DynamicQuotes{})
		if err != nil {
			dqr.l.Err(err).Msg("http - v1 - dynamicQuotes - dqr.d.Calculate")
			return echo.NewHTTPError(http.StatusInternalServerError, "failed getting dynamic quotes")
		}

		var resp DynamicQuotesResponse
		for _, dq := range dqs {
			if len(dq.Quotes) != 0 {
				var q Quote
				q.Currency = dq.CurrencyName
				q.MinQuote = dq.Min
				q.MaxQuote = dq.Max
				q.AvgQuote = dq.Avg
				resp.Quotes = append(resp.Quotes, q)
			}
		}

		return c.JSON(http.StatusOK, resp)
	}
}
