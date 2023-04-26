package webapi

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zd4r/cbr_exchange_rate/internal/entity"
	"golang.org/x/text/encoding/charmap"
)

const DateRange = 90

type CbrWebAPI struct {
	httpClient *http.Client
}

func New() *CbrWebAPI {
	return &CbrWebAPI{
		httpClient: &http.Client{},
	}
}

type CurrenciesResp struct {
	XMLName xml.Name `xml:"Valuta"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Items   []struct {
		Text        string `xml:",chardata"`
		ID          string `xml:"ID,attr"`
		Name        string `xml:"Name"`
		EngName     string `xml:"EngName"`
		Nominal     string `xml:"Nominal"`
		ParentCode  string `xml:"ParentCode"`
		ISONumCode  string `xml:"ISO_Num_Code"`
		ISOCharCode string `xml:"ISO_Char_Code"`
	} `xml:"Item"`
}

func (c *CbrWebAPI) GetCurrencies(dqs entity.DynamicQuotes) (entity.DynamicQuotes, error) {
	req, err := http.NewRequest(http.MethodGet, "http://www.cbr.ru/scripts/XML_valFull.asp", nil)
	if err != nil {
		return entity.DynamicQuotes{}, fmt.Errorf("CbrWebAPI - GetCurrencies - http.NewRequest: %w", err)
	}

	req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return entity.DynamicQuotes{}, fmt.Errorf("CbrWebAPI - GetCurrencies - c.httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	var respData CurrenciesResp
	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = d.Decode(&respData)
	if err != nil {
		return entity.DynamicQuotes{}, fmt.Errorf("CbrWebAPI - GetCurrencies - xml.Decode: %w", err)
	}

	for _, item := range respData.Items {
		var dq entity.DynamicQuote
		dq.CurrencyName = item.Name
		dq.ID = item.ID
		dqs = append(dqs, dq)
	}

	return dqs, nil
}

type CurrencyDynamicQuoteResp struct {
	XMLName    xml.Name `xml:"ValCurs"`
	Text       string   `xml:",chardata"`
	ID         string   `xml:"ID,attr"`
	DateRange1 string   `xml:"DateRange1,attr"`
	DateRange2 string   `xml:"DateRange2,attr"`
	Name       string   `xml:"name,attr"`
	Records    []struct {
		Text    string `xml:",chardata"`
		Date    string `xml:"Date,attr"`
		ID      string `xml:"Id,attr"`
		Nominal string `xml:"Nominal"`
		Value   string `xml:"Value"`
	} `xml:"Record"`
}

func (c *CbrWebAPI) GetCurrencyDynamicQuote(dq entity.DynamicQuote) (entity.DynamicQuote, error) {
	dateReq1, dateReq2 := dateRange(DateRange)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.cbr.ru/scripts/XML_dynamic.asp?date_req1=%s&date_req2=%s&VAL_NM_RQ=%s", dateReq1, dateReq2, dq.ID), nil)
	if err != nil {
		return entity.DynamicQuote{}, fmt.Errorf("CbrWebAPI - GetCurrencyDynamicQuote - http.NewRequest: %w", err)
	}

	req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return entity.DynamicQuote{}, fmt.Errorf("CbrWebAPI - GetCurrencyDynamicQuote - c.httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	var respData CurrencyDynamicQuoteResp
	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = d.Decode(&respData)
	if err != nil {
		return entity.DynamicQuote{}, fmt.Errorf("CbrWebAPI - GetCurrencyDynamicQuote - xml.Decode: %w", err)
	}

	for _, record := range respData.Records {
		var q entity.Quote
		v, err := strconv.ParseFloat(strings.Replace(record.Value, ",", ".", -1), 32)
		if err != nil {
			return entity.DynamicQuote{}, fmt.Errorf("CbrWebAPI - GetCurrencyDynamicQuote - strconv.ParseFloat record.Value: %w", err)
		}
		n, err := strconv.ParseFloat(strings.Replace(record.Nominal, ",", ".", -1), 32)
		if err != nil {
			return entity.DynamicQuote{}, fmt.Errorf("CbrWebAPI - GetCurrencyDynamicQuote - strconv.ParseFloat record.Nominal: %w", err)
		}
		q.Price = v / n
		q.Time = record.Date
		dq.Quotes = append(dq.Quotes, q)
	}

	return dq, nil
}

func dateRange(days int) (string, string) {
	dateReq2 := time.Now()
	dateReq1 := dateReq2.AddDate(0, 0, -days)
	return dateReq1.Format("02/01/2006"), dateReq2.Format("02/01/2006")
}
