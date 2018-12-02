package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

func TestSibylDatabase_DumpCredsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpCredsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_DumpHistoryRecordsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpHistoryRecordsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_DumpOptionQuoteRecordsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpOptionQuoteRecordsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_DumpStableOptionQuoteRecordsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpStableOptionQuoteRecordsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_DumpStableStockQuoteRecordsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpStableStockQuoteRecordsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_DumpStockQuoteRecordsToFile(t *testing.T) {
	filePathname := "./testFile.out"
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sd, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("expected not to fail:%v", err)
		return
	}

	if err := sd.DumpStockQuoteRecordsToFile(ctx, filePathname); err != nil {
		t.Errorf("expected not to fail:%v", err)
	} else {
		if _, err := os.Stat(filePathname); os.IsExist(err) {
			if err := os.Remove(filePathname); err != nil {
				t.Errorf("there was an error while trying to remove %v :%v", filePathname, err)
			}
		} else {
			t.Errorf("expected file to exist however no file found")
		}
	}
}

func TestSibylDatabase_LoadOptionQuoteRecords(t *testing.T) {
	quotes := []*core.SibylOptionQuoteRecord{
		{
			Ask:                sql.NullFloat64{4.3, true},
			AskTime:            sql.NullInt64{1559, true},
			AskSize:            sql.NullInt64{20, true},
			Bid:                sql.NullFloat64{0, true},
			BidTime:            sql.NullInt64{1559, true},
			BidSize:            sql.NullInt64{1559, true},
			Change:             sql.NullFloat64{4.1, true},
			Delta:              sql.NullFloat64{-0.25, true},
			EquityType:         core.EquityType("CALL"),
			Expiration:         core.NewDateTypeFromTime(time.Now()),
			Gamma:              sql.NullFloat64{-0.25, true},
			HighPrice:          sql.NullFloat64{0.06059, true},
			ImpliedVolatility:  sql.NullFloat64{4.3, true},
			LastTradePrice:     sql.NullFloat64{0.32521, true},
			LastTradeTimestamp: sql.NullInt64{1559, true},
			LastTradeVolume:    sql.NullInt64{9, true},
			LowPrice:           sql.NullFloat64{0.32521, true},
			OpenInterest:       sql.NullInt64{2, true},
			Rho:                sql.NullFloat64{285, true},
			StrikePrice:        15.0,
			Symbol:             core.StockSymbolType("TESTTEST"),
			Theta:              sql.NullFloat64{0, false},
			Timestamp:          core.NewTimestampTypeFromUnix(1231),
			Vega:               sql.NullFloat64{153.0, true},
		},
	}
	var sd *SibylDatabase
	var err error
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	if sd, err = ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress); err != nil {
		t.Errorf("cound not connect to sibyl database:%v", err)
		return
	}

	if err = sd.LoadOptionQuoteRecords(ctx, quotes); err != nil {
		t.Errorf("there was an error while adding the test quote, IT MAY REQUIRE MANUAL REMOVAL search for Symbol = 'TESTTEST'  : %v", err)
		return
	}
	execStr := fmt.Sprintf("DELETE FROM `%v`.`%v` where `stockSymbol` = 'TESTTEST';", SibylDatabaseName, OptionQuotesTableName)
	if _, err = sd.DBConn.ExecContext(ctx, execStr); err != nil {
		t.Errorf("there was an error while trying to remove the test quote, IT WILL REQUIRE MANUAL REMOVAL search for `stockSymbol` = 'TESTTEST'  : %v", err)
		return
	}
	sd.Close()
}

func TestSibylDatabase_LoadOptionQuoteRecords_largevolume(t *testing.T) {
	singleQuote := &core.SibylOptionQuoteRecord{
		Ask:                sql.NullFloat64{4.3, true},
		AskTime:            sql.NullInt64{1559, true},
		AskSize:            sql.NullInt64{20, true},
		Bid:                sql.NullFloat64{0, true},
		BidTime:            sql.NullInt64{1559, true},
		BidSize:            sql.NullInt64{1559, true},
		Change:             sql.NullFloat64{4.1, true},
		Delta:              sql.NullFloat64{-0.25, true},
		EquityType:         core.EquityType("CALL"),
		Expiration:         core.NewDateTypeFromTime(time.Now()),
		Gamma:              sql.NullFloat64{-0.25, true},
		HighPrice:          sql.NullFloat64{0.06059, true},
		ImpliedVolatility:  sql.NullFloat64{4.3, true},
		LastTradePrice:     sql.NullFloat64{0.32521, true},
		LastTradeTimestamp: sql.NullInt64{1559, true},
		LastTradeVolume:    sql.NullInt64{9, true},
		LowPrice:           sql.NullFloat64{0.32521, true},
		OpenInterest:       sql.NullInt64{2, true},
		Rho:                sql.NullFloat64{285, true},
		StrikePrice:        15.0,
		Symbol:             core.StockSymbolType("TESTTEST"),
		Theta:              sql.NullFloat64{0, false},
		Timestamp:          core.NewTimestampTypeFromUnix(1231),
		Vega:               sql.NullFloat64{153.0, true},
	}

	size := 60000
	quotes := make([]*core.SibylOptionQuoteRecord, 0, size)
	for i := 0; i < size; i++ {
		quotes = append(quotes, singleQuote)
	}

	var sd *SibylDatabase
	var err error
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	if sd, err = ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress); err != nil {
		t.Errorf("cound not connect to sibyl database:%v", err)
		return
	}

	if err = sd.LoadOptionQuoteRecords(ctx, quotes); err != nil {
		t.Errorf("there was an error while adding the test quote, IT MAY REQUIRE MANUAL REMOVAL search for Symbol = 'TESTTEST'  : %v", err)
		return
	}
	execStr := fmt.Sprintf("DELETE FROM `%v`.`%v` where `stockSymbol` = 'TESTTEST';", SibylDatabaseName, OptionQuotesTableName)
	if _, err = sd.DBConn.ExecContext(ctx, execStr); err != nil {
		t.Errorf("there was an error while trying to remove the test quote, IT WILL REQUIRE MANUAL REMOVAL search for `stockSymbol` = 'TESTTEST'  : %v", err)
		return
	}
	sd.Close()
}

func TestSibylDatabase_LoadIntradays(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	intraday := []*core.SibylIntradayRecord{{
		HighPrice: sql.NullFloat64{1539374400, true},
		LastPrice: sql.NullFloat64{142.07, true},
		LowPrice:  sql.NullFloat64{142.07, true},
		OpenPrice: sql.NullFloat64{142.07, true},
		Symbol:    core.StockSymbolType("TESTTEST"),
		Timestamp: core.NewTimestampTypeFromUnix(0),
		Volume:    sql.NullInt64{982778, true},
	}}

	var sd *SibylDatabase
	var err error
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if sd, err = ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress); err != nil {
		t.Errorf("cound not connect to sibyl database:%v", err)
		return
	}

	if err = sd.LoadIntradayRecords(ctx, intraday); err != nil {
		t.Errorf("there was an error while adding the test quote, IT MAY REQUIRE MANUAL REMOVAL search for `id` LIKE 'TESTTEST' : %v", err)
		return
	}

	execStr := fmt.Sprintf("DELETE FROM `%v`.`%v` where `id` LIKE 'TESTTEST%%';", SibylDatabaseName, IntradayTableName)
	if _, err = sd.DBConn.ExecContext(ctx, execStr); err != nil {
		t.Errorf("there was an error while trying to remove the test quote, IT WILL REQUIRE MANUAL REMOVAL search for `id` LIKE 'TESTTEST%%'  : %v", err)
		return
	}
	sd.Close()
}

func TestSibylDatabase_SetAndGetOptionsForStock(t *testing.T) {
	stockSymbol := core.StockSymbolType("TESTTEST")
	optionSymbols := []*core.OptionSymbolType{
		{Symbol: "TESTTEST", OptionType: "PUT", Expiration: core.NewDateType(2019, 1, 18), StrikePrice: 65},
		{Symbol: "TESTTEST", OptionType: "PUT", Expiration: core.NewDateType(2019, 1, 18), StrikePrice: 70},
		{Symbol: "TESTTEST", OptionType: "PUT", Expiration: core.NewDateType(2019, 1, 18), StrikePrice: 72},
		{Symbol: "TESTTEST", OptionType: "PUT", Expiration: core.NewDateType(2019, 1, 18), StrikePrice: 74.5},
		{Symbol: "TESTTEST", OptionType: "PUT", Expiration: core.NewDateType(2019, 1, 18), StrikePrice: 80},
	}

	var sd *SibylDatabase
	var err error
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if sd, err = ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress); err != nil {
		t.Errorf("cound not connect to sibyl database:%v", err)
		return
	}

	err = sd.SetOptionsForStock(ctx, stockSymbol, optionSymbols)
	if err != nil {
		t.Errorf("found an error %v", err)
		return
	}

	var options []*core.OptionSymbolType
	options, err = sd.GetOptionsFor(ctx, map[core.StockSymbolType]bool{stockSymbol: true})
	if err != nil {
		t.Errorf("expected no error but found %v", err)
	}

	if len(options) != len(optionSymbols) {
		t.Errorf("Expected %v options but found %v", len(optionSymbols), len(options))
	}
	_ = sd.SetOptionsForStock(ctx, stockSymbol, []*core.OptionSymbolType{})

	options, err = sd.GetOptionsFor(ctx, map[core.StockSymbolType]bool{stockSymbol: true})

	if len(options) != 0 {
		t.Errorf("Expected 0 options but found %v it will require manually deleting options with stock symbol 'TESTTEST'", len(options))
	}
}

func TestSibylDatabase_GetAllStocks(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	db, err := ConnectAndEnsureSibylDatabase(ctx, DefaultDatabaseServerAddress)
	if err != nil {
		t.Errorf("cound not connect to sibyl database:%v", err)
		return
	}
	//first we put in a test stock
	stockSymbol := core.StockSymbolType("TESTTEST")

	if err := db.StockAdd(ctx, stockSymbol); err != nil {
		t.Errorf("Expected to not have an issue with adding test stock")
		t.Errorf("MANUAL REMOVE MAY BE NEEDED")
		return
	}

	if stockRecords, err := db.GetAllStockRecords(ctx); err != nil {
		t.Errorf("Quoter: had a problem getting list of stocks: %v", err)
	} else {
		foundIt := false
		for _, stock := range stockRecords {
			if stock.Symbol == stockSymbol {
				foundIt = true
				break
			}
		}
		if !foundIt {
			t.Errorf("expected to find the stock symbol but didn't")
		}
	}

	if err := db.StockDelete(ctx, stockSymbol); err != nil {
		t.Errorf("Expected to not have an issue with adding test stock")
		t.Errorf("MANUAL REMOVE NEEDED")
	}
}
