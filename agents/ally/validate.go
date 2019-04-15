package ally

import (
	"context"
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

func (ag *AllyAgent) VerifyStockSymbol(ctx context.Context, symbol core.StockSymbolType) (good, hasOptions bool, exchange, exchangeName, name string, err error) {
	//set defaults for returns
	good, hasOptions, exchange, exchangeName, name, err = false, false, "", "", "", nil
	if result, err := ag.verifyGetQuote(ctx, symbol); err != nil {
		err = fmt.Errorf("VerifyStockSymbol: error while checking symbol: %v : %v", symbol, err)
	} else {
		if result.Symbol != symbol {
			err = fmt.Errorf("did not recieve the expect results, expected stock symbol %v but found %v", symbol, result.Symbol)
		} else if result.Timestamp == core.NewTimestampTypeFromUnix(0) {
			err = fmt.Errorf("INVALID STOCK: expected a valid timestamp for the result for stock %v", symbol)
		} else {
			good, hasOptions, exchange, exchangeName, name = true, result.HasOptions, result.Exchange, result.ExchangeDescription, result.Name
		}
	}

	return
}

type agentVerify struct {
	Exchange, ExchangeDescription, Name string
	Symbol                              core.StockSymbolType
	Timestamp                           core.TimestampType
	HasOptions                          bool
}

type jsonVerifyQuote struct {
	Exchange            string `json:"exch"`
	ExchangeDescription string `json:"exch_desc"`
	Name                string `json:"name"`
	OptionFlag          string `json:"op_flag"`
	Symbol              string `json:"symbol"`
	Timestamp           string `json:"timestamp"`
}

var verifyQuoteRequestFields = []string{
	"exch",      //Exchange
	"exch_desc", //ExchangeDescription
	"name",      //Name
	"op_flag",   //OptionFlag : "1" for yes "0" for no
	"symbol",    //Symbol
	"timestamp", //Timestamp
}

type allyJsonVerifyQuoteResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quotetype string          `json:"quotetype"`
			Quote     jsonVerifyQuote `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

func createAllyJsonVerifyQuoteResponse(response string) (*allyJsonVerifyQuoteResponse, error) {
	var allyJsonResponse allyJsonVerifyQuoteResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonVerifyQuoteResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonVerifyQuoteResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func allyJsonQuoteToAgentVerify(quote jsonVerifyQuote) (*agentVerify, error) {

	timestamp, err := strconv.ParseInt(quote.Timestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the timestamp for %v was %v and could not be converted", quote.Symbol, quote.Timestamp)
	}
	/////////////////////////
	exchange := quote.Exchange
	if exchange == "na" || exchange == "" {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the exchange was empty for %v and is expected to be non-empty", quote.Symbol)
	}
	/////////////////////////
	exchangeDescription := quote.ExchangeDescription
	if exchangeDescription == "na" || exchangeDescription == "" {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the exchange description was empty for %v and is expected to be non-empty", quote.Symbol)
	}
	/////////////////////////
	name := quote.Name
	if name == "na" || name == "" {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the name was empty for %v and is expected to be non-empty", quote.Symbol)
	}
	/////////////////////////

	hasOptions := false
	if quote.OptionFlag != "1" && quote.OptionFlag != "0" {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the options flag was empty for %v and is expected to be non-empty", quote.Symbol)
	}
	if quote.OptionFlag == "1" {
		hasOptions = true
	}

	return &agentVerify{
		Exchange:            exchange,
		ExchangeDescription: exchangeDescription,
		HasOptions:          hasOptions,
		Name:                name,
		Symbol:              core.StockSymbolType(quote.Symbol),
		Timestamp:           core.NewTimestampTypeFromUnix(timestamp),
	}, nil
}

func (ag *AllyAgent) verifyGetQuote(ctx context.Context, symbol core.StockSymbolType) (*agentVerify, error) {
	startTime := time.Now()
	data := url.Values{}
	//add the symbols to query

	symbolStringSlice := []string{string(symbol)}

	// add the symbols to the query
	data.Add("symbols", strings.Join(symbolStringSlice, ","))
	//next specify the fields we're interested in; this will help reduce bandwidth usage
	data.Add("fids", strings.Join(verifyQuoteRequestFields, ","))

	request, err := http.NewRequest(http.MethodPost, "https://api.tradeking.com/v1/market/ext/quotes.json", strings.NewReader(data.Encode()))
	if logrus.GetLevel() == logrus.DebugLevel {
		if dump, err := httputil.DumpRequest(request, true); err != nil {
			logrus.Debugf("verifyGetQuote: the request:%v", string(dump))
		} else {
			logrus.Errorf("verifyGetQuote: there was a problem with dumping the request: %v", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("verifyGetQuote: request creation error: %v", err)
	}

	ag.rateLimitMarketLowPriority.Take(ctx) // this is a lower priority
	ag.rateLimitMarketCalls.Take(ctx)       // and it's a market call
	ag.concurrentLimit.Take(ctx)            // and we limit concurrent requests
	resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
	ag.concurrentLimit.Return()
	if err != nil {
		return nil, fmt.Errorf("verifyGetQuote: client error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verifyGetQuote: client error with status code %v: %v", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.Debugf("verifyGetQuote: response body: %v", string(body))
	}

	response, err := createAllyJsonVerifyQuoteResponse(string(body))
	if err != nil {
		return nil, fmt.Errorf("verifyGetQuote: %v", err)
	}

	errStrings := []string{}
	toReturn := []*agentVerify{}

	if sq, err := allyJsonQuoteToAgentVerify(response.Response.Quotes.Quote); err != nil {
		return nil, fmt.Errorf("verifyGetQuote: had errors while parsing quotes: %v", strings.Join(errStrings, ";"))
	} else {
		//see note above where we populate the symbols to send to get quotes
		// only want the first one
		logrus.Debugf("verifyGetQuote: finished found %v quotes in %s", len(toReturn), time.Since(startTime))
		return sq, nil
	}
}
