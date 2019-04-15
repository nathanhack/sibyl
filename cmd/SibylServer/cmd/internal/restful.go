package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/nathanhack/sibyl/rest"
	"github.com/sirupsen/logrus"
)

type ServerContext struct {
	Ctx            context.Context
	db             *database.SibylDatabase
	stockValidator *StockValidator
}

func makeServer(serverContext *ServerContext, serverAddress string) (*http.Server, error) {
	router := mux.NewRouter()

	router.HandleFunc("/stocks/get", serverContext.StocksGetAll).Methods(http.MethodGet)
	router.HandleFunc("/stocks/add/{stockSymbol}", serverContext.StockAdd).Methods(http.MethodPost)
	router.HandleFunc("/stocks/delete/{stockSymbol}", serverContext.StockDelete).Methods(http.MethodDelete)

	router.HandleFunc("/stocks/enable/downloading/{stockSymbol}", serverContext.StockEnableDownloading).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/downloading/{stockSymbol}", serverContext.StockDisableDownloading).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/history/{stockSymbol}", serverContext.StockEnableHistory).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/history/{stockSymbol}", serverContext.StockDisableHistory).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/intraday/{stockSymbol}", serverContext.StockEnableIntraday).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/intraday/{stockSymbol}", serverContext.StockDisableIntraday).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/quotes/{stockSymbol}", serverContext.StockEnableQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/quotes/{stockSymbol}", serverContext.StockDisableQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/stableQuotes/{stockSymbol}", serverContext.StockEnableStableQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/stableQuotes/{stockSymbol}", serverContext.StockDisableStableQuotes).Methods(http.MethodPut)

	router.HandleFunc("/stocks/enable/all/downloading", serverContext.StockEnableAllDownloading).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all/downloading", serverContext.StockDisableAllDownloading).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/all/history", serverContext.StockEnableAllHistory).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all/history", serverContext.StockDisableAllHistory).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/all/intraday", serverContext.StockEnableAllIntraday).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all/intraday", serverContext.StockDisableAllIntraday).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/all/quotes", serverContext.StockEnableAllQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all/quotes", serverContext.StockDisableAllQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/all/stableQuotes", serverContext.StockEnableAllStableQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all/stableQuotes", serverContext.StockDisableAllStableQuotes).Methods(http.MethodPut)
	router.HandleFunc("/stocks/enable/all", serverContext.StockEnableAll).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all", serverContext.StockDisableAll).Methods(http.MethodPut)

	router.HandleFunc("/stocks/revalidate/{stockSymbol}", serverContext.StockRevalidate).Methods(http.MethodPut)

	router.HandleFunc("/history/{stockSymbol}/{startTimestamp}/{endTimestamp}", serverContext.HistoryGet).Methods(http.MethodGet)

	router.HandleFunc("/intraday/{stockSymbol}/{startTimestamp}/{endTimestamp}", serverContext.IntradayGet).Methods(http.MethodGet)

	router.HandleFunc("/database/download/history/{lastID:.*}", serverContext.DatabaseDownloadHistory).Methods(http.MethodGet)
	router.HandleFunc("/database/download/intraday/{lastID:.*}", serverContext.DatabaseDownloadIntraday).Methods(http.MethodGet)
	router.HandleFunc("/database/download/stocks/quote/{lastID:.*}", serverContext.DatabaseDownloadStockQuote).Methods(http.MethodGet)
	router.HandleFunc("/database/download/stocks/stable/{lastID:.*}", serverContext.DatabaseDownloadStockStable).Methods(http.MethodGet)
	router.HandleFunc("/database/download/options/quote/{lastID:.*}", serverContext.DatabaseDownloadOptionsQuote).Methods(http.MethodGet)
	router.HandleFunc("/database/download/options/stable/{lastID:.*}", serverContext.DatabaseDownloadOptionStable).Methods(http.MethodGet)

	router.HandleFunc("/database/upload/history", serverContext.DatabaseUploadHistory).Methods(http.MethodPost)
	router.HandleFunc("/database/upload/intraday", serverContext.DatabaseUploadIntraday).Methods(http.MethodPost)
	router.HandleFunc("/database/upload/stocks/quote", serverContext.DatabaseUploadStockQuote).Methods(http.MethodPost)
	router.HandleFunc("/database/upload/stocks/stable", serverContext.DatabaseUploadStockStable).Methods(http.MethodPost)
	router.HandleFunc("/database/upload/options/quote", serverContext.DatabaseUploadOptionsQuote).Methods(http.MethodPost)
	router.HandleFunc("/database/upload/options/stable", serverContext.DatabaseUploadOptionStable).Methods(http.MethodPost)

	//Agent related///////////////
	router.HandleFunc("/agent/ally/{consumerKey}/{consumerSecret}/{oAuthToken}/{oAuthTokenSecret}", serverContext.AgentAllyCreds).Methods(http.MethodPost)
	router.HandleFunc("/agent/use/ally", serverContext.AgentUseAlly).Methods(http.MethodPut)
	router.HandleFunc("/agent/creds", serverContext.AgentCreds).Methods(http.MethodGet)

	router.HandleFunc("/agent/use/tdameritrade", serverContext.AgentUseTdAmeritrade).Methods(http.MethodPut)

	urlR, err := url.Parse(serverAddress)
	logrus.Infof("Server on: %v", urlR)
	if err != nil {
		logrus.Errorf("Failed to parse localhost URL")
		return nil, fmt.Errorf("failed to parse localhost URL")
	}
	server := &http.Server{
		Addr: urlR.Host,
		// Good practice to set timeouts to avoid Slowloris attacks.
		// TODO make these timeouts part of the future configurations
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	return server, nil
}
func (sc *ServerContext) StockDisableAll(writer http.ResponseWriter, request *http.Request) {
	err := sc.db.StockDisableAll(sc.Ctx)

	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableAll(writer http.ResponseWriter, request *http.Request) {
	err := sc.db.StockEnableAll(sc.Ctx)
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableStableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableStableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableAllStableQuotes(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockDisableStableQuotes(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockEnableStableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableStableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableAllStableQuotes(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockEnableStableQuotes(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockDisableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))

	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableAllQuotes(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockDisableQuotes(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockEnableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableAllQuotes(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockEnableQuotes(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockDisableIntraday(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableIntradayHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableAllIntraday(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockDisableIntradayHistory(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockEnableIntraday(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableIntradayHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableAllIntraday(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockEnableIntradayHistory(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockDisableHistory(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableAllHistory(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockDisableHistory(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockEnableHistory(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableAllHistory(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, stockRecord := range records {
		err := sc.db.StockEnableHistory(sc.Ctx, stockRecord.Symbol)
		errs := make([]string, 0)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockDelete(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	stockSymbol := params["stockSymbol"]
	err := sc.db.StockDelete(sc.Ctx, core.StockSymbolType(stockSymbol))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockRevalidate(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockRevalidate(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StocksGetAll(writer http.ResponseWriter, request *http.Request) {
	dbStocks, err := sc.db.GetAllStockRecords(sc.Ctx)

	stocks := make([]rest.Stock, 0, len(dbStocks))
	for _, dbStock := range dbStocks {
		stocks = append(stocks, rest.Stock{
			Exchange:              string(dbStock.Exchange),
			ExchangeDescription:   string(dbStock.ExchangeDescription),
			DownloadStatus:        string(dbStock.DownloadStatus),
			HasOptions:            dbStock.HasOptions,
			HistoryStatus:         string(dbStock.HistoryStatus),
			IntradayHistoryStatus: string(dbStock.IntradayHistoryStatus),
			IntradayHistoryState:  string(dbStock.IntradayHistoryState),
			Name:                  string(dbStock.Name),
			QuotesStatus:          string(dbStock.QuotesStatus),
			StableQuotesStatus:    string(dbStock.StableQuotesStatus),
			Symbol:                string(dbStock.Symbol),
			Validation:            string(dbStock.ValidationStatus),
		})
	}

	toReturn := rest.StocksResponse{
		Stocks:     stocks,
		ErrorState: errToRestErrorState(err),
	}
	if err != nil {
		logrus.Errorf("StockGetAll: had an error: %v", err)
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(toReturn)
}

func errToRestErrorState(err error) rest.ErrorState {
	if err != nil {
		return rest.ErrorState{Error: err.Error(), ErrorReturned: true}
	}
	return rest.ErrorState{Error: "", ErrorReturned: false}
}

func errsToRestErrorState(errs []string) rest.ErrorState {
	if len(errs) > 0 {
		return rest.ErrorState{Error: fmt.Sprint(errs), ErrorReturned: true}
	}
	return rest.ErrorState{Error: "", ErrorReturned: false}
}

func (sc *ServerContext) StockAdd(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	stockSymbol := params["stockSymbol"]
	err := sc.db.StockAdd(sc.Ctx, core.StockSymbolType(stockSymbol))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableDownloading(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableDownloading(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableAllDownloading(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, stockRecord := range records {
		err := sc.db.StockEnableDownloading(sc.Ctx, stockRecord.Symbol)

		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true

	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) StockDisableDownloading(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableDownloading(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	if err == nil {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableAllDownloading(writer http.ResponseWriter, request *http.Request) {
	records, err := sc.db.GetAllStockRecords(sc.Ctx)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	errs := make([]string, 0)
	for _, record := range records {
		err := sc.db.StockDisableDownloading(sc.Ctx, record.Symbol)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	sc.stockValidator.RequestUpdate <- true
	//now write it out as the response
	json.NewEncoder(writer).Encode(errsToRestErrorState(errs))
}

func (sc *ServerContext) AgentAllyCreds(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// there will be four params from the endpoint .../{consumerKey}/{consumerSecret}/{oAuthToken}/{oAuthTokenSecret}

	creds, err := sc.db.GetCreds(sc.Ctx)
	if err != nil {
		//if we're here then the database should be working the error
		// would be due an empty database so we'll substitute in with the default
		creds = core.DefaultSibylCreds()
	}

	//now create a new cred struct with the updated data
	creds = core.NewSibylCreds(
		creds.AgentSelection(),
		params["consumerKey"],
		params["consumerSecret"],
		params["oAuthToken"],
		params["oAuthTokenSecret"],
		creds.UrlRedirect(),
		creds.AccessToken(),
		creds.RefreshToken(),
		creds.ExpireTimestamp(),
		creds.RefreshExpireTimestamp(),
	)

	//send the error value as the result
	json.NewEncoder(writer).Encode(errToRestErrorState(sc.db.LoadCreds(sc.Ctx, creds)))
}

func (sc *ServerContext) AgentUseAlly(writer http.ResponseWriter, request *http.Request) {
	creds, err := sc.db.GetCreds(sc.Ctx)
	if err != nil {
		//if we're here then the database should be working the error
		// would be due an empty database so we'll substitute in with the default
		creds = core.DefaultSibylCreds()
	}

	//now create a new cred struct with the updated data
	creds = core.NewSibylCreds(
		core.AgentSelectionAlly,
		creds.ConsumerKey(),
		creds.ConsumerSecret(),
		creds.Token(),
		creds.TokenSecret(),
		creds.UrlRedirect(),
		creds.AccessToken(),
		creds.RefreshToken(),
		creds.ExpireTimestamp(),
		creds.RefreshExpireTimestamp(),
	)

	json.NewEncoder(writer).Encode(errToRestErrorState(sc.db.LoadCreds(sc.Ctx, creds)))
}

func (sc *ServerContext) AgentUseTdAmeritrade(writer http.ResponseWriter, request *http.Request) {
	creds, err := sc.db.GetCreds(sc.Ctx)
	if err != nil {
		//if we're here then the database should be working the error
		// would be due an empty database so we'll substitute in with the default
		creds = core.DefaultSibylCreds()
	}

	//now create a new cred struct with the updated data
	creds = core.NewSibylCreds(
		core.AgentSelectionTDAmeritrade,
		creds.ConsumerKey(),
		creds.ConsumerSecret(),
		creds.Token(),
		creds.TokenSecret(),
		creds.UrlRedirect(),
		creds.AccessToken(),
		creds.RefreshToken(),
		creds.ExpireTimestamp(),
		creds.RefreshExpireTimestamp(),
	)

	json.NewEncoder(writer).Encode(errToRestErrorState(sc.db.LoadCreds(sc.Ctx, creds)))
}
func (sc *ServerContext) HistoryGet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// .../{stockSymbol}/{startTimestamp}/{endTimestamp}
	var startTimestamp, endTimestamp int64
	var err error

	stockSymbol := core.StockSymbolType(params["stockSymbol"])
	startTimestamp, err = strconv.ParseInt(params["startTimestamp"], 10, 64)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}
	endTimestamp, err = strconv.ParseInt(params["endTimestamp"], 10, 64)
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	var historyRecords []*core.SibylHistoryRecord
	historyRecords, err = sc.db.GetHistory(sc.Ctx, stockSymbol, core.NewDateTypeFromUnix(startTimestamp), core.NewDateTypeFromUnix(endTimestamp))
	if err != nil {
		json.NewEncoder(writer).Encode(errToRestErrorState(err))
		return
	}

	restHistories := make([]rest.History, 0, len(historyRecords))
	for _, record := range historyRecords {
		restHistories = append(restHistories, rest.History{
			ClosePrice: nullFloat64ToString(record.ClosePrice, ""),
			HighPrice:  nullFloat64ToString(record.HighPrice, ""),
			LowPrice:   nullFloat64ToString(record.LowPrice, ""),
			OpenPrice:  nullFloat64ToString(record.OpenPrice, ""),
			Symbol:     string(record.Symbol),
			Timestamp:  record.Timestamp.Unix(),
			Volume:     nullInt64ToString(record.Volume, ""),
		})
	}

	json.NewEncoder(writer).Encode(rest.Histories{Histories: restHistories, ErrorState: errToRestErrorState(nil)})
}
func (sc *ServerContext) IntradayGet(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	// .../{stockSymbol}/{startTimestamp}/{endTimestamp}
	var startTimestamp, endTimestamp int64
	var err error

	stockSymbol := core.StockSymbolType(params["stockSymbol"])
	startTimestamp, err = strconv.ParseInt(params["startTimestamp"], 10, 64)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.Intradays{Intradays: []rest.Intraday{}, ErrorState: errToRestErrorState(err)})
		return
	}
	endTimestamp, err = strconv.ParseInt(params["endTimestamp"], 10, 64)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.Intradays{Intradays: []rest.Intraday{}, ErrorState: errToRestErrorState(err)})
		return
	}

	var intradayRecords []*core.SibylIntradayRecord
	intradayRecords, err = sc.db.GetIntraday(sc.Ctx, stockSymbol, core.NewTimestampTypeFromUnix(startTimestamp), core.NewTimestampTypeFromUnix(endTimestamp))
	if err != nil {
		json.NewEncoder(writer).Encode(rest.Intradays{Intradays: []rest.Intraday{}, ErrorState: errToRestErrorState(err)})
		return
	}

	restIntradays := make([]rest.Intraday, 0, len(intradayRecords))
	for _, record := range intradayRecords {
		restIntradays = append(restIntradays, rest.Intraday{
			HighPrice: nullFloat64ToString(record.HighPrice, ""),
			LastPrice: nullFloat64ToString(record.LastPrice, ""),
			LowPrice:  nullFloat64ToString(record.LowPrice, ""),
			OpenPrice: nullFloat64ToString(record.OpenPrice, ""),
			Symbol:    string(record.Symbol),
			Timestamp: record.Timestamp.Unix(),
			Volume:    nullInt64ToString(record.Volume, ""),
		})
	}

	json.NewEncoder(writer).Encode(rest.Intradays{Intradays: restIntradays, ErrorState: errToRestErrorState(nil)})
}

func (sc *ServerContext) DatabaseDownloadHistory(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	nextLastID, buffer, err := sc.db.DumpRangeHistoryRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadHistory: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastHistoryRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{Histories: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadHistory: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadHistory: successfully downloaded(bytes:%v)", len(buffer))
}

func (sc *ServerContext) DatabaseDownloadIntraday(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	nextLastID, buffer, err := sc.db.DumpRangeIntradayRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadIntraday: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastIntradayRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{Intradays: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadIntraday: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadIntraday: successfully downloaded(bytes:%v)", len(buffer))
}

func (sc *ServerContext) DatabaseDownloadStockQuote(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	nextLastID, buffer, err := sc.db.DumpRangeStockQuoteRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadStockQuote: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastStockQuoteRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{StockQuotes: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadStockQuote: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadStockQuote: successfully downloaded(bytes:%v)", len(buffer))
}

func (sc *ServerContext) DatabaseDownloadStockStable(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	nextLastID, buffer, err := sc.db.DumpRangeStableStockQuoteRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadStockStable: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastStableStockQuoteRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{StockStableQuotes: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadStockStable: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadStockStable: successfully downloaded(bytes:%v)", len(buffer))
}

func (sc *ServerContext) DatabaseDownloadOptionsQuote(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	nextLastID, buffer, err := sc.db.DumpRangeOptionQuoteRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadOptionsQuote: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastOptionQuoteRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{OptionQuotes: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadOptionsQuote: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadOptionsQuote: successfully downloaded(bytes:%v)", len(buffer))
}

func (sc *ServerContext) DatabaseDownloadOptionStable(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	nextLastID, buffer, err := sc.db.DumpRangeStableOptionQuoteRecordsToBuffer(sc.Ctx, params["lastID"], 10000)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadOptionStable: had a error: %v", err)
		return
	}

	hasMore := false
	if nextLastID != "" {
		if dbLastID, err := sc.db.LastStableOptionQuoteRecordID(sc.Ctx); err == nil {
			hasMore = nextLastID != dbLastID
		} else {
			hasMore = true
		}
	}

	if err := json.NewEncoder(writer).Encode(rest.DatabaseRecords{OptionStableQuotes: buffer, LastID: nextLastID, More: hasMore, ErrorState: errToRestErrorState(nil)}); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseDownloadOptionStable: had a error: %v", err)
		return
	}
	logrus.Infof("DatabaseDownloadOptionStable: successfully downloaded(bytes:%v)", len(buffer))
}

func readDatabaseRecords(request *http.Request) (*rest.DatabaseRecords, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("problems with reading body of request: %v", err)
	}
	request.Body.Close()

	var databaseRecords rest.DatabaseRecords
	if err := json.Unmarshal(body, &databaseRecords); err != nil {
		return nil, fmt.Errorf("problems unmarshalling json: %v", err)
	}
	return &databaseRecords, nil
}

func (sc *ServerContext) DatabaseUploadHistory(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadHistory: had a error: %v", err)
		return
	}

	if err := sc.db.LoadHistoryRecordsFromFileContents(sc.Ctx, databaseRecords.Histories); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadHistory: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadHistory: successfully upload(bytes:%v)", len(databaseRecords.Histories))
	return
}

func (sc *ServerContext) DatabaseUploadIntraday(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadIntraday: had a error: %v", err)
		return
	}

	if err := sc.db.LoadIntradayRecordsFromFileContents(sc.Ctx, databaseRecords.Intradays); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadIntraday: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadIntraday: successfully upload(bytes:%v)", len(databaseRecords.Intradays))
	return

}

func (sc *ServerContext) DatabaseUploadStockQuote(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadStockQuote: had a error: %v", err)
		return
	}

	if err := sc.db.LoadStockQuoteRecordsFromFileContents(sc.Ctx, databaseRecords.StockQuotes); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadStockQuote: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadStockQuote: successfully upload(bytes:%v)", len(databaseRecords.StockQuotes))
	return

}

func (sc *ServerContext) DatabaseUploadStockStable(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadStockStable: had a error: %v", err)
		return
	}

	if err := sc.db.LoadStableStockQuoteRecordsFromFileContents(sc.Ctx, databaseRecords.StockStableQuotes); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadStockStable: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadStockStable: successfully upload(bytes:%v)", len(databaseRecords.StockStableQuotes))
	return

}

func (sc *ServerContext) DatabaseUploadOptionsQuote(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadOptionsQuote: had a error: %v", err)
		return
	}

	if err := sc.db.LoadOptionQuoteRecordsFromFileContents(sc.Ctx, databaseRecords.OptionQuotes); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadOptionsQuote: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadOptionsQuote: successfully upload(bytes:%v)", len(databaseRecords.OptionQuotes))
	return
}

func (sc *ServerContext) DatabaseUploadOptionStable(writer http.ResponseWriter, request *http.Request) {
	databaseRecords, err := readDatabaseRecords(request)
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadOptionStable: had a error: %v", err)
		return
	}

	if err := sc.db.LoadStableOptionQuoteRecordsFromFileContents(sc.Ctx, databaseRecords.OptionStableQuotes); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		logrus.Errorf("DatabaseUploadOptionStable: had a error: %v", err)
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	logrus.Infof("DatabaseUploadOptionStable: successfully upload(bytes:%v)", len(databaseRecords.OptionStableQuotes))
	return
}

func (sc *ServerContext) AgentCreds(writer http.ResponseWriter, request *http.Request) {

	if creds, err := sc.db.GetCreds(sc.Ctx); err == nil {
		json.NewEncoder(writer).Encode(rest.Creds{
			AgentSelection:  string(creds.AgentSelection()),
			ConsumerKey:     creds.ConsumerKey(),
			ConsumerSecret:  creds.ConsumerSecret(),
			Token:           creds.Token(),
			TokenSecret:     creds.TokenSecret(),
			UrlRedirect:     creds.UrlRedirect(),
			AccessToken:     creds.AccessToken(),
			RefreshToken:    creds.RefreshToken(),
			ExpireTimestamp: creds.ExpireTimestamp(),
			ErrorState:      errToRestErrorState(err)})
	} else {
		json.NewEncoder(writer).Encode(rest.Creds{ErrorState: errToRestErrorState(err)})
	}
}

func nullFloat64ToString(v sql.NullFloat64, nullString string) string {
	if v.Valid {
		return fmt.Sprintf("%v", v.Float64)
	}
	return nullString
}

func nullInt64ToString(v sql.NullInt64, nullString string) string {
	if v.Valid {
		return fmt.Sprintf("%v", v.Int64)
	}
	return nullString
}
