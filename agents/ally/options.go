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
	"strings"
	"time"
)

type allyJsonOptionResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quote []struct {
				Symbol string `json:"symbol"`
			} `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

func createAllyJsonOptionResponse(response string) (*allyJsonOptionResponse, error) {
	var allyJsonResponse allyJsonOptionResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonOptionResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonOptionResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func (ag *AllyAgent) GetStockOptionSymbols(ctx context.Context, symbol core.StockSymbolType) ([]*core.OptionSymbolType, error) {
	startTime := time.Now()
	queryNow := fmt.Sprintf("xdat-gte:%v", time.Now().Format("20060102"))

	data := url.Values{}
	data.Add("symbol", string(symbol))
	data.Add("query", queryNow) //just give us stuff from today forward
	data.Add("fids", "symbol")  // and don't waste bandwidth, just send the symbol

	request, err := http.NewRequest(http.MethodPost, "https://api.tradeking.com/v1/market/options/search.json", strings.NewReader(data.Encode()))

	if logrus.GetLevel() == logrus.DebugLevel {
		if dump, err := httputil.DumpRequest(request, true); err != nil {
			logrus.Errorf("GetStockOptionSymbols: there was a problem with dumping the request: %v", err)
		} else {
			logrus.Debugf("GetStockOptionSymbols: the request:%v", string(dump))
		}
	}

	if err != nil {
		return []*core.OptionSymbolType{}, fmt.Errorf("GetStockOptionSymbols: request creation error: %v", err)
	}

	ag.rateLimitMarketLowPriority.Take(ctx) // this is a lower priority
	ag.rateLimitMarketCalls.Take(ctx)       // and it's a market call
	ag.concurrentLimit.Take(ctx)            // and we limit concurrent requests
	resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
	ag.concurrentLimit.Return()
	if err != nil {
		return []*core.OptionSymbolType{}, fmt.Errorf("GetStockOptionSymbols: client error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*core.OptionSymbolType{}, fmt.Errorf("GetStockOptionSymbols: client error for %v with status code %v: %v", symbol, resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.Debugf("GetStockOptionSymbols: response body: %v", string(body))
	}

	response, err := createAllyJsonOptionResponse(string(body))
	if err != nil {
		return []*core.OptionSymbolType{}, fmt.Errorf("GetStockOptionSymbols: %v", err)
	}

	errStrings := []string{}
	toReturn := make([]*core.OptionSymbolType, 0, len(response.Response.Quotes.Quote))
	for _, quote := range response.Response.Quotes.Quote {
		if optionSymbol, err := toOptionSymbol(quote.Symbol); err != nil {
			errStrings = append(errStrings, err.Error())
		} else {
			if optionSymbol.Symbol != symbol {
				//if the feed is sending a bunch of bad data this can explode with errors
				// so we limit it to 10
				if len(errStrings) < 10 {
					errStrings = append(errStrings, fmt.Sprintf("bad options %v", quote.Symbol))
				}
			} else {
				toReturn = append(toReturn, optionSymbol)
			}
		}
	}

	if len(errStrings) != 0 {
		return toReturn, fmt.Errorf("GetStockOptionSymbols: had errors while parsing quotes: %v", strings.Join(errStrings, ";"))
	}

	logrus.Debugf("GetStockOptionSymbols: finished found %v options for stock %v in %s", len(toReturn), symbol, time.Since(startTime))
	return toReturn, nil
}
