package ally

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nathanhack/sibyl/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

type jsonData struct {
	ClosePrice string `json:"close"`
	HighPrice  string `json:"high"`
	LowPrice   string `json:"low"`
	OpenPrice  string `json:"open"`
	Timestamp  string `json:"date"`
	Volume     string `json:"volume"`
}

type allyJsonHistoryResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Timeseries  struct {
			Symbol    string `json:"symbol"`
			Startdate string `json:"startdate"`
			Enddate   string `json:"enddate"`
			Series    struct {
				Data []jsonData `json:"data"`
			} `json:"series"`
		} `json:"timeseries"`
		Error string `json:"error"`
	} `json:"response"`
}

type allyJsonHistoryResponseSingle struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Timeseries  struct {
			Symbol    string `json:"symbol"`
			Startdate string `json:"startdate"`
			Enddate   string `json:"enddate"`
			Series    struct {
				Data jsonData `json:"data"`
			} `json:"series"`
		} `json:"timeseries"`
		Error string `json:"error"`
	} `json:"response"`
}

func createAllyJsonHistoryResponse(response string) (*allyJsonHistoryResponse, error) {
	var allyJsonResponse allyJsonHistoryResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonHistoryResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonHistoryResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func createAllyJsonHistoryResponseSingle(response string) (*allyJsonHistoryResponse, error) {
	var allyJsonResponse allyJsonHistoryResponseSingle
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return nil, fmt.Errorf("createAllyJsonHistoryResponseSingle: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return nil, fmt.Errorf("createAllyJsonHistoryResponseSingle: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	toReturn := &allyJsonHistoryResponse{}
	toReturn.Response.ID = allyJsonResponse.Response.ID
	toReturn.Response.Elapsedtime = allyJsonResponse.Response.Elapsedtime
	toReturn.Response.Timeseries.Symbol = allyJsonResponse.Response.Timeseries.Symbol
	toReturn.Response.Timeseries.Startdate = allyJsonResponse.Response.Timeseries.Startdate
	toReturn.Response.Timeseries.Enddate = allyJsonResponse.Response.Timeseries.Enddate
	toReturn.Response.Timeseries.Series.Data = []jsonData{allyJsonResponse.Response.Timeseries.Series.Data}
	toReturn.Response.Error = allyJsonResponse.Response.Error
	return toReturn, nil
}

func (ag *AllyAgent) GetHistory(ctx context.Context, symbol core.StockSymbolType, interval core.HistoryInterval, startDate, endDate core.DateType) ([]*core.SibylHistoryRecord, error) {
	//NOTE this only works on stocks
	startTime := time.Now()

	allyHistoryUrl, _ := url.ParseRequestURI("https://api.tradeking.com/v1/market/historical/search.json")
	data := allyHistoryUrl.Query()
	data.Add("symbols", string(symbol))
	data.Add("interval", string(interval)) // can be  “daily”, “weekly”, “monthly”, or “yearly”.
	data.Add("startdate", startDate.Time().Format("2006-01-02"))
	data.Add("enddate", endDate.Time().Format("2006-01-02"))
	allyHistoryUrl.RawQuery = data.Encode()

	request, err := http.NewRequest(http.MethodGet, allyHistoryUrl.String(), strings.NewReader(data.Encode()))
	if logrus.GetLevel() == logrus.DebugLevel {
		if dump, err := httputil.DumpRequest(request, true); err != nil {
			logrus.Errorf("GetHistory: there was a problem with dumping the request for %v: %v", symbol, err)
		} else {
			logrus.Debugf("GetHistory: the request: %v", string(dump))
		}
	}

	if err != nil {
		return []*core.SibylHistoryRecord{}, fmt.Errorf("GetHistory: request creation error on %v: %v", symbol, err)
	}

	ag.rateLimitMarketLowPriority.Take(ctx) // this is a lower priority
	ag.rateLimitMarketCalls.Take(ctx)       // and it's a market call
	ag.concurrentLimit.Take(ctx)            // and we limit concurrent requests
	resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
	ag.concurrentLimit.Return()
	if err != nil {
		return []*core.SibylHistoryRecord{}, fmt.Errorf("GetHistory: client error for %v: %v", symbol, err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*core.SibylHistoryRecord{}, fmt.Errorf("GetHistory: client error for %v with status code %v: %v", symbol, resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	logrus.Debugf("GetHistory: response body: %v", string(body))

	response, err := createAllyJsonHistoryResponse(string(body))
	if err != nil {
		//if we failed we'll try an run it against the singleton version
		var err1 error
		if response, err1 = createAllyJsonHistoryResponseSingle(string(body)); err1 != nil {
			return []*core.SibylHistoryRecord{}, fmt.Errorf("GetHistory: %v", err)
		}
	}

	errStrings := []string{}
	toReturn := []*core.SibylHistoryRecord{}
	for _, quote := range response.Response.Timeseries.Series.Data {
		if sq, err := allyJsonHistoryToSibylHistory(symbol, interval, quote); err != nil {
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
		return toReturn, fmt.Errorf("GetHistory: had errors while parsing quotes for %v: %v", symbol, strings.Join(errStrings, ";"))
	}

	logrus.Debugf("GetHistory: finished getting history for %v in %s", symbol, time.Since(startTime))
	return toReturn, nil
}

func allyJsonHistoryToSibylHistory(symbol core.StockSymbolType, interval core.HistoryInterval, quote jsonData) (*core.SibylHistoryRecord, error) {
	var err error

	///////////////////////
	var Timestamp time.Time
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, fmt.Errorf("allyJsonHistoryToSibylHistory: failed to create location for %v: %v", symbol, err)
	}
	if Timestamp, err = time.ParseInLocation("2006-01-02", quote.Timestamp, location); err != nil {
		//it must have a timestamp otherwise it breaks down
		return nil, fmt.Errorf("allyJsonHistoryToSibylHistory: json must have a valid timestamp for %v but found %v", symbol, quote.Timestamp)
	}
	/////////////////////////
	var ClosePrice sql.NullFloat64
	if ClosePrice.Float64, err = strconv.ParseFloat(quote.ClosePrice, 64); err == nil {
		ClosePrice.Valid = true
	}
	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(quote.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(quote.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	/////////////////////////
	var openPrice sql.NullFloat64
	if openPrice.Float64, err = strconv.ParseFloat(quote.OpenPrice, 64); err == nil {
		openPrice.Valid = true
	}
	///////////////////////
	var Volume sql.NullInt64
	if Volume.Int64, err = strconv.ParseInt(quote.Volume, 10, 64); err == nil {
		Volume.Valid = true
	} else if tmp, err := strconv.ParseFloat(quote.Volume, 64); err == nil {
		Volume.Int64 = int64(tmp)
		Volume.Valid = true
	}

	i := core.HistoryStatusDaily
	switch interval {
	case core.WeeklyInterval:
		i = core.HistoryStatusWeekly
	case core.MonthlyInterval:
		i = core.HistoryStatusMonthly
	case core.YearlyInterval:
		i = core.HistoryStatusYearly
	}

	return &core.SibylHistoryRecord{
		ClosePrice: ClosePrice,
		HighPrice:  HighPrice,
		Interval:   i,
		LowPrice:   LowPrice,
		OpenPrice:  openPrice,
		Symbol:     symbol,
		Timestamp:  core.NewDateTypeFromTime(Timestamp),
		Volume:     Volume,
	}, nil

}
