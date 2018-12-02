package ally

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nathanhack/sibyl/agents/internal/ratelimiter"
	"github.com/nathanhack/sibyl/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	consumerKey      = "ADD_HERE_FOR_TESTING"
	consumerSecret   = "ADD_HERE_FOR_TESTING"
	authToken        = "ADD_HERE_FOR_TESTING"
	authTokenSecrete = "ADD_HERE_FOR_TESTING"

	stockSymbol0                          = "MSFT"
	stockOptionSymbol0                    = "MSFT181130C00115000"
	stockOptionSymbol0StrikePrice         = 115.0
	stockOptionSymbol0ExpirationDateYear  = 2018
	stockOptionSymbol0ExpirationDateMonth = 10
	stockOptionSymbol0ExpirationDateDay   = 11
	stockSymbol1                          = "FB"
	listOfAtLeast30kOptionSymbols         = ""
)

func Test_AllySymbolToSibylSymbol(t *testing.T) {
	symbol, err := toOptionSymbol(stockOptionSymbol0)
	if err != nil {
		t.Errorf("had an error while parsing a valid Option string: %v", err)
		return
	}

	if symbol.Symbol != stockSymbol0 {
		t.Errorf("should have not found a stock %+v", symbol)
		return
	}

	if symbol.OptionType != core.CallEquity {
		t.Errorf("should have found a option %+v", symbol)
		return
	}

	if symbol.StrikePrice != stockOptionSymbol0StrikePrice {
		t.Errorf("should have found a option strike price of 115.00 instead found %v : %+v", symbol.StrikePrice, symbol)
		return
	}
	expectedDate := core.NewDateType(stockOptionSymbol0ExpirationDateYear, stockOptionSymbol0ExpirationDateMonth, stockOptionSymbol0ExpirationDateDay)
	if symbol.Expiration != expectedDate {
		t.Errorf("should have found a option ExpirationTimestamp of %v instead found %v : %+v", expectedDate, symbol.Expiration, symbol)
		return
	}

}

func TestAllyAgent_GetOptionQuotes(t *testing.T) {
	symbolsStrings := strings.Split(listOfAtLeast30kOptionSymbols, ",")
	symbols := map[core.OptionSymbolType]bool{}
	for _, symbol := range symbolsStrings {
		s, _ := toOptionSymbol(symbol)
		symbols[*s] = true
	}

	t.Logf("Total symbols to look for %v", len(symbols))

	//symbol, err := toOptionSymbol(stockOptionSymbol0)
	//if err != nil {
	//	t.Errorf("had an error while parsing a valid Option string: %v", err)
	//	return
	//}

	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	for i := 0; i < 5; i++ {

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
		//_, quotes, err := ac.GetQuotes(ctx, map[core.StockSymbolType]bool{}, map[core.OptionSymbolType]bool{*symbol: true})
		_, quotes, err := ac.GetQuotes(ctx, map[core.StockSymbolType]bool{}, symbols)
		if err != nil {
			t.Errorf("had an error while executing GetAllyQuotes(): %v", err)
			//return
		}
		if len(quotes) != len(symbols) {
			t.Errorf("expoected %v but found %v", len(symbols), len(quotes))
		} else {
			t.Logf("expoected %v but found %v", len(symbols), len(quotes))
		}
		//t.Logf("%+v", quotes)
	}
}

func TestAllyAgent_GetStockQuotes(t *testing.T) {
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	for i := 0; i < 30; i++ {

		quotes, _, err := ac.GetQuotes(ctx, map[core.StockSymbolType]bool{core.StockSymbolType(stockSymbol0): true}, map[core.OptionSymbolType]bool{})
		if err != nil {
			t.Errorf("had an error while executing GetAllyQuotes(): %v", err)
			return
		}

		if len(quotes) != 1 {
			t.Errorf("expected 1 quote but found %v", len(quotes))
		}
	}

	quotes, _, err := ac.GetQuotes(ctx, map[core.StockSymbolType]bool{core.StockSymbolType(stockSymbol0): true, core.StockSymbolType(stockSymbol1): true}, map[core.OptionSymbolType]bool{})
	if err != nil {
		t.Errorf("had an error while executing GetAllyQuotes(): %v", err)
		return
	}

	if len(quotes) != 2 {
		t.Errorf("expected 2 quote but found %v", len(quotes))
	}

}

func TestAllyAgent_GetStableOptionQuote(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	symbol, err := toOptionSymbol(stockOptionSymbol0)

	if err != nil {
		t.Errorf("had an error while parsing a valid Option string: %v", err)
		return
	}

	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	_, quotes, err := ac.GetStableQuotes(ctx, map[core.StockSymbolType]bool{}, map[core.OptionSymbolType]bool{*symbol: true})
	if err != nil {
		t.Errorf("had an error while executing GetStableQuote(): %v", err)
		return
	}

	if len(quotes) != 1 {
		t.Errorf("expected only 1 result but found %v", len(quotes))
	}

	if quotes[0].Symbol != stockSymbol0 {
		t.Errorf("expected %v as the underlying stock but found %v", stockSymbol0, quotes[0].Symbol)
	}

}

func TestAllyAgent_GetStableStockQuote(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	quotes, _, err := ac.GetStableQuotes(ctx, map[core.StockSymbolType]bool{core.StockSymbolType(stockSymbol0): true}, map[core.OptionSymbolType]bool{})
	if err != nil {
		t.Errorf("had an error while executing GetStableQuote():%v", err)
		return
	}

	if len(quotes) != 1 {
		t.Errorf("expected 1 quote but found %v", len(quotes))
	}
}

func TestAllyAgent_GetHistory(t *testing.T) {
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	endDate := core.NewDateTypeFromTime(time.Now())
	startDate := endDate.AddDate(-20, 0, 0)
	quotes, err := ac.GetHistory(ctx, core.StockSymbolType(stockSymbol0), core.DailyTicks, startDate, endDate)

	if err != nil {
		t.Errorf("had an error while executing GetHistory():%v", err)
		return
	}

	if len(quotes) != 2 {
		t.Errorf("had an error while executing GetHistory() expected 1 result found %v : %+v", len(quotes), quotes)
		return
	}
	for _, quote := range quotes {
		t.Logf("%+v", *quote)
		t.Logf("%+v", quote)
	}
}

func TestAllyAgent_GetIntraday(t *testing.T) {
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	quotes, err := ac.GetIntraday(ctx, core.StockSymbolType(stockSymbol0), core.MinuteTicks, core.NewTimestampTypeFromTime(time.Now()).AddDate(0, 0, -58), core.NewTimestampTypeFromTime(time.Now()))

	if err != nil {
		t.Errorf("had an error while executing GetStableQuote(): %v", err)
		return
	}

	if len(quotes) == 0 {
		t.Errorf("had an error while executing GetStableQuote() expected results but found %v : %+v", len(quotes), quotes)
		return
	}
}

func TestAllyAgent_GetStockOptionSymbols(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	options, err := ac.GetStockOptionSymbols(ctx, core.StockSymbolType(stockSymbol0))
	if err != nil {
		t.Errorf("had an error while executing GetStableQuote():%v", err)
		return
	}

	if len(options) == 0 {
		t.Errorf("had an error while executing GetStableQuote() expected some results found %v : %+v", len(options), options)
		return
	}
}

func TestAllyAgent_VerifyStockSymbol(t *testing.T) {
	ac := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	good, hasOptions, _, _, _, err := ac.VerifyStockSymbol(ctx, stockSymbol0)
	if err != nil {
		t.Errorf("had an error while executing GetAllyQuotes(): %v", err)
		return
	}
	if !good {
		t.Errorf("expected a good return")
		return
	}
	if !hasOptions {
		t.Errorf("expected the return to have options")
	}
}

func experiment_FindQuoteRate(t *testing.T) {
	ag := NewAllyAgent(consumerKey, consumerSecret, authToken, authTokenSecrete)
	symbolsStrings := strings.Split(listOfAtLeast30kOptionSymbols, ",")
	symbols := map[core.OptionSymbolType]bool{}
	for _, symbol := range symbolsStrings {
		s, _ := toOptionSymbol(symbol)
		symbols[*s] = true
	}

	symbolStringSlice := make([]string, 0, len(symbols))
	for optionSymbol := range symbols {
		symbolStringSlice = append(symbolStringSlice, toAllySymbol(optionSymbol))
	}
	fmt.Println("Number of symbols:", len(symbolStringSlice))

	sizeAndTimes := map[int][]time.Duration{}
	rateL := ratelimiter.New(1)

	for currentSize := 1000; currentSize < 4000; currentSize += 150 {
		for i := 0; i < 20; i++ {
			if _, has := sizeAndTimes[currentSize]; !has {
				sizeAndTimes[currentSize] = []time.Duration{}
			}

			data := url.Values{}
			data.Add("symbols", strings.Join(symbolStringSlice[:currentSize], ","))
			//next specify the fields we're interested in; this will help reduce bandwidth usage
			data.Add("fids", strings.Join(quoteRequestFields, ","))

			request, err := http.NewRequest(http.MethodPost, "https://api.tradeking.com/v1/market/ext/quotes.xml", strings.NewReader(data.Encode()))
			if err != nil {
				t.Errorf("error making the request %v", err)
				return
			}
			startD := time.Now()
			_ = rateL.Take(context.Background())
			//to make sure this doesn't stall forever (as when there isn't an internet connect)
			// we will kill it after a minute it shouldn't happen as Ally uses a 30 sec time.
			ctx, _ := context.WithTimeout(context.Background(), time.Minute)
			resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
			stopD := time.Now()
			if err != nil {
				fmt.Printf("error : %v\n", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("error StatusCode: %v\n", resp.StatusCode)
				continue
			}

			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if _, err := createAllyXMLStockQuoteResponse(string(body)); err != nil {
				fmt.Printf("error : %v\n", err)
				continue
			}
			timePerSymbol := time.Duration(float64(stopD.Sub(startD)) / float64(currentSize))
			sizeAndTimes[currentSize] = append(sizeAndTimes[currentSize], timePerSymbol)

			fmt.Printf("Round %v index :%v : %v\n", i, currentSize, timePerSymbol)
		}
	}

	st := make(map[int]time.Duration)
	ts := make(map[time.Duration]int)
	times := make([]time.Duration, 0)

	for k, v := range sizeAndTimes {
		num := float64(0)
		for _, i := range v {
			num += float64(i)
		}
		vv := time.Duration(num / float64(len(v)))
		st[k] = vv
		ts[vv] = k
		times = append(times, vv)

	}

	fmt.Println(len(ts), len(st))

	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})

	spew.Dump(st)

	for _, k := range times {
		fmt.Println(k, " : ", ts[k])
	}

}
