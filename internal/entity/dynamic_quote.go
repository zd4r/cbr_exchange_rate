package entity

import "math"

type DynamicQuotes []DynamicQuote

type DynamicQuote struct {
	ID           string
	CurrencyName string
	Min          MinPrice
	Max          MaxPrice
	Avg          AvgPrice
	Quotes       []Quote
}

type MinPrice struct {
	Value float64 `json:"value"`
	Date  string  `json:"date"`
}

type MaxPrice struct {
	Value float64 `json:"value"`
	Date  string  `json:"date"`
}

type AvgPrice struct {
	Value float64 `json:"value"`
}

type Quote struct {
	Price float64
	Time  string
}

func (dq *DynamicQuote) CountMin() {
	dq.Min.Value = math.MaxFloat64
	for _, v := range dq.Quotes {
		if v.Price < dq.Min.Value {
			dq.Min.Value = roundFloat(v.Price, 4)
			dq.Min.Date = v.Time
		}
	}
}

func (dq *DynamicQuote) CountMax() {
	dq.Max.Value = -math.MaxFloat64
	for _, v := range dq.Quotes {
		if v.Price > dq.Max.Value {
			dq.Max.Value = roundFloat(v.Price, 4)
			dq.Max.Date = v.Time
		}
	}
}

func (dq *DynamicQuote) CountAvg() {
	dq.Avg.Value = 0
	for _, v := range dq.Quotes {
		dq.Avg.Value += v.Price
	}
	dq.Avg.Value = roundFloat(dq.Avg.Value/float64(len(dq.Quotes)), 4)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
