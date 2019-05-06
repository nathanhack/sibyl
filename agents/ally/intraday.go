package ally

import (
	"context"
	"database/sql"
	"encoding/json"
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

type allyJsonIntradayResponseSingle struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quote jsonIntraday `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

func createAllyJsonIntradayResponse(response string) (*allyJsonIntradayResponse, error) {
	var allyJsonResponse allyJsonIntradayResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return nil, fmt.Errorf("createAllyJsonIntradayResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return nil, fmt.Errorf("createAllyJsonIntradayResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func createAllyJsonIntradayResponseSingle(response string) (*allyJsonIntradayResponse, error) {
	var allyJsonResponse allyJsonIntradayResponseSingle
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return nil, fmt.Errorf("createAllyJsonIntradayResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return nil, fmt.Errorf("createAllyJsonIntradayResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}

	var v allyJsonIntradayResponse
	v.Response.ID = allyJsonResponse.Response.ID
	v.Response.Elapsedtime = allyJsonResponse.Response.Elapsedtime
	v.Response.Quotes.Quote = []jsonIntraday{allyJsonResponse.Response.Quotes.Quote}
	v.Response.Error = allyJsonResponse.Response.Error
	return &v, nil

}

func (ag *AllyAgent) GetIntraday(ctx context.Context, symbol core.StockSymbolType, interval core.IntradayInterval, startDate, endDate core.TimestampType) ([]*core.SibylIntradayRecord, error) {
	//this only works for Stocks
	startTime := time.Now()

	allyHistoryUrl, _ := url.ParseRequestURI("https://api.tradeking.com/v1/market/timesales.json")
	data := allyHistoryUrl.Query()
	data.Add("symbols", string(symbol))

	// can be  "5min", "1min", "tick" (5min is the Ally default)
	data.Add("interval", string(interval))
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

	ag.rateLimitMarketLowPriority.Take(ctx) // this is a lower priority
	ag.rateLimitMarketCalls.Take(ctx)       // and it's a market call
	ag.concurrentLimit.Take(ctx)            // and we limit concurrent requests
	resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
	ag.concurrentLimit.Return()
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
		//if we failed we'll try an run it against the singleton version
		var err1 error
		if response, err1 = createAllyJsonIntradayResponseSingle(string(body)); err1 != nil {
			return []*core.SibylIntradayRecord{}, fmt.Errorf("GetIntraday: %v", err)
		}
	}

	errStrings := []string{}
	toReturn := []*core.SibylIntradayRecord{}
	for _, quote := range response.Response.Quotes.Quote {
		if sq, err := allyJsonIntradayToSibylIntraday(symbol, interval, quote); err != nil {
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

	logrus.Debugf("GetIntraday: finished getting %v Intraday history (%v - %v) for %v in %s", interval, startDate.Time().Format("2006-01-02"), endDate.Time().Format("2006-01-02"), symbol, time.Since(startTime))
	return toReturn, nil
}
func allyJsonIntradayToSibylIntraday(symbol core.StockSymbolType, interval core.IntradayInterval, intraday jsonIntraday) (*core.SibylIntradayRecord, error) {
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

	i := core.IntradayStatusTicks
	switch interval {
	//case core.TickInterval:
	case core.OneMinInterval:
		i = core.IntradayStatus1Min
	case core.FiveMinInterval:
		i = core.IntradayStatus5Min
	}

	return &core.SibylIntradayRecord{
		HighPrice: HighPrice,
		LastPrice: LastPrice,
		Interval:  i,
		LowPrice:  LowPrice,
		OpenPrice: OpenPrice,
		Symbol:    symbol,
		Timestamp: core.NewTimestampTypeFromTime(intraday.DateTimestamp),
		Volume:    Volume,
	}, nil
}
