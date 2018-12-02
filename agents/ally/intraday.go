package ally

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type jsonIntraday struct {
	Date          string    `json:"date"`     //The date "YYYY-MM-DD" for this intraday data point
	DateTimestamp time.Time `json:"datetime"` //the timestamp "YYYY-MM-DDTHH:MM:SSZ" (ex. 2018-10-10T13:30:00Z)  for this data point
	HighPrice     string    `json:"hi"`       //The high price
	IncrVl        string    `json:"incr_vl"`
	LastPrice     string    `json:"last"`
	LowPrice      string    `json:"lo"`
	OpenPrice     string    `json:"opn"`
	Timestamp     time.Time `json:"timestamp"`
	Volume        string    `json:"vl"`
}

type allyJsonIntradayResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quote []jsonIntraday `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

type xmlIntraday struct {
	Date          string `xml:"date"`
	DateTimestamp string `xml:"datetime"`
	HighPrice     string `xml:"hi"`
	IncrVl        string `xml:"incr_vl"`
	LastPrice     string `xml:"last"`
	LowPrice      string `xml:"lo"`
	OpenPrice     string `xml:"opn"`
	Timestamp     string `xml:"timestamp"`
	Volume        string `xml:"vl"`
}

type allXMLIntradayResponse struct {
	XMLName     xml.Name `xml:"response"`
	Text        string   `xml:",chardata"`
	ID          string   `xml:"id,attr"`
	Elapsedtime string   `xml:"elapsedtime"`
	Quotes      struct {
		Text  string        `xml:",chardata"`
		Quote []xmlIntraday `xml:"quote"`
	} `xml:"quotes"`
	Error string `xml:"error"`
}

func createAllyJsonIntradayResponse(response string) (*allyJsonIntradayResponse, error) {
	var allyJsonResponse allyJsonIntradayResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonIntradayResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonIntradayResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func (ag *AllyAgent) GetIntraday(ctx context.Context, symbol core.StockSymbolType, tickSize core.HistoryTick, startDate, endDate core.TimestampType) ([]*core.SibylIntradayRecord, error) {
	//this only works for Stocks
	startTime := time.Now()
	if tickSize != core.MinuteTicks && tickSize != core.FiveMinuteTicks {
		return nil, fmt.Errorf("GetIntraday: tickSize must be 1 or 5 minutes")
	}

	allyHistoryUrl, _ := url.ParseRequestURI("https://api.tradeking.com/v1/market/timesales.json")
	data := allyHistoryUrl.Query()
	data.Add("symbols", string(symbol))

	if tickSize == core.MinuteTicks {
		data.Add("interval", "1min") // can be  "5min", "1min", "tick" (5min is the default)
	} else if tickSize == core.FiveMinuteTicks {
		data.Add("interval", "5min")
	}

	data.Add("startdate", startDate.Time().Format("2006-01-02"))
	data.Add("enddate", endDate.Time().Format("2006-01-02"))

	allyHistoryUrl.RawQuery = data.Encode()

	request, err := http.NewRequest(http.MethodGet, allyHistoryUrl.String(), strings.NewReader(data.Encode()))
	if logrus.GetLevel() == logrus.DebugLevel {
		if dump, err := httputil.DumpRequest(request, true); err != nil {
			logrus.Errorf("GetIntraday: there was a problem with dumping the request: %v", err)
		} else {
			logrus.Debugf("GetIntraday: the request:%v", string(dump))
		}
	}

	if err != nil {
		return []*core.SibylIntradayRecord{}, fmt.Errorf("GetIntraday: request creation error: %v", err)
	}

	ag.rateLimitMarketLowPriority.Take(ctx)
	ag.rateLimitMarketCalls.Take(ctx)
	resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
	if err != nil {
		return []*core.SibylIntradayRecord{}, fmt.Errorf("GetIntraday: client error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*core.SibylIntradayRecord{}, fmt.Errorf("GetIntraday: client error for %v with status code %v: %v", symbol, resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	logrus.Debugf("GetIntraday: response body: %v", string(body))

	response, err := createAllyJsonIntradayResponse(string(body))
	if err != nil {
		return []*core.SibylIntradayRecord{}, fmt.Errorf("GetIntraday: %v", err)
	}

	errStrings := []string{}
	toReturn := []*core.SibylIntradayRecord{}
	for _, quote := range response.Response.Quotes.Quote {
		if sq, err := allyJsonIntradayToSibylIntraday(symbol, quote); err != nil {
			//if the feed is sending a bunch of bad data this can explode with errors
			// so we limit it to 10
			if len(errStrings) < 10 {
				errStrings = append(errStrings, err.Error())
			}
		} else {
			toReturn = append(toReturn, sq)
		}
	}
	if len(errStrings) != 0 {
		return toReturn, fmt.Errorf("GetIntraday: had errors while parsing quotes: %v", strings.Join(errStrings, ";"))
	}

	logrus.Infof("GetIntraday: finished getting intraday history for %v in %s", symbol, time.Since(startTime))
	return toReturn, nil
}
func allyJsonIntradayToSibylIntraday(symbol core.StockSymbolType, intraday jsonIntraday) (*core.SibylIntradayRecord, error) {
	var err error

	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(intraday.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var LastPrice sql.NullFloat64
	if LastPrice.Float64, err = strconv.ParseFloat(intraday.LastPrice, 64); err == nil {
		LastPrice.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(intraday.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	///////////////////////
	var OpenPrice sql.NullFloat64
	if OpenPrice.Float64, err = strconv.ParseFloat(intraday.OpenPrice, 64); err == nil {
		OpenPrice.Valid = true
	}
	///////////////////////
	var Volume sql.NullInt64
	if Volume.Int64, err = strconv.ParseInt(intraday.Volume, 10, 64); err == nil {
		Volume.Valid = true
	}
	///////////////////////

	return &core.SibylIntradayRecord{
		HighPrice: HighPrice,
		LastPrice: LastPrice,
		LowPrice:  LowPrice,
		OpenPrice: OpenPrice,
		Symbol:    symbol,
		Timestamp: core.NewTimestampTypeFromTime(intraday.DateTimestamp),
		Volume:    Volume,
	}, nil
}
