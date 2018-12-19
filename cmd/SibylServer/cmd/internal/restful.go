package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/nathanhack/sibyl/rest"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
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
	router.HandleFunc("/stocks/enable/all", serverContext.StockEnableAll).Methods(http.MethodPut)
	router.HandleFunc("/stocks/disable/all", serverContext.StockDisableAll).Methods(http.MethodPut)

	router.HandleFunc("/stocks/revalidate/{stockSymbol}", serverContext.StockRevalidate).Methods(http.MethodPut)

	router.HandleFunc("/history/{stockSymbol}/{startTimestamp}/{endTimestamp}", serverContext.HistoryGet).Methods(http.MethodGet)

	router.HandleFunc("/intraday/{stockSymbol}/{startTimestamp}/{endTimestamp}", serverContext.IntradayGet).Methods(http.MethodGet)

	router.HandleFunc("/database/download/history", serverContext.DatabaseDownloadHistory).Methods(http.MethodGet)
	router.HandleFunc("/database/download/intraday", serverContext.DatabaseDownloadIntraday).Methods(http.MethodGet)
	router.HandleFunc("/database/download/stocks/quote", serverContext.DatabaseDownloadStockQuote).Methods(http.MethodGet)
	router.HandleFunc("/database/download/stocks/stable", serverContext.DatabaseDownloadStockStable).Methods(http.MethodGet)
	router.HandleFunc("/database/download/options/quote", serverContext.DatabaseDownloadOptionsQuote).Methods(http.MethodGet)
	router.HandleFunc("/database/download/options/stable", serverContext.DatabaseDownloadOptionStable).Methods(http.MethodGet)

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
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	return server, nil
}
func (sc *ServerContext) StockDisableAll(writer http.ResponseWriter, request *http.Request) {
	err := sc.db.StockDisableAll(sc.Ctx)
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {

		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableAll(writer http.ResponseWriter, request *http.Request) {
	err := sc.db.StockEnableAll(sc.Ctx)
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDisableStableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableStableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockEnableStableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableStableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockDisableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableQuotes(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableQuotes(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockDisableIntraday(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableIntradayHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableIntraday(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableIntradayHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockDisableHistory(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableHistory(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableHistory(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}

func (sc *ServerContext) StockDelete(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	stockSymbol := params["stockSymbol"]
	err := sc.db.StockDelete(sc.Ctx, core.StockSymbolType(stockSymbol))

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockRevalidate(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockRevalidate(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))

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
		return rest.ErrorState{err.Error(), true}
	}
	return rest.ErrorState{"", false}
}

func (sc *ServerContext) StockAdd(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	stockSymbol := params["stockSymbol"]
	err := sc.db.StockAdd(sc.Ctx, core.StockSymbolType(stockSymbol))

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockEnableDownloading(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockEnableDownloading(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		sc.stockValidator.RequestUpdate <- true
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
}
func (sc *ServerContext) StockDisableDownloading(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := sc.db.StockDisableDownloading(sc.Ctx, core.StockSymbolType(params["stockSymbol"]))
	errs := make([]string, 0)
	if err != nil {
		errs = append(errs, err.Error())
	}

	//now write it out as the response
	json.NewEncoder(writer).Encode(errToRestErrorState(err))
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

	sc.stockValidator.RequestUpdate <- true

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

	sc.stockValidator.RequestUpdate <- true

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
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpHistoryRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{Histories: string(bytes), ErrorState: errToRestErrorState(nil)})
}

func (sc *ServerContext) DatabaseDownloadIntraday(writer http.ResponseWriter, request *http.Request) {
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpIntradayRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{Intradays: string(bytes), ErrorState: errToRestErrorState(nil)})
}

func (sc *ServerContext) DatabaseDownloadStockQuote(writer http.ResponseWriter, request *http.Request) {
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpStockQuoteRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{StockQuotes: string(bytes), ErrorState: errToRestErrorState(nil)})

}

func (sc *ServerContext) DatabaseDownloadStockStable(writer http.ResponseWriter, request *http.Request) {
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpStableStockQuoteRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{StockStableQuotes: string(bytes), ErrorState: errToRestErrorState(nil)})
}

func (sc *ServerContext) DatabaseDownloadOptionsQuote(writer http.ResponseWriter, request *http.Request) {
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpOptionQuoteRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{OptionQuotes: string(bytes), ErrorState: errToRestErrorState(nil)})
}

func (sc *ServerContext) DatabaseDownloadOptionStable(writer http.ResponseWriter, request *http.Request) {
	//so in order to do this we use the dump functionality, pass in a tmp file then read the file
	// in and send it in the message (this doesn't have to be fast but it does have to work)
	tmpFile, err := ioutil.TempFile("./", "dump*")
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}
	defer os.Remove(tmpFile.Name())

	err = tmpFile.Close()
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	sc.db.DumpStableOptionQuoteRecordsToFile(sc.Ctx, tmpFile.Name())
	bytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{OptionStableQuotes: string(bytes), ErrorState: errToRestErrorState(nil)})
}

func databaseUploadStage(request *http.Request) (string, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}
	request.Body.Close()

	var databaseRecords rest.DatabaseRecords
	if err := json.Unmarshal(body, &databaseRecords); err != nil {
		return "", err
	}

	var fileBody string
	if len(databaseRecords.Histories) != 0 {
		fileBody = databaseRecords.Histories
	} else if len(databaseRecords.Intradays) != 0 {
		fileBody = databaseRecords.Intradays
	} else if len(databaseRecords.OptionStableQuotes) != 0 {
		fileBody = databaseRecords.OptionStableQuotes
	} else if len(databaseRecords.OptionQuotes) != 0 {
		fileBody = databaseRecords.OptionQuotes
	} else if len(databaseRecords.StockQuotes) != 0 {
		fileBody = databaseRecords.StockQuotes
	} else if len(databaseRecords.StockStableQuotes) != 0 {
		fileBody = databaseRecords.StockStableQuotes
	} else {
		return "", fmt.Errorf("nothing to parse")
	}

	tmpFile, err := ioutil.TempFile("./", "load*")
	if err != nil {
		return "", err
	}

	n, err := tmpFile.WriteString(fileBody)
	if err != nil {
		return "", err
	}
	if n != len(databaseRecords.Histories) {
		return "", fmt.Errorf("unable to stage before ingesting")
	}

	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil

}

func (sc *ServerContext) DatabaseUploadHistory(writer http.ResponseWriter, request *http.Request) {

	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadHistoryRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	return

}
func (sc *ServerContext) DatabaseUploadIntraday(writer http.ResponseWriter, request *http.Request) {
	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadIntradayRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	return
}
func (sc *ServerContext) DatabaseUploadStockQuote(writer http.ResponseWriter, request *http.Request) {
	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadStockQuoteRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	return
}
func (sc *ServerContext) DatabaseUploadStockStable(writer http.ResponseWriter, request *http.Request) {
	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadStableStockQuoteRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	return
}
func (sc *ServerContext) DatabaseUploadOptionsQuote(writer http.ResponseWriter, request *http.Request) {
	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadOptionQuoteRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
	return
}
func (sc *ServerContext) DatabaseUploadOptionStable(writer http.ResponseWriter, request *http.Request) {
	if fileName, err := databaseUploadStage(request); err != nil {
		json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
		return
	} else {
		defer os.Remove(fileName)
		if err := sc.db.LoadStableOptionQuoteRecordsFromFile(sc.Ctx, fileName); err != nil {
			json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(err)})
			return
		}
	}

	json.NewEncoder(writer).Encode(rest.DatabaseRecords{ErrorState: errToRestErrorState(nil)})
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
