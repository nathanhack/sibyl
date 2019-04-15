package ally

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nathanhack/sibyl/core"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

// for Ally we have a unified query, for both stock and options
// the response takes care of both, so we sort out which one is a stock or option
// once we get the results from the server

type jsonQuote struct {
	Ask                 string `json:"ask"`
	AskTime             string `json:"ask_time"`
	AskSize             string `json:"asksz"`
	Beta                string `json:"beta"`
	Bid                 string `json:"bid"`
	BidTime             string `json:"bid_time"`
	BidSize             string `json:"bidsz"`
	Change              string `json:"chg"`
	ChangeSign          string `json:"chg_sign"`
	Delta               string `json:"idelta"`
	Gamma               string `json:"igamma"`
	HighPrice           string `json:"hi"`
	ImpliedVolatility   string `json:"imp_volatility"`
	LastTrade           string `json:"last"`
	LastTradeTimestamp  string `json:"datetime"`
	LastTradeVolume     string `json:"incr_vl"`
	LowPrice            string `json:"lo"`
	OpenInterest        string `json:"openinterest"`
	Rho                 string `json:"irho"`
	Symbol              string `json:"symbol"`
	Theta               string `json:"itheta"`
	Timestamp           string `json:"timestamp"`
	Vega                string `json:"ivega"`
	Volume              string `json:"vl"`
	VolWeightedAvgPrice string `json:"vwap"`
}

var quoteRequestFields = []string{
	"ask",      //Ask
	"ask_time", //AskTime
	"asksz",    //AskSize
	"beta",     //Beta
	"bid",      //Bid
	"bid_time", //BidTime
	"bidsz",    //BidSize

	"chg",      //Change
	"chg_sign", //ChangeSign

	"idelta",         //Delta
	"igamma",         //Gamma
	"hi",             //HighPrice
	"imp_volatility", //ImpliedVolatility
	"last",           //LastTrade
	"datetime",       //LastTradeTimestamp
	"incr_vl",        //LastTradeVolume
	"lo",             //LowPrice
	"openinterest",   //OpenInterest
	"irho",           //Rho
	"symbol",         //Symbol
	"itheta",         //Theta
	"timestamp",      //TimeStamp
	"ivega",          //Vega
	"vl",             //Volume
	"vwap",           //VolWeightedAvgPrice
}

type allyJsonQuoteResponse struct {
	Response struct {
		ID          string `json:"@id"`
		Elapsedtime string `json:"elapsedtime"`
		Quotes      struct {
			Quotetype string      `json:"quotetype"`
			Quote     []jsonQuote `json:"quote"`
		} `json:"quotes"`
		Error string `json:"error"`
	} `json:"response"`
}

type xmlQuote struct {
	Ask                 string `xml:"ask"`
	AskTime             string `xml:"ask_time"`
	AskSize             string `xml:"asksz"`
	Beta                string `xml:"beta"`
	Bid                 string `xml:"bid"`
	BidTime             string `xml:"bid_time"`
	BidSize             string `xml:"bidsz"`
	Change              string `xml:"chg"`
	ChangeSign          string `xml:"chg_sign"`
	LastTradeTimestamp  string `xml:"datetime"`
	HighPrice           string `xml:"hi"`
	Delta               string `xml:"idelta"`
	Gamma               string `xml:"igamma"`
	ImpliedVolatility   string `xml:"imp_volatility"`
	LastTradeVolume     string `xml:"incr_vl"`
	Rho                 string `xml:"irho"`
	Theta               string `xml:"itheta"`
	Vega                string `xml:"ivega"`
	LastTrade           string `xml:"last"`
	LowPrice            string `xml:"lo"`
	OpenInterest        string `xml:"openinterest"`
	Symbol              string `xml:"symbol"`
	Timestamp           string `xml:"timestamp"`
	Volume              string `xml:"vl"`
	VolWeightedAvgPrice string `xml:"vwap"`
}

type allyXMLQuoteResponse struct {
	XMLName     xml.Name `xml:"response"`
	ID          string   `xml:"id,attr"`
	Elapsedtime string   `xml:"elapsedtime"`
	Quotes      struct {
		Quotetype string     `xml:"quotetype"`
		Quote     []xmlQuote `xml:"quote"`
	} `xml:"quotes"`
	Error string `xml:"error"`
}

func createAllyJsonStockQuoteResponse(response string) (*allyJsonQuoteResponse, error) {
	var allyJsonResponse allyJsonQuoteResponse
	if err := json.Unmarshal([]byte(response), &allyJsonResponse); err != nil {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonStockQuoteResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyJsonResponse.Response.Error != "Success" {
		return &allyJsonResponse, fmt.Errorf("createAllyJsonStockQuoteResponse: error: %v : %v", allyJsonResponse.Response.Error, string(response))
	}
	return &allyJsonResponse, nil
}

func createAllyXMLStockQuoteResponse(response string) (*allyXMLQuoteResponse, error) {
	var allyXMLResponse allyXMLQuoteResponse
	if err := xml.Unmarshal([]byte(response), &allyXMLResponse); err != nil {
		return &allyXMLResponse, fmt.Errorf("createAllyXMLStockQuoteResponse: unmarshal failure: %v\n on response: %v", err, response)
	}

	if allyXMLResponse.Error != "Success" {
		return &allyXMLResponse, fmt.Errorf("createAllyXMLStockQuoteResponse: error: %v : %v", allyXMLResponse.Error, string(response))
	}
	return &allyXMLResponse, nil
}

///////////////////////////
//below is both json and xml versions, currently it looks like XML is more optimized on the Ally server

func allyJsonQuoteToSibylStockQuoteRecord(quote jsonQuote) (*core.SibylStockQuoteRecord, error) {
	var err error
	timeRegex := regexp.MustCompile(`(?P<hour>\d{2}):(?P<min>\d{2})`)
	var match []string

	///////////////////////
	var Timestamp int64
	if Timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonStockQuoteToSibylQuote: the timestamp for %v was %v and could not convert", quote.Symbol, quote.Timestamp)
	}
	if Timestamp == 0 {
		//if the timestamp was erroneous so bail
		return nil, fmt.Errorf("allyJsonStockQuoteToSibylQuote: the timestamp for %v was erroneous: %v", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var Ask sql.NullFloat64
	if Ask.Float64, err = strconv.ParseFloat(quote.Ask, 64); err == nil {
		Ask.Valid = true
	}
	///////////////////////
	var AskTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.AskTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				AskTime.Int64 = hour*100 + min
				AskTime.Valid = true
			}
		}
	}
	///////////////////////
	var AskSize sql.NullInt64
	if AskSize.Int64, err = strconv.ParseInt(quote.AskSize, 10, 64); err == nil {
		AskSize.Valid = true
	}
	///////////////////////
	var Beta sql.NullFloat64
	if Beta.Float64, err = strconv.ParseFloat(quote.Beta, 64); err == nil {
		Beta.Valid = true
	}
	///////////////////////
	var Bid sql.NullFloat64
	if Bid.Float64, err = strconv.ParseFloat(quote.Bid, 64); err == nil {
		Bid.Valid = true
	}
	///////////////////////
	var BidTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.BidTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				BidTime.Int64 = hour*100 + min
				BidTime.Valid = true
			}
		}
	}
	///////////////////////
	var BidSize sql.NullInt64
	if BidSize.Int64, err = strconv.ParseInt(quote.BidSize, 10, 64); err == nil {
		BidSize.Valid = true
	}
	///////////////////////
	var Change sql.NullFloat64
	if Change.Float64, err = strconv.ParseFloat(quote.Change, 64); err == nil {
		if strings.ToLower(quote.ChangeSign) != "na" {
			Change.Valid = true
			if quote.ChangeSign == "d" {
				Change.Float64 *= -1
			}
		}
	}
	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(quote.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var LastTradePrice sql.NullFloat64
	if LastTradePrice.Float64, err = strconv.ParseFloat(quote.LastTrade, 64); err == nil {
		LastTradePrice.Valid = true
	}
	///////////////////////
	var LastTradeTimestamp sql.NullInt64
	// example: 2018-09-21T00:00:00-04:00
	if t, err := time.Parse("2006-01-02T15:04:05-07:00", quote.LastTradeTimestamp); err == nil {
		LastTradeTimestamp.Int64 = t.Unix() //number of seconds since epoch
		LastTradeTimestamp.Valid = true
	}
	///////////////////////
	var LastTradeVolume sql.NullInt64
	if LastTradeVolume.Int64, err = strconv.ParseInt(quote.LastTradeVolume, 10, 64); err == nil {
		LastTradeVolume.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(quote.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	///////////////////////
	var Volume sql.NullInt64
	if Volume.Int64, err = strconv.ParseInt(quote.Volume, 10, 64); err == nil {
		Volume.Valid = true
	}
	///////////////////////
	var VolWeightedAvgPrice sql.NullFloat64
	if VolWeightedAvgPrice.Float64, err = strconv.ParseFloat(quote.VolWeightedAvgPrice, 64); err == nil {
		VolWeightedAvgPrice.Valid = true
	}

	return &core.SibylStockQuoteRecord{
		Ask:                 Ask,
		AskTime:             AskTime,
		AskSize:             AskSize,
		Beta:                Beta,
		Bid:                 Bid,
		BidTime:             BidTime,
		BidSize:             BidSize,
		Change:              Change,
		HighPrice:           HighPrice,
		LastTradePrice:      LastTradePrice,
		LastTradeTimestamp:  LastTradeTimestamp,
		LastTradeVolume:     LastTradeVolume,
		LowPrice:            LowPrice,
		Symbol:              core.StockSymbolType(quote.Symbol),
		Timestamp:           core.NewTimestampTypeFromUnix(Timestamp),
		Volume:              Volume,
		VolWeightedAvgPrice: VolWeightedAvgPrice,
	}, nil
}

func allyJsonQuoteToSibylOptionQuoteRecord(quote jsonQuote) (*core.SibylOptionQuoteRecord, error) {
	var err error
	timeRegex := regexp.MustCompile(`(?P<hour>\d{2}):(?P<min>\d{2})`)
	var match []string

	///////////////////////
	var optionSymbolType *core.OptionSymbolType
	if optionSymbolType, err = toOptionSymbol(quote.Symbol); err != nil {
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: had an error with the symbol %v error: %v", quote.Symbol, err)
	}
	///////////////////////
	var Timestamp int64
	if Timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: had an error processing the timestamp for the symbol %v error: %v", quote.Symbol, err)
	}
	if Timestamp == 0 {
		//if the timestamp was erroneous so bail
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: the timestamp for %v was erroneous: %v", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var Ask sql.NullFloat64
	if Ask.Float64, err = strconv.ParseFloat(quote.Ask, 64); err == nil {
		Ask.Valid = true
	}
	///////////////////////
	var AskTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.AskTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				AskTime.Int64 = hour*100 + min
				AskTime.Valid = true
			}
		}
	}
	///////////////////////
	var AskSize sql.NullInt64
	if AskSize.Int64, err = strconv.ParseInt(quote.AskSize, 10, 64); err == nil {
		AskSize.Valid = true
	}
	///////////////////////
	var Bid sql.NullFloat64
	if Bid.Float64, err = strconv.ParseFloat(quote.Bid, 64); err == nil {
		Bid.Valid = true
	}
	///////////////////////
	var BidTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.BidTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				BidTime.Int64 = hour*100 + min
				BidTime.Valid = true
			}
		}
	}
	///////////////////////
	var BidSize sql.NullInt64
	if BidSize.Int64, err = strconv.ParseInt(quote.BidSize, 10, 64); err == nil {
		BidSize.Valid = true
	}
	///////////////////////
	var Change sql.NullFloat64
	if Change.Float64, err = strconv.ParseFloat(quote.Change, 64); err == nil {
		if strings.ToLower(quote.ChangeSign) != "na" {
			Change.Valid = true
			if quote.ChangeSign == "d" {
				Change.Float64 *= -1
			}
		}
	}
	///////////////////////
	var Delta sql.NullFloat64
	if Delta.Float64, err = strconv.ParseFloat(quote.Delta, 64); err == nil {
		Delta.Valid = true

	}
	///////////////////////
	var Gamma sql.NullFloat64
	if Gamma.Float64, err = strconv.ParseFloat(quote.Gamma, 64); err == nil {
		Gamma.Valid = true
	}
	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(quote.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var ImpliedVolatility sql.NullFloat64
	if ImpliedVolatility.Float64, err = strconv.ParseFloat(quote.ImpliedVolatility, 64); err == nil {
		ImpliedVolatility.Valid = true
	}
	///////////////////////
	var LastTradePrice sql.NullFloat64
	if LastTradePrice.Float64, err = strconv.ParseFloat(quote.LastTrade, 64); err == nil {
		LastTradePrice.Valid = true
	}
	///////////////////////
	var LastTradeTimestamp sql.NullInt64
	// example: 2018-09-21T00:00:00-04:00
	if t, err := time.Parse("2006-01-02T15:04:05-07:00", quote.LastTradeTimestamp); err == nil {
		LastTradeTimestamp.Int64 = t.Unix() //number of seconds since epoch
		LastTradeTimestamp.Valid = true
	}
	///////////////////////
	var LastTradeVolume sql.NullInt64
	if LastTradeVolume.Int64, err = strconv.ParseInt(quote.LastTradeVolume, 10, 64); err == nil {
		LastTradeVolume.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(quote.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	///////////////////////
	var OpenInterest sql.NullInt64
	if OpenInterest.Int64, err = strconv.ParseInt(quote.OpenInterest, 10, 64); err == nil {
		OpenInterest.Valid = true
	}
	///////////////////////
	var Rho sql.NullFloat64
	if Rho.Float64, err = strconv.ParseFloat(quote.Rho, 64); err == nil {
		Rho.Valid = true
	}
	///////////////////////
	var Theta sql.NullFloat64
	if Theta.Float64, err = strconv.ParseFloat(quote.Theta, 64); err == nil {
		Theta.Valid = true
	}
	///////////////////////
	var Vega sql.NullFloat64
	if Vega.Float64, err = strconv.ParseFloat(quote.Vega, 64); err == nil {
		Vega.Valid = true
	}

	return &core.SibylOptionQuoteRecord{
		Ask:                Ask,
		AskTime:            AskTime,
		AskSize:            AskSize,
		Bid:                Bid,
		BidTime:            BidTime,
		BidSize:            BidSize,
		Change:             Change,
		Delta:              Delta,
		EquityType:         optionSymbolType.OptionType,
		Expiration:         optionSymbolType.Expiration,
		Gamma:              Gamma,
		HighPrice:          HighPrice,
		ImpliedVolatility:  ImpliedVolatility,
		LastTradePrice:     LastTradePrice,
		LastTradeVolume:    LastTradeVolume,
		LastTradeTimestamp: LastTradeTimestamp,
		LowPrice:           LowPrice,
		OpenInterest:       OpenInterest,
		Rho:                Rho,
		Symbol:             optionSymbolType.Symbol,
		StrikePrice:        optionSymbolType.StrikePrice,
		Theta:              Theta,
		Timestamp:          core.NewTimestampTypeFromUnix(Timestamp),
		Vega:               Vega,
	}, nil
}

func allyXMLQuoteToSibylStockQuoteRecord(quote xmlQuote) (*core.SibylStockQuoteRecord, error) {
	var err error
	timeRegex := regexp.MustCompile(`(?P<hour>\d{2}):(?P<min>\d{2})`)
	var match []string

	///////////////////////
	var Timestamp int64
	if Timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonStockQuoteToSibylQuote: the timestamp for %v was %v and could not convert", quote.Symbol, quote.Timestamp)
	}
	if Timestamp == 0 {
		//if the timestamp was erroneous so bail
		return nil, fmt.Errorf("allyJsonStockQuoteToSibylQuote: the timestamp for %v was erroneous: %v", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var Ask sql.NullFloat64
	if Ask.Float64, err = strconv.ParseFloat(quote.Ask, 64); err == nil {
		Ask.Valid = true
	}
	///////////////////////
	var AskTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.AskTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				AskTime.Int64 = hour*100 + min
				AskTime.Valid = true
			}
		}
	}
	///////////////////////
	var AskSize sql.NullInt64
	if AskSize.Int64, err = strconv.ParseInt(quote.AskSize, 10, 64); err == nil {
		AskSize.Valid = true
	}
	///////////////////////
	var Beta sql.NullFloat64
	if Beta.Float64, err = strconv.ParseFloat(quote.Beta, 64); err == nil {
		Beta.Valid = true
	}
	///////////////////////
	var Bid sql.NullFloat64
	if Bid.Float64, err = strconv.ParseFloat(quote.Bid, 64); err == nil {
		Bid.Valid = true
	}
	///////////////////////
	var BidTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.BidTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				BidTime.Int64 = hour*100 + min
				BidTime.Valid = true
			}
		}
	}
	///////////////////////
	var BidSize sql.NullInt64
	if BidSize.Int64, err = strconv.ParseInt(quote.BidSize, 10, 64); err == nil {
		BidSize.Valid = true
	}
	///////////////////////
	var Change sql.NullFloat64
	if Change.Float64, err = strconv.ParseFloat(quote.Change, 64); err == nil {
		if strings.ToLower(quote.ChangeSign) != "na" {
			Change.Valid = true
			if quote.ChangeSign == "d" {
				Change.Float64 *= -1
			}
		}
	}
	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(quote.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var LastTradePrice sql.NullFloat64
	if LastTradePrice.Float64, err = strconv.ParseFloat(quote.LastTrade, 64); err == nil {
		LastTradePrice.Valid = true
	}
	///////////////////////
	var LastTradeTimestamp sql.NullInt64
	// example: 2018-09-21T00:00:00-04:00
	if t, err := time.Parse("2006-01-02T15:04:05-07:00", quote.LastTradeTimestamp); err == nil {
		LastTradeTimestamp.Int64 = t.Unix() //number of seconds since epoch
		LastTradeTimestamp.Valid = true
	}
	///////////////////////
	var LastTradeVolume sql.NullInt64
	if LastTradeVolume.Int64, err = strconv.ParseInt(quote.LastTradeVolume, 10, 64); err == nil {
		LastTradeVolume.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(quote.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	///////////////////////
	var Volume sql.NullInt64
	if Volume.Int64, err = strconv.ParseInt(quote.Volume, 10, 64); err == nil {
		Volume.Valid = true
	}
	///////////////////////
	var VolWeightedAvgPrice sql.NullFloat64
	if VolWeightedAvgPrice.Float64, err = strconv.ParseFloat(quote.VolWeightedAvgPrice, 64); err == nil {
		VolWeightedAvgPrice.Valid = true
	}

	return &core.SibylStockQuoteRecord{
		Ask:                 Ask,
		AskTime:             AskTime,
		AskSize:             AskSize,
		Beta:                Beta,
		Bid:                 Bid,
		BidTime:             BidTime,
		BidSize:             BidSize,
		Change:              Change,
		HighPrice:           HighPrice,
		LastTradePrice:      LastTradePrice,
		LastTradeTimestamp:  LastTradeTimestamp,
		LastTradeVolume:     LastTradeVolume,
		LowPrice:            LowPrice,
		Symbol:              core.StockSymbolType(quote.Symbol),
		Timestamp:           core.NewTimestampTypeFromUnix(Timestamp),
		Volume:              Volume,
		VolWeightedAvgPrice: VolWeightedAvgPrice,
	}, nil
}

func allyXMLQuoteToSibylOptionQuoteRecord(quote xmlQuote) (*core.SibylOptionQuoteRecord, error) {
	var err error
	timeRegex := regexp.MustCompile(`(?P<hour>\d{2}):(?P<min>\d{2})`)
	var match []string

	///////////////////////
	var optionSymbolType *core.OptionSymbolType
	if optionSymbolType, err = toOptionSymbol(quote.Symbol); err != nil {
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: had an error with the symbol %v error: %v", quote.Symbol, err)
	}
	///////////////////////
	var Timestamp int64
	if Timestamp, err = strconv.ParseInt(quote.Timestamp, 10, 64); err != nil {
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: had an error processing the timestamp for the symbol %v error: %v", quote.Symbol, err)
	}
	if Timestamp == 0 {
		//if the timestamp was erroneous so bail
		return nil, fmt.Errorf("allyJsonOptionQuoteToSibylQuote: the timestamp for %v was erroneous: %v", quote.Symbol, quote.Timestamp)
	}
	///////////////////////
	var Ask sql.NullFloat64
	if Ask.Float64, err = strconv.ParseFloat(quote.Ask, 64); err == nil {
		Ask.Valid = true
	}
	///////////////////////
	var AskTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.AskTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				AskTime.Int64 = hour*100 + min
				AskTime.Valid = true
			}
		}
	}
	///////////////////////
	var AskSize sql.NullInt64
	if AskSize.Int64, err = strconv.ParseInt(quote.AskSize, 10, 64); err == nil {
		AskSize.Valid = true
	}
	///////////////////////
	var Bid sql.NullFloat64
	if Bid.Float64, err = strconv.ParseFloat(quote.Bid, 64); err == nil {
		Bid.Valid = true
	}
	///////////////////////
	var BidTime sql.NullInt64
	match = timeRegex.FindStringSubmatch(quote.BidTime)
	if match != nil {
		if hour, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			if min, err := strconv.ParseInt(match[2], 10, 64); err == nil {
				BidTime.Int64 = hour*100 + min
				BidTime.Valid = true
			}
		}
	}
	///////////////////////
	var BidSize sql.NullInt64
	if BidSize.Int64, err = strconv.ParseInt(quote.BidSize, 10, 64); err == nil {
		BidSize.Valid = true
	}
	///////////////////////
	var Change sql.NullFloat64
	if Change.Float64, err = strconv.ParseFloat(quote.Change, 64); err == nil {
		if strings.ToLower(quote.ChangeSign) != "na" {
			Change.Valid = true
			if quote.ChangeSign == "d" {
				Change.Float64 *= -1
			}
		}
	}
	///////////////////////
	var Delta sql.NullFloat64
	if Delta.Float64, err = strconv.ParseFloat(quote.Delta, 64); err == nil {
		Delta.Valid = true

	}
	///////////////////////
	var Gamma sql.NullFloat64
	if Gamma.Float64, err = strconv.ParseFloat(quote.Gamma, 64); err == nil {
		Gamma.Valid = true
	}
	///////////////////////
	var HighPrice sql.NullFloat64
	if HighPrice.Float64, err = strconv.ParseFloat(quote.HighPrice, 64); err == nil {
		HighPrice.Valid = true
	}
	///////////////////////
	var ImpliedVolatility sql.NullFloat64
	if ImpliedVolatility.Float64, err = strconv.ParseFloat(quote.ImpliedVolatility, 64); err == nil {
		ImpliedVolatility.Valid = true
	}
	///////////////////////
	var LastTradePrice sql.NullFloat64
	if LastTradePrice.Float64, err = strconv.ParseFloat(quote.LastTrade, 64); err == nil {
		LastTradePrice.Valid = true
	}
	///////////////////////
	var LastTradeTimestamp sql.NullInt64
	// example: 2018-09-21T00:00:00-04:00
	if t, err := time.Parse("2006-01-02T15:04:05-07:00", quote.LastTradeTimestamp); err == nil {
		LastTradeTimestamp.Int64 = t.Unix() //number of seconds since epoch
		LastTradeTimestamp.Valid = true
	}
	///////////////////////
	var LastTradeVolume sql.NullInt64
	if LastTradeVolume.Int64, err = strconv.ParseInt(quote.LastTradeVolume, 10, 64); err == nil {
		LastTradeVolume.Valid = true
	}
	///////////////////////
	var LowPrice sql.NullFloat64
	if LowPrice.Float64, err = strconv.ParseFloat(quote.LowPrice, 64); err == nil {
		LowPrice.Valid = true
	}
	///////////////////////
	var OpenInterest sql.NullInt64
	if OpenInterest.Int64, err = strconv.ParseInt(quote.OpenInterest, 10, 64); err == nil {
		OpenInterest.Valid = true
	}
	///////////////////////
	var Rho sql.NullFloat64
	if Rho.Float64, err = strconv.ParseFloat(quote.Rho, 64); err == nil {
		Rho.Valid = true
	}
	///////////////////////
	var Theta sql.NullFloat64
	if Theta.Float64, err = strconv.ParseFloat(quote.Theta, 64); err == nil {
		Theta.Valid = true
	}
	///////////////////////
	var Vega sql.NullFloat64
	if Vega.Float64, err = strconv.ParseFloat(quote.Vega, 64); err == nil {
		Vega.Valid = true
	}

	return &core.SibylOptionQuoteRecord{
		Ask:                Ask,
		AskTime:            AskTime,
		AskSize:            AskSize,
		Bid:                Bid,
		BidTime:            BidTime,
		BidSize:            BidSize,
		Change:             Change,
		Delta:              Delta,
		EquityType:         optionSymbolType.OptionType,
		Expiration:         optionSymbolType.Expiration,
		Gamma:              Gamma,
		HighPrice:          HighPrice,
		ImpliedVolatility:  ImpliedVolatility,
		LastTradePrice:     LastTradePrice,
		LastTradeVolume:    LastTradeVolume,
		LastTradeTimestamp: LastTradeTimestamp,
		LowPrice:           LowPrice,
		OpenInterest:       OpenInterest,
		Rho:                Rho,
		Symbol:             optionSymbolType.Symbol,
		StrikePrice:        optionSymbolType.StrikePrice,
		Theta:              Theta,
		Timestamp:          core.NewTimestampTypeFromUnix(Timestamp),
		Vega:               Vega,
	}, nil
}

///////////////////////////

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ag *AllyAgent) GetQuotes(ctx context.Context, stockSymbols map[core.StockSymbolType]bool, optionSymbols map[core.OptionSymbolType]bool) ([]*core.SibylStockQuoteRecord, []*core.SibylOptionQuoteRecord, error) {
	startTime := time.Now()
	emptyStockRecords := make([]*core.SibylStockQuoteRecord, 0)
	emptyOptionRecords := make([]*core.SibylOptionQuoteRecord, 0)

	if len(stockSymbols) == 0 && len(optionSymbols) == 0 {
		logrus.Debugf("GetQuotes: nothing to find, finished in %s", time.Since(startTime))
		return emptyStockRecords, emptyOptionRecords, nil
	}

	//add the symbols to query
	symbolStringSlice := make([]string, 0, len(stockSymbols)+len(optionSymbols))
	for stockSymbol := range stockSymbols {
		symbolStringSlice = append(symbolStringSlice, string(stockSymbol))
	}
	sortedOptions := make([]core.OptionSymbolType, 0, len(optionSymbols))

	for optionSymbol := range optionSymbols {
		sortedOptions = append(sortedOptions, optionSymbol)
	}
	// because ally has a timeout we want to sort so the priority is the following
	// stock then options.  Since options are less like to change a fequently when they are further out in time
	// options are prioritized by expiration date.
	sort.Slice(sortedOptions, func(i, j int) bool {
		return sortedOptions[i].Expiration.Unix() < sortedOptions[j].Expiration.Unix()
	})

	for _, optionSymbol := range sortedOptions {
		symbolStringSlice = append(symbolStringSlice, toAllySymbol(optionSymbol))
	}
	// since in general anything index over 60k will most like never get hit
	// so we will randomly mix those up so we get them occasionally
	if len(symbolStringSlice) > 50000 {
		// it is assumed that we the last 10k (60k-50k) will get
		// through we'll shuffle everything beyond 50k
		rand.Shuffle(len(symbolStringSlice)-50000, func(i, j int) {
			ii := i + 50000
			jj := j + 50000
			tmp := symbolStringSlice[ii]
			symbolStringSlice[ii] = symbolStringSlice[jj]
			symbolStringSlice[jj] = tmp
		})
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
	// request will time out.  After manually testing it, the best seems
	// to be around the 8k mark
	//TODO maybe consider a monte carlo method to explore other values and act greedy as on online policy
	maxSymbolCount := 8000

	if len(symbolStringSlice)%maxSymbolCount == 1 {
		//and we make sure that the last request has more than one symbol in it
		//we make all of the requests a bit over stuffed but the maxSymbolCount is below
		// the real max so should be a problem
		maxSymbolCount++
	}

	symbolGroups := make([][]string, 0)
	for i := 0; i < len(symbolStringSlice); i += maxSymbolCount {
		endIndex := min(i+maxSymbolCount, len(symbolStringSlice))
		symbolGroups = append(symbolGroups, symbolStringSlice[i:endIndex:endIndex])
	}

	//ally seems to only like one of these requests at a time
	// so we'll do them serially... since things are randomized and the
	// law of large number it should matter way too much
	errStrings := []string{}
	toReturnStockRecords := make([]*core.SibylStockQuoteRecord, 0, len(stockSymbols))
	toReturnOptionRecords := make([]*core.SibylOptionQuoteRecord, 0, len(optionSymbols))
	waitGroup := sync.WaitGroup{}
	fieldsString := strings.Join(quoteRequestFields, ",")
	for _, symbols := range symbolGroups {
		//TODO consider ways to reduce the repetitive cost of creating a new request
		data := url.Values{}
		// add the symbols to the query
		data.Add("symbols", strings.Join(symbols, ","))
		//next specify the fields we're interested in; this will help reduce bandwidth usage
		data.Add("fids", fieldsString)

		request, err := http.NewRequest(http.MethodPost, "https://api.tradeking.com/v1/market/ext/quotes.xml", strings.NewReader(data.Encode()))
		if logrus.GetLevel() == logrus.DebugLevel {
			if dump, err := httputil.DumpRequest(request, true); err != nil {
				logrus.Errorf("GetQuotes: there was a problem with dumping the request: %v", err)
			} else {
				logrus.Debugf("GetQuotes: the request:%v", string(dump))
			}
		}

		if err != nil {
			return emptyStockRecords, emptyOptionRecords, fmt.Errorf("GetQuotes: request creation error: %v", err)
		}

		//we rate limit the quote calls and give each request one shot to be successful
		ag.rateLimitMarketCalls.Take(ctx) // we limit frequeny of market calls
		ag.concurrentLimit.Take(ctx)      // and we limit concurrent requests
		resp, err := ctxhttp.Do(ctx, ag.httpClient, request)
		ag.concurrentLimit.Return()
		if err != nil {
			errStrings = append(errStrings, fmt.Sprintf("GetQuotes: client error: %v", err))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			//TODO consider creating a queue for retries (for 429 since they return fast), and if there is enough time at the end retry them
			errStrings = append(errStrings, fmt.Sprintf("GetQuotes: client error with status code %v: %v", resp.StatusCode, resp.Status))
			continue
		}

		waitGroup.Add(1)
		go func(resp *http.Response) {
			//TODO performance: this loop is slow consider something faster
			//currently we spin the processing of the response off into a go routine of its
			// own because it's kind of slow and impacts the request rate (which is the long pole in the tent)
			defer waitGroup.Done()

			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()

			logrus.Debugf("GetQuotes: response body(%v): %v", len(body), string(body))

			//response, err := createAllyJsonStockQuoteResponse(string(body))
			response, err := createAllyXMLStockQuoteResponse(string(body))
			if err != nil {
				tmp := err.Error()
				errStrings = append(errStrings, fmt.Sprintf("GetQuotes: unmarshal error: %v", tmp[:min(len(tmp), 1000)]))
				return
			}

			for _, quote := range response.Quotes.Quote {
				//so for each quote we could have either a stock or option
				// to check we take the quotes' symbol and see if it's in the the stocks list
				if _, has := stockSymbols[core.StockSymbolType(quote.Symbol)]; has {
					if sq, err := allyXMLQuoteToSibylStockQuoteRecord(quote); err != nil {
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
					if sq, err := allyXMLQuoteToSibylOptionQuoteRecord(quote); err != nil {
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
		}(resp)

	}
	waitGroup.Wait()

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
		err = fmt.Errorf("GetQuotes: had errors while parsing quotes: %v", strings.Join(errStrings, ";"))
	}

	logrus.Debugf("GetQuotes: finished found %v stock quotes and %v option quotes in %s", len(toReturnStockRecords), len(toReturnOptionRecords), time.Since(startTime))

	return toReturnStockRecords, toReturnOptionRecords, err
}
