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
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// for Ally we have a unified query, for both stock and options
// the response takes care of both, so we sort out which one is a stock or option
// once we get the results from the server

type jsonStableQuote struct {
	AnnualDividend    string `json:"iad"`
	BookValue         string `json:"prbook"`
	ClosePrice        string `json:"cl"`
	ContractSize      string `json:"contract_size"`
	Div               string `json:"div"`
	DivExDate         string `json:"divexdate"`
	DivFreq           string `json:"divfreq"`
	DivPayDate        string `json:"divpaydt"`
	Eps               string `json:"eps"`
	HighPrice52Wk     string `json:"wk52hi"`
	HighPrice52WkDate string `json:"wk52hidate"`
	LowPrice52Wk      string `json:"wk52lo"`
	LowPrice52WkDate  string `json:"wk52lodate"`
	Multiplier        string `json:"prem_mult"`
	OpenPrice         string `json:"opn"`
	PriceEarnings     string `json:"pe"`
	SharesOutstanding string `json:"sho"`
	Symbol            string `json:"symbol"`
	Timestamp         string `json:"timestamp"`
	Volatility        string `json:"volatility12"`
	Yield             string `json:"yield"`
}

var stableQuoteRequestFields = []string{
	"iad",           //AnnualDividend
	"prbook",        //BookValue
	"cl",            //ClosePrice
	"contract_size", //ContractSize
	"div",           //Div
	"divexdate",     //DivExDate
	"divfreq",       //DivFreq
	"divpaydt",      //DivPayDate
	"eps",           //Eps
	"wk52hi",        //HighPrice52Wk
	"wk52hidate",    //HighPrice52WkDate
	"wk52lo",        //LowPrice52Wk
	"wk52lodate",    //LowPrice52WkDate
	"prem_mult",     //Multiplier
	"opn",           //OpenPrice
	"pe",            //PriceEarnings
	"sho",           //SharesOutstanding
	"symbol",        //Symbol
	"timestamp",     //Timestamp
	"volatility12",  //Volatility
	"yield",         //Yield
}

type allyJsonStableQuoteResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quotetype string            `json:"quotetype"`
			Quote     []jsonStableQuote `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

type xmlStableQuote struct {
	AnnualDividend    string `xml:"iad"`
	BookValue         string `xml:"prbook"`
	ClosePrice        string `xml:"cl"`
	ContractSize      string `xml:"contract_size"`
	Div               string `xml:"div"`
	DivExDate         string `xml:"divexdate"`
	DivFreq           string `xml:"divfreq"`
	DivPayDate        string `xml:"divpaydt"`
	Eps               string `xml:"eps"`
	HighPrice52Wk     string `xml:"wk52hi"`
	HighPrice52WkDate string `xml:"wk52hidate"`
	LowPrice52Wk      string `xml:"wk52lo"`
	LowPrice52WkDate  string `xml:"wk52lodate"`
	Multiplier        string `xml:"prem_mult"`
	OpenPrice         string `xml:"opn"`
	PriceEarnings     string `xml:"pe"`
	SharesOutstanding string `xml:"sho"`
	Symbol            string `xml:"symbol"`
	Timestamp         string `xml:"timestamp"`
	Volatility        string `xml:"volatility12"`
	Yield             string `xml:"yield"`
}

type allyXMLStableQuoteResponse struct {
	XMLName     xml.Name `xml:"response"`
	ID          string   `xml:"id,attr"`
	Elapsedtime string   `xml:"elapsedtime"`
	Quotes      struct {
		Quotetype string           `xml:"quotetype"`
		Quote     []xmlStableQuote `xml:"quote"`
	} `xml:"quotes"`
	Error string `xml:"error"`
}

func createAllyJsonStableQuoteResponse(response string) (*allyJsonStableQuoteResponse, error) {
	var allyJsonResponse allyJsonStableQuoteResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonStableQuoteResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonStableQuoteResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func createAllyXMLStockStableQuoteResponse(response string) (*allyXMLStableQuoteResponse, error) {
	var allyXMLResponse allyXMLStableQuoteResponse
	if err := xml.Unmarshal([]byte(response), &allyXMLResponse); err != nil {
		return &allyXMLResponse, fmt.Errorf("createAllyXMLStockStableQuoteResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyXMLResponse.Error != "Success" {
		return &allyXMLResponse, fmt.Errorf("createAllyXMLStockStableQuoteResponse: error: %v : %v", allyXMLResponse.Error, string(response))
	}
	return &allyXMLResponse, nil
}

///////////////////////////

func allyJsonStableQuoteToSibylStableStockQuoteRecord(quote jsonStableQuote) (*core.SibylStableStockQuoteRecord, error) {
	var err error

	///////////////////////
	var timestamp int64
	if timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the time stamp for %v was %v and could not convert", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var annualDividend sql.NullFloat64
	if annualDividend.Float64, err = strconv.ParseFloat(quote.AnnualDividend, 64); err == nil {
		annualDividend.Valid = true
	}
	/////////////////////////
	var bookValue sql.NullFloat64
	if bookValue.Float64, err = strconv.ParseFloat(quote.BookValue, 64); err == nil {
		bookValue.Valid = true
	}
	/////////////////////////
	var closePrice sql.NullFloat64
	if closePrice.Float64, err = strconv.ParseFloat(quote.ClosePrice, 64); err == nil {
		closePrice.Valid = true
	}
	/////////////////////////
	var div sql.NullFloat64
	if div.Float64, err = strconv.ParseFloat(quote.Div, 64); err == nil {
		div.Valid = true
	}
	/////////////////////////
	var divFreq core.NullDivFrequency
	divFreq.DivFreq = core.DivFrequencyType(strings.ToUpper(quote.DivFreq))
	if divFreq.DivFreq == core.AnnualDiv ||
		divFreq.DivFreq == core.SemiAnnualDiv ||
		divFreq.DivFreq == core.QuarterlyDiv ||
		divFreq.DivFreq == core.MonthlyDiv ||
		divFreq.DivFreq == core.NoDiv {
		divFreq.Valid = true
	}
	/////////////////////////
	var divExTimestamp sql.NullInt64
	//format ex. 20180820
	if t, err := time.Parse("20060102", quote.DivExDate); err == nil {
		divExTimestamp.Int64 = t.Unix() //number of seconds since epoch
		divExTimestamp.Valid = true
	}
	/////////////////////////
	var divPayTimestamp sql.NullInt64
	//format ex. 20180820
	if t, err := time.Parse("20060102", quote.DivPayDate); err == nil {
		divPayTimestamp.Int64 = t.Unix() //number of seconds since epoch
		divPayTimestamp.Valid = true
	}
	/////////////////////////
	var eps sql.NullFloat64
	if eps.Float64, err = strconv.ParseFloat(quote.Eps, 64); err == nil {
		eps.Valid = true
	}
	/////////////////////////
	var highPrice52Wk sql.NullFloat64
	if highPrice52Wk.Float64, err = strconv.ParseFloat(quote.HighPrice52Wk, 64); err == nil {
		highPrice52Wk.Valid = true
	}
	/////////////////////////
	var highPrice52WkTimestamp sql.NullInt64
	// example : 20180116
	if t, err := time.Parse("20060102", quote.HighPrice52WkDate); err == nil {
		highPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		highPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var lowPrice52Wk sql.NullFloat64
	if lowPrice52Wk.Float64, err = strconv.ParseFloat(quote.LowPrice52Wk, 64); err == nil {
		lowPrice52Wk.Valid = true
	}
	/////////////////////////
	var lowPrice52WkTimestamp sql.NullInt64
	// example : 20170928
	if t, err := time.Parse("20060102", quote.LowPrice52WkDate); err == nil {
		lowPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		lowPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var openPrice sql.NullFloat64
	if openPrice.Float64, err = strconv.ParseFloat(quote.OpenPrice, 64); err == nil {
		openPrice.Valid = true
	}
	/////////////////////////
	var priceEarnings sql.NullFloat64
	if priceEarnings.Float64, err = strconv.ParseFloat(quote.PriceEarnings, 64); err == nil {
		priceEarnings.Valid = true
	}
	/////////////////////////
	var sharesOutstanding sql.NullInt64
	if sharesOutstanding.Int64, err = strconv.ParseInt(quote.SharesOutstanding, 10, 64); err == nil {
		sharesOutstanding.Valid = true
	}
	/////////////////////////
	var volatility sql.NullFloat64
	if volatility.Float64, err = strconv.ParseFloat(quote.Volatility, 64); err == nil {
		volatility.Valid = true
	}
	/////////////////////////
	var yield sql.NullFloat64
	if yield.Float64, err = strconv.ParseFloat(quote.Yield, 64); err == nil {
		yield.Valid = true
	}

	return &core.SibylStableStockQuoteRecord{
		AnnualDividend:         annualDividend,
		BookValue:              bookValue,
		ClosePrice:             closePrice,
		Div:                    div,
		DivFreq:                divFreq,
		DivExTimestamp:         divExTimestamp,
		DivPayTimestamp:        divPayTimestamp,
		Eps:                    eps,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		OpenPrice:              openPrice,
		PriceEarnings:          priceEarnings,
		SharesOutstanding:      sharesOutstanding,
		Symbol:                 core.StockSymbolType(quote.Symbol),
		Timestamp:              core.NewDateTypeFromUnix(timestamp),
		Volatility:             volatility,
		Yield:                  yield,
	}, nil
}

func allyJsonStableQuoteToSibylStableOptionQuoteRecord(quote jsonStableQuote) (*core.SibylStableOptionQuoteRecord, error) {
	var err error

	///////////////////////
	var optionSymbolType *core.OptionSymbolType
	if optionSymbolType, err = toOptionSymbol(quote.Symbol); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableOptionQuote: had an error with the symbol %v error: %v", quote.Symbol, err)
	}
	///////////////////////
	var timestamp int64
	if timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableOptionQuote: had an error processing the timestamp for the symbol %v error: %v", quote.Symbol, err)
	}
	/////////////////////////
	var closePrice sql.NullFloat64
	if closePrice.Float64, err = strconv.ParseFloat(quote.ClosePrice, 64); err == nil {
		closePrice.Valid = true
	}
	/////////////////////////
	var contractSize sql.NullInt64
	if contractSize.Int64, err = strconv.ParseInt(quote.ContractSize, 10, 64); err == nil {
		contractSize.Valid = true
	}
	/////////////////////////
	var highPrice52Wk sql.NullFloat64
	if highPrice52Wk.Float64, err = strconv.ParseFloat(quote.HighPrice52Wk, 64); err == nil {
		highPrice52Wk.Valid = true
	}
	/////////////////////////
	var highPrice52WkTimestamp sql.NullInt64
	// example : 20180116
	if t, err := time.Parse("20060102", quote.HighPrice52WkDate); err == nil {
		highPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		highPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var lowPrice52Wk sql.NullFloat64
	if lowPrice52Wk.Float64, err = strconv.ParseFloat(quote.LowPrice52Wk, 64); err == nil {
		lowPrice52Wk.Valid = true
	}
	/////////////////////////
	var lowPrice52WkTimestamp sql.NullInt64
	// example : 20170928
	if t, err := time.Parse("20060102", quote.LowPrice52WkDate); err == nil {
		lowPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		lowPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var multiplier sql.NullInt64
	if multiplier.Int64, err = strconv.ParseInt(quote.Multiplier, 10, 64); err == nil {
		multiplier.Valid = true
	}
	/////////////////////////
	var openPrice sql.NullFloat64
	if openPrice.Float64, err = strconv.ParseFloat(quote.OpenPrice, 64); err == nil {
		openPrice.Valid = true
	}
	/////////////////////////

	return &core.SibylStableOptionQuoteRecord{
		ClosePrice:             closePrice,
		ContractSize:           contractSize,
		EquityType:             optionSymbolType.OptionType,
		Expiration:             optionSymbolType.Expiration,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		Multiplier:             multiplier,
		OpenPrice:              openPrice,
		Symbol:                 optionSymbolType.Symbol,
		StrikePrice:            optionSymbolType.StrikePrice,
		Timestamp:              core.NewDateTypeFromUnix(timestamp),
	}, nil
}

func allyXMLStableQuoteToSibylStableStockQuoteRecord(quote xmlStableQuote) (*core.SibylStableStockQuoteRecord, error) {
	var err error

	///////////////////////
	var timestamp int64
	if timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableStockQuote: the time stamp for %v was %v and could not convert", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var annualDividend sql.NullFloat64
	if annualDividend.Float64, err = strconv.ParseFloat(quote.AnnualDividend, 64); err == nil {
		annualDividend.Valid = true
	}
	/////////////////////////
	var bookValue sql.NullFloat64
	if bookValue.Float64, err = strconv.ParseFloat(quote.BookValue, 64); err == nil {
		bookValue.Valid = true
	}
	/////////////////////////
	var closePrice sql.NullFloat64
	if closePrice.Float64, err = strconv.ParseFloat(quote.ClosePrice, 64); err == nil {
		closePrice.Valid = true
	}
	/////////////////////////
	var div sql.NullFloat64
	if div.Float64, err = strconv.ParseFloat(quote.Div, 64); err == nil {
		div.Valid = true
	}
	/////////////////////////
	var divFreq core.NullDivFrequency
	divFreq.DivFreq = core.DivFrequencyType(strings.ToUpper(quote.DivFreq))
	if divFreq.DivFreq == core.AnnualDiv ||
		divFreq.DivFreq == core.SemiAnnualDiv ||
		divFreq.DivFreq == core.QuarterlyDiv ||
		divFreq.DivFreq == core.MonthlyDiv ||
		divFreq.DivFreq == core.NoDiv {
		divFreq.Valid = true
	}
	/////////////////////////
	var divExTimestamp sql.NullInt64
	//format ex. 20180820
	if t, err := time.Parse("20060102", quote.DivExDate); err == nil {
		divExTimestamp.Int64 = t.Unix() //number of seconds since epoch
		divExTimestamp.Valid = true
	}
	/////////////////////////
	var divPayTimestamp sql.NullInt64
	//format ex. 20180820
	if t, err := time.Parse("20060102", quote.DivPayDate); err == nil {
		divPayTimestamp.Int64 = t.Unix() //number of seconds since epoch
		divPayTimestamp.Valid = true
	}
	/////////////////////////
	var eps sql.NullFloat64
	if eps.Float64, err = strconv.ParseFloat(quote.Eps, 64); err == nil {
		eps.Valid = true
	}
	/////////////////////////
	var highPrice52Wk sql.NullFloat64
	if highPrice52Wk.Float64, err = strconv.ParseFloat(quote.HighPrice52Wk, 64); err == nil {
		highPrice52Wk.Valid = true
	}
	/////////////////////////
	var highPrice52WkTimestamp sql.NullInt64
	// example : 20180116
	if t, err := time.Parse("20060102", quote.HighPrice52WkDate); err == nil {
		highPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		highPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var lowPrice52Wk sql.NullFloat64
	if lowPrice52Wk.Float64, err = strconv.ParseFloat(quote.LowPrice52Wk, 64); err == nil {
		lowPrice52Wk.Valid = true
	}
	/////////////////////////
	var lowPrice52WkTimestamp sql.NullInt64
	// example : 20170928
	if t, err := time.Parse("20060102", quote.LowPrice52WkDate); err == nil {
		lowPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		lowPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var openPrice sql.NullFloat64
	if openPrice.Float64, err = strconv.ParseFloat(quote.OpenPrice, 64); err == nil {
		openPrice.Valid = true
	}
	/////////////////////////
	var priceEarnings sql.NullFloat64
	if priceEarnings.Float64, err = strconv.ParseFloat(quote.PriceEarnings, 64); err == nil {
		priceEarnings.Valid = true
	}
	/////////////////////////
	var sharesOutstanding sql.NullInt64
	if sharesOutstanding.Int64, err = strconv.ParseInt(quote.SharesOutstanding, 10, 64); err == nil {
		sharesOutstanding.Valid = true
	}
	/////////////////////////
	var volatility sql.NullFloat64
	if volatility.Float64, err = strconv.ParseFloat(quote.Volatility, 64); err == nil {
		volatility.Valid = true
	}
	/////////////////////////
	var yield sql.NullFloat64
	if yield.Float64, err = strconv.ParseFloat(quote.Yield, 64); err == nil {
		yield.Valid = true
	}

	return &core.SibylStableStockQuoteRecord{
		AnnualDividend:         annualDividend,
		BookValue:              bookValue,
		ClosePrice:             closePrice,
		Div:                    div,
		DivFreq:                divFreq,
		DivExTimestamp:         divExTimestamp,
		DivPayTimestamp:        divPayTimestamp,
		Eps:                    eps,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		OpenPrice:              openPrice,
		PriceEarnings:          priceEarnings,
		SharesOutstanding:      sharesOutstanding,
		Symbol:                 core.StockSymbolType(quote.Symbol),
		Timestamp:              core.NewDateTypeFromUnix(timestamp),
		Volatility:             volatility,
		Yield:                  yield,
	}, nil
}

func allyXMLStableQuoteToSibylStableOptionQuoteRecord(quote xmlStableQuote) (*core.SibylStableOptionQuoteRecord, error) {
	var err error

	///////////////////////
	var optionSymbolType *core.OptionSymbolType
	if optionSymbolType, err = toOptionSymbol(quote.Symbol); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableOptionQuote: had an error with the symbol %v error: %v", quote.Symbol, err)
	}
	///////////////////////
	var timestamp int64
	if timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonQuoteToSibylStableOptionQuote: had an error processing the timestamp for the symbol %v error: %v", quote.Symbol, err)
	}
	/////////////////////////
	var closePrice sql.NullFloat64
	if closePrice.Float64, err = strconv.ParseFloat(quote.ClosePrice, 64); err == nil {
		closePrice.Valid = true
	}
	/////////////////////////
	var contractSize sql.NullInt64
	if contractSize.Int64, err = strconv.ParseInt(quote.ContractSize, 10, 64); err == nil {
		contractSize.Valid = true
	}
	/////////////////////////
	var highPrice52Wk sql.NullFloat64
	if highPrice52Wk.Float64, err = strconv.ParseFloat(quote.HighPrice52Wk, 64); err == nil {
		highPrice52Wk.Valid = true
	}
	/////////////////////////
	var highPrice52WkTimestamp sql.NullInt64
	// example : 20180116
	if t, err := time.Parse("20060102", quote.HighPrice52WkDate); err == nil {
		highPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		highPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var lowPrice52Wk sql.NullFloat64
	if lowPrice52Wk.Float64, err = strconv.ParseFloat(quote.LowPrice52Wk, 64); err == nil {
		lowPrice52Wk.Valid = true
	}
	/////////////////////////
	var lowPrice52WkTimestamp sql.NullInt64
	// example : 20170928
	if t, err := time.Parse("20060102", quote.LowPrice52WkDate); err == nil {
		lowPrice52WkTimestamp.Int64 = t.Unix() //number of seconds since epoch
		lowPrice52WkTimestamp.Valid = true
	}
	/////////////////////////
	var multiplier sql.NullInt64
	if multiplier.Int64, err = strconv.ParseInt(quote.Multiplier, 10, 64); err == nil {
		multiplier.Valid = true
	}
	/////////////////////////
	var openPrice sql.NullFloat64
	if openPrice.Float64, err = strconv.ParseFloat(quote.OpenPrice, 64); err == nil {
		openPrice.Valid = true
	}
	/////////////////////////

	return &core.SibylStableOptionQuoteRecord{
		ClosePrice:             closePrice,
		ContractSize:           contractSize,
		EquityType:             optionSymbolType.OptionType,
		Expiration:             optionSymbolType.Expiration,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		Multiplier:             multiplier,
		OpenPrice:              openPrice,
		Symbol:                 optionSymbolType.Symbol,
		StrikePrice:            optionSymbolType.StrikePrice,
		Timestamp:              core.NewDateTypeFromUnix(timestamp),
	}, nil
}

//////////////////////////
func (ag *AllyAgent) GetStableQuotes(ctx context.Context, stockSymbols map[core.StockSymbolType]bool, optionSymbols map[core.OptionSymbolType]bool) ([]*core.SibylStableStockQuoteRecord, []*core.SibylStableOptionQuoteRecord, error) {
	startTime := time.Now()
	emptyStockRecords := make([]*core.SibylStableStockQuoteRecord, 0)
	emptyOptionRecords := make([]*core.SibylStableOptionQuoteRecord, 0)

	if len(stockSymbols) == 0 && len(optionSymbols) == 0 {
		logrus.Debugf("GetStableQuotes: nothing to find, finished in %s", time.Since(startTime))
		return emptyStockRecords, emptyOptionRecords, nil
	}

	//add the symbols to query
	symbolStringSlice := make([]string, 0, len(stockSymbols)+len(optionSymbols))
	for stockSymbol := range stockSymbols {
		symbolStringSlice = append(symbolStringSlice, string(stockSymbol))
	}
	for optionSymbol := range optionSymbols {
		symbolStringSlice = append(symbolStringSlice, toAllySymbol(optionSymbol))
	}

	singleQuote := false
	if len(symbolStringSlice) == 1 {
		//in the case that we only need one quote we add it twice this is a hack
		// because if we only send one then we receive a different structure,
		// most of the time we will be using this for more than one quote so this is just an
		// easier thing to do.. still hacky ..
		//TODO find a better solution for this case
		singleQuote = true
		symbolStringSlice = append(symbolStringSlice, symbolStringSlice[0])
	}

	//Ally doesn't say this anywhere but it seems like after 30 seconds a
	// request will time out.  After manually testing it, the best seems to be around the 8k mark however
	// unlike the normal GetQuotes() we dont' care about speed here so long as we get the job done.  To that
	// end we'll drop it down to 1k so that it doesn't effect GetQuotes() when then are occurring around the same time.
	// we'll want to do a shuffle because we'll be sending them off requests in batches
	// and we want to ensure the we randomize (to help randomize failures - since we'll be staggering the requests)
	maxSymbolCount := 1000
	if len(symbolStringSlice) > maxSymbolCount {
		rand.Shuffle(len(symbolStringSlice), func(i, j int) {
			tmp := symbolStringSlice[i]
			symbolStringSlice[i] = symbolStringSlice[j]
			symbolStringSlice[j] = tmp
		})
	}

	//we make the groups have about the same amount of symbols
	groupSize := int(float64(len(symbolStringSlice)) / math.Ceil(float64(len(symbolStringSlice))/float64(maxSymbolCount)))
	if len(symbolStringSlice)%groupSize == 1 {
		//and we make sure that the last request has more than one symbol in it
		//we make all of the requests a bit over stuffed but the maxSymbolCount is below
		// the real max so should be a problem
		groupSize++
	}

	symbolGroups := make([][]string, 0)
	for i := 0; i < len(symbolStringSlice); i += groupSize {
		endIndex := min(i+groupSize, len(symbolStringSlice))
		symbolGroups = append(symbolGroups, symbolStringSlice[i:endIndex:endIndex])
	}

	//ally seems to only like one of these requests at a time
	// so we'll do them serially... since things are randomized and the
	// law of large number it should matter way too much
	errStrings := []string{}
	toReturnStockRecords := make([]*core.SibylStableStockQuoteRecord, 0, len(stockSymbols))
	toReturnOptionRecords := make([]*core.SibylStableOptionQuoteRecord, 0, len(optionSymbols))

	for _, symbols := range symbolGroups {
		data := url.Values{}
		// add the symbols to the query
		data.Add("symbols", strings.Join(symbols, ","))
		//next specify the fields we're interested in; this will help reduce bandwidth usage
		data.Add("fids", strings.Join(stableQuoteRequestFields, ","))

		request, err := http.NewRequest(http.MethodPost, "https://api.tradeking.com/v1/market/ext/quotes.xml", strings.NewReader(data.Encode()))
		if logrus.GetLevel() == logrus.DebugLevel {
			if dump, err := httputil.DumpRequest(request, true); err != nil {
				logrus.Debugf("GetStableQuotes: the request:%v", string(dump))
			} else {
				logrus.Errorf("GetStableQuotes: there was a problem with dumping the request: %v", err)
			}
		}

		if err != nil {
			return emptyStockRecords, emptyOptionRecords, fmt.Errorf("GetStableQuotes: request creation error: %v", err)
		}

		ag.quoteFlowLimit.Take(ctx)
		_ = ag.rateLimitMarketCalls.Take(ctx) //we rate limit this call
		resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
		ag.quoteFlowLimit.Return()
		if err != nil {
			return emptyStockRecords, emptyOptionRecords, fmt.Errorf("GetStableQuotes: client error: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return emptyStockRecords, emptyOptionRecords, fmt.Errorf("GetStableQuotes: client error with status code %v: %v", resp.StatusCode, resp.Status)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if logrus.GetLevel() == logrus.DebugLevel {
			logrus.Debugf("GetStableQuotes: response body: %v", string(body))
		}

		response, err := createAllyXMLStockStableQuoteResponse(string(body))
		if err != nil {
			return emptyStockRecords, emptyOptionRecords, fmt.Errorf("GetStableQuotes: %v", err)
		}

		//TODO performance this loop is slow consider something faster
		for _, quote := range response.Quotes.Quote {
			select {
			case <-ctx.Done():
				return emptyStockRecords, emptyOptionRecords, fmt.Errorf("context canceled")
			default:
			}

			//so for each quote we could have either a stock or option
			// to check we take the quotes' symbol and see if it's in the the stocks list
			if _, has := stockSymbols[core.StockSymbolType(quote.Symbol)]; has {
				if sq, err := allyXMLStableQuoteToSibylStableStockQuoteRecord(quote); err != nil {
					//if the feed is sending a bunch of bad data this can explode with errors
					// so we limit it to 10
					if len(errStrings) < 10 {
						errStrings = append(errStrings, err.Error())
					}
				} else {
					toReturnStockRecords = append(toReturnStockRecords, sq)
				}

			} else {
				//else it must be an option
				if sq, err := allyXMLStableQuoteToSibylStableOptionQuoteRecord(quote); err != nil {
					//if the feed is sending a bunch of bad data this can explode with errors
					// so we limit it to 10
					if len(errStrings) < 10 {
						errStrings = append(errStrings, err.Error())
					}
				} else {
					toReturnOptionRecords = append(toReturnOptionRecords, sq)
				}
			}
		}
	}

	if singleQuote {
		//see note above where we populate the symbols to send to get quotes
		// only want the first one, so we find which one then use that
		if len(toReturnStockRecords) > 1 {
			toReturnStockRecords = toReturnStockRecords[:1]
		} else if len(toReturnOptionRecords) > 1 {
			toReturnOptionRecords = toReturnOptionRecords[:1]
		}

	}

	var err error = nil
	if len(errStrings) != 0 {
		err = fmt.Errorf("GetStableQuotes: had errors while parsing quotes: %v", strings.Join(errStrings, ";"))
	}

	logrus.Debugf("GetStableQuotes: finished found %v stock quotes and %v option quotes in %s", len(toReturnStockRecords), len(toReturnOptionRecords), time.Since(startTime))
	return toReturnStockRecords, toReturnOptionRecords, err
}
