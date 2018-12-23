package database

/*

When setting up the MYsql server consider setting the following values:

// this one helps speed up shutting down the server BIG deal when DB gets big
// see this site for more https://www.speedemy.com/how-to-speed-up-mysql-restart/
SET GLOBAL innodb_max_dirty_pages_pct = 0;

and these MUST be set:

SET GLOBAL local_infile = 'ON';
SHOW GLOBAL VARIABLES LIKE 'local_infile';
---------------------

For setting up the user the following commands in order can be used (or use the GUI):

CREATE USER 'sibyl'@'localhost' IDENTIFIED BY 'pa$$word';
GRANT CREATE ON *.* TO `sibyl`@'localhost';
GRANT DELETE ON *.* TO `sibyl`@'localhost';
GRANT DROP ON *.* TO `sibyl`@'localhost';
GRANT INDEX ON *.* TO `sibyl`@'localhost';
GRANT INSERT ON *.* TO `sibyl`@'localhost';
GRANT SELECT ON *.* TO `sibyl`@'localhost';
GRANT UPDATE ON *.* TO `sibyl`@'localhost';
FLUSH PRIVILEGES;

*/

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nathanhack/sibyl/agents/ally"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database/internal/scanners"
	"github.com/sirupsen/logrus"
)

const (
	DefaultDatabaseServerAddress = "localhost:3306"
	SibylDBUser                  = "sibyl"
	SibylDBUserPassword          = "pa$$word"
	SibylDatabaseName            = "sibyl"
	CredsTableName               = "creds"
	StocksTableName              = "stocks"
	OptionQuotesTableName        = "optionQuotes"
	StockQuotesTableName         = "stockQuotes"
	StableOptionQuotesTableName  = "stableOptionQuotes"
	StableStockQuotesTableName   = "stableStockQuotes"
	HistoryTableName             = "history"
	OptionsTableName             = "options"
	IntradayTableName            = "intraday"

	credsTableCreate = "CREATE TABLE IF NOT EXISTS `" + SibylDatabaseName + "`.`" + CredsTableName + "` (" +
		"`id` ENUM('1') NOT NULL, " +
		"`agentSelection` ENUM('none', 'ally_invest', 'td_ameritrade') NOT NULL DEFAULT 'none'," +
		"`customerKey` VARCHAR(255) NOT NULL DEFAULT '\"\"'," +
		"`customerSecret` VARCHAR(45) NOT NULL DEFAULT '\"\"'," +
		"`token` VARCHAR(45) NOT NULL DEFAULT '\"\"'," +
		"`tokenSecret` VARCHAR(45) NOT NULL DEFAULT '\"\"'," +
		"`urlRedirect` VARCHAR(255) NOT NULL DEFAULT '\"\"'," +
		"`accessToken` VARCHAR(1065) NOT NULL DEFAULT '\"\"'," +
		"`refreshToken` VARCHAR(1065) NOT NULL DEFAULT '\"\"'," +
		"`expireTimestamp` INT(18) NOT NULL DEFAULT 0," +
		"`refreshExpireTimestamp` INT(18) NOT NULL DEFAULT 0, " +
		"PRIMARY KEY(`id`));"

	optionQuotesTableCreate = "CREATE TABLE IF NOT EXISTS  `" + SibylDatabaseName + "`.`" + OptionQuotesTableName + "` ( " +
		"`id` VARCHAR(45) NOT NULL, " +
		"`ask` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`askTime` INT(18) NULL DEFAULT NULL, " +
		"`askSize` INT(18) NULL DEFAULT NULL, " +
		"`bid` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`bidTime` INT(18) NULL DEFAULT NULL, " +
		"`bidSize` INT(18) NULL DEFAULT NULL, " +
		"`change` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`delta` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`equityType` VARCHAR(5) NULL DEFAULT NULL, " +
		"`expiration` INT(18) NOT NULL, " +
		"`gamma` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`highPrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`impliedVolatility` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`lastTradePrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`lastTradeTimestamp` INT(18) NULL DEFAULT NULL, " +
		"`lastTradeVolume` DECIMAL(18) NULL DEFAULT NULL, " +
		"`lowPrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`openInterest` DECIMAL(18) NULL DEFAULT NULL, " +
		"`rho` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`strikePrice` DECIMAL(36,18) NOT NULL, " +
		"`symbol` VARCHAR(45) NOT NULL, " +
		"`theta` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`timestamp` INT(18) NOT NULL, " +
		"`vega` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC, `equityType` ASC, `timestamp` ASC) VISIBLE, " +
		"INDEX `index1` (`symbol` ASC, `timestamp` ASC, `equityType` ASC, `expiration` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	stockQuotesTableCreate = "CREATE TABLE IF NOT EXISTS  `" + SibylDatabaseName + "`.`" + StockQuotesTableName + "` ( " +
		"`id` VARCHAR(45) NOT NULL, " +
		"`ask` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`askTime` INT(18) NULL DEFAULT NULL, " +
		"`askSize` INT(18) NULL DEFAULT NULL, " +
		"`beta` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`bid` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`bidTime` INT(18) NULL DEFAULT NULL, " +
		"`bidSize` INT(18) NULL DEFAULT NULL, " +
		"`change` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`highPrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`lastTradePrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`lastTradeTimestamp` INT(18) NULL DEFAULT NULL, " +
		"`lastTradeVolume` DECIMAL(18) NULL DEFAULT NULL, " +
		"`lowPrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"`symbol` VARCHAR(45) NOT NULL, " +
		"`timestamp` INT(18) NOT NULL, " +
		"`volume` DECIMAL(18) NULL DEFAULT NULL, " +
		"`volWeightedAvgPrice` DECIMAL(36,18) NULL DEFAULT NULL, " +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC) VISIBLE, " +
		"INDEX `index1` (`symbol` ASC, `timestamp` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	stableOptionQuotesTableCreate = "CREATE TABLE IF NOT EXISTS  `" + SibylDatabaseName + "`.`" + StableOptionQuotesTableName + "` ( " +
		"`id` VARCHAR(45) NOT NULL, " +
		"`closePrice` DECIMAL(36,18) NULL default NULL," +
		"`contractSize` DECIMAL(18) NULL default NULL," +
		"`equityType` ENUM('CALL','PUT') NOT NULL," +
		"`expiration` INT(18) NOT NULL," +
		"`highPrice52Wk` DECIMAL(36,18) NULL default NULL," +
		"`highPrice52WkTimestamp` INT(18) NULL default NULL," +
		"`lowPrice52Wk` DECIMAL(36,18) NULL default NULL," +
		"`lowPrice52WkTimestamp` INT(18) NULL default NULL," +
		"`multiplier` INT(18) NULL default NULL," +
		"`openPrice` DECIMAL(36,18) NULL default NULL," +
		"`strikePrice` DECIMAL(36,18) NOT NULL," +
		"`symbol` VARCHAR(20) NOT NULL," +
		"`timestamp` INT(18) NOT NULL," +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC,`equityType` ASC, `timestamp` ASC) VISIBLE, " +
		"INDEX `index1` (`symbol` ASC,`equityType` ASC,`expiration` ASC, `strikePrice` ASC, `timestamp` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	stableStockQuotesTableCreate = "CREATE TABLE IF NOT EXISTS  `" + SibylDatabaseName + "`.`" + StableStockQuotesTableName + "` ( " +
		"`id` VARCHAR(45) NOT NULL, " +
		"`annualDividend` DECIMAL(36,18) NULL default NULL," +
		"`bookValue` DECIMAL(36,18) NULL default NULL," +
		"`closePrice` DECIMAL(36,18) NULL default NULL," +
		"`div` DECIMAL(36,18) NULL default NULL," +
		"`divExTimestamp` INT(18) NULL default NULL," +
		"`divFreq` ENUM('A','S','Q','M','N') NULL default NULL," +
		"`divPayTimestamp` INT(18) NULL default NULL," +
		"`eps` DECIMAL(36,18) NULL default NULL," +
		"`highPrice52Wk` DECIMAL(36,18) NULL default NULL," +
		"`highPrice52WkTimestamp` INT(18) NULL default NULL," +
		"`lowPrice52Wk` DECIMAL(36,18) NULL default NULL," +
		"`lowPrice52WkTimestamp` INT(18) NULL default NULL," +
		"`openPrice` DECIMAL(36,18) NULL default NULL," +
		"`priceEarnings` DECIMAL(36,18) NULL default NULL," +
		"`sharesOutstanding` DECIMAL(18) NULL default NULL," +
		"`symbol` VARCHAR(20) NOT NULL ," +
		"`timestamp` INT(18) NOT NULL," +
		"`volatility` DECIMAL(36,18) NULL default NULL," +
		"`yield` DECIMAL(36,18) NULL default NULL," +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC, `timestamp` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	historyTableCreate = "CREATE TABLE IF NOT EXISTS `" + SibylDatabaseName + "`.`" + HistoryTableName + "` (" +
		"`id` VARCHAR(45) NOT NULL, " +
		"`closePrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`highPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`lowPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`openPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`symbol` VARCHAR(45) NOT NULL ," +
		"`timestamp` INT(18) NOT NULL ," +
		"`volume` INT(18) NULL DEFAULT NULL," +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id`ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC) VISIBLE,INDEX `index1` (`symbol` ASC, `timestamp` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	intradayTableCreate = "CREATE TABLE IF NOT EXISTS `" + SibylDatabaseName + "`.`" + IntradayTableName + "` (" +
		"`id` VARCHAR(45) NOT NULL, " +
		"`highPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`lastPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`lowPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`openPrice` DECIMAL(36,18) NULL DEFAULT NULL," +
		"`symbol` VARCHAR(45) NOT NULL," +
		"`timestamp` INT(18) NOT NULL," +
		"`volume` INT(18) NULL DEFAULT NULL," +
		"PRIMARY KEY(`id`),UNIQUE INDEX `id_UNIQUE` (`id`ASC) VISIBLE, " +
		"INDEX `index0` (`symbol` ASC) VISIBLE,INDEX " +
		"`index1` (`symbol` ASC, `timestamp` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	stockOfInterestCreate = "CREATE TABLE IF NOT EXISTS `" + SibylDatabaseName + "`.`" + StocksTableName + "` (" +
		"`id` INT(18) NOT NULL AUTO_INCREMENT, " +
		"`downloadStatus` ENUM('enabled','disabled') NOT NULL DEFAULT 'disabled'," +
		"`exchange` VARCHAR(15) NOT NULL DEFAULT '\"\"'," +
		"`exchangeDescription` VARCHAR(60) NOT NULL DEFAULT '\"\"'," +
		"`hasOptions` ENUM('yes','no') NOT NULL DEFAULT 'no'," +
		"`historyStatus` ENUM('enabled','disabled') NOT NULL DEFAULT 'disabled'," +
		"`intradayHistoryStatus` ENUM('enabled','disabled') NOT NULL DEFAULT 'disabled'," +
		"`name` VARCHAR(100) NOT NULL DEFAULT '\"\"'," +
		"`quotesStatus` ENUM('enabled','disabled') NOT NULL DEFAULT 'disabled'," +
		"`stableQuotesStatus` ENUM('enabled','disabled') NOT NULL DEFAULT 'disabled'," +
		"`symbol` VARCHAR(45) NOT NULL," +
		"`validationStatus` ENUM('pending','valid','invalid') NOT NULL DEFAULT 'pending'," +
		"PRIMARY KEY(`id`),UNIQUE INDEX `symbol_UNIQUE` (`symbol` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"

	optionsTableCreate = "CREATE TABLE IF NOT EXISTS `" + SibylDatabaseName + "`.`" + OptionsTableName + "` (" +
		"`id` INT(18) NOT NULL AUTO_INCREMENT," +
		"`expiration` INT(18) NOT NULL," +
		"`optionType` ENUM('CALL', 'PUT') NOT NULL," +
		"`strikePrice` DECIMAL(36,18) NOT NULL," +
		"`symbol` VARCHAR(45) NOT NULL," +
		"PRIMARY KEY (`id`), INDEX `index0` (`symbol` ASC) VISIBLE) ROW_FORMAT = COMPRESSED;"
)

type SibylDatabase struct {
	DBConn *sql.DB
}

func ConnectAndEnsureSibylDatabase(ctx context.Context, address string) (*SibylDatabase, error) {
	toReturn := SibylDatabase{}

	if err := toReturn.connect(SibylDBUser, SibylDBUserPassword, address); err != nil {
		return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while opening connection: %v \n"+
			"If this is the first connecting to the database make sure the user:sibyl exists with privileges: CREATE, DELETE, DROP, INDEX, INSERT, SELECT, UPDATE."+
			"Additionally, ensure system variable 'local_infile' is 'ON'.", err)
	}

	if err := toReturn.verifySQLConnection(ctx); err != nil {
		return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while verifing connection: %v", err)
	}

	//now that we're connected to the database backend
	//time to verify and/or create the database and tables needed

	if !toReturn.dbExists(ctx, SibylDatabaseName) {
		if err := toReturn.createDatabase(ctx, SibylDatabaseName); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v database: %v", SibylDatabaseName, err)
		}
	}

	//TODO consider change out to use EnsureTableExists()
	if !toReturn.hasTable(ctx, SibylDatabaseName, CredsTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, credsTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table: %v", CredsTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, StocksTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, stockOfInterestCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table: %v", StocksTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, OptionQuotesTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, optionQuotesTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table : %v", OptionQuotesTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, StockQuotesTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, stockQuotesTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table : %v", StockQuotesTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, StableOptionQuotesTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, stableOptionQuotesTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table : %v", StableOptionQuotesTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, StableStockQuotesTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, stableStockQuotesTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table : %v", StableStockQuotesTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, HistoryTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, historyTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table: %v", HistoryTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, IntradayTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, intradayTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table: %v", IntradayTableName, err)
		}
	}

	if !toReturn.hasTable(ctx, SibylDatabaseName, OptionsTableName) {
		if _, err := toReturn.DBConn.ExecContext(ctx, optionsTableCreate); err != nil {
			return nil, fmt.Errorf("ConnectAndEnsureSibylDatabase: found an error while creating %v table: %v", OptionsTableName, err)
		}
	}

	return &toReturn, nil
}

type DatabaseStringer interface {
	StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string
}

func (sd *SibylDatabase) verifySQLConnection(ctx context.Context) error {

	if err := sd.DBConn.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

type LoadDupAction string

const (
	Replace  LoadDupAction = "REPLACE"
	Ignore   LoadDupAction = "IGNORE"
	NoAction LoadDupAction = ""
)

func inList(item string, list []string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}

func (sd *SibylDatabase) loadRecords(ctx context.Context, records []DatabaseStringer, databaseName, tableName string, combineIntoID []string, recordFieldNames []string, action LoadDupAction) error {
	if len(records) == 0 {
		//there's no error for not passing in anything to insert
		return nil
	}
	if len(recordFieldNames) == 0 {
		return fmt.Errorf("loadRecords: no field names passed in, requires at least one")
	}

	recordStrBuilder := strings.Builder{}
	for _, record := range records {
		//IMPORTANT NOTE : we use the \N to denote NULL field values
		recordStrBuilder.WriteString(record.StringBlindWithDelimiter(";", "\\N", false))
		recordStrBuilder.WriteString("\n")
	}
	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.Debugf("loadRecords: %v", recordStrBuilder.String())
	}
	buf := bytes.NewBufferString(recordStrBuilder.String())

	// create and assign the Reader
	mysql.RegisterReaderHandler("test", func() io.Reader {
		return buf
	})

	variables := make(map[string]string)
	for i, name := range combineIntoID {
		variables[name] = fmt.Sprintf("var%v", i)
	}

	insertCommandBuilder := strings.Builder{}
	insertCommandBuilder.WriteString("LOAD DATA LOCAL INFILE 'Reader::test' ")
	insertCommandBuilder.WriteString(string(action))
	insertCommandBuilder.WriteString(" INTO TABLE `")
	insertCommandBuilder.WriteString(databaseName)
	insertCommandBuilder.WriteString("`.`")
	insertCommandBuilder.WriteString(tableName)
	insertCommandBuilder.WriteString("` FIELDS TERMINATED BY ';' (")
	if inList(recordFieldNames[0], combineIntoID) {
		insertCommandBuilder.WriteString("@")
		insertCommandBuilder.WriteString(variables[recordFieldNames[0]])
	} else {
		insertCommandBuilder.WriteString("`")
		insertCommandBuilder.WriteString(recordFieldNames[0])
		insertCommandBuilder.WriteString("`")
	}
	for _, fieldName := range recordFieldNames[1:] {
		if inList(fieldName, combineIntoID) {
			insertCommandBuilder.WriteString(",@")
			insertCommandBuilder.WriteString(variables[fieldName])
		} else {
			insertCommandBuilder.WriteString(",`")
			insertCommandBuilder.WriteString(fieldName)
			insertCommandBuilder.WriteString("`")
		}
	}
	insertCommandBuilder.WriteString(")")

	if len(combineIntoID) > 0 {
		// we take these and concat into 'id'
		insertCommandBuilder.WriteString(" SET `id` = concat(")
		insertCommandBuilder.WriteString("@")
		insertCommandBuilder.WriteString(variables[combineIntoID[0]])
		for _, name := range combineIntoID[1:] {
			insertCommandBuilder.WriteString(",@")
			insertCommandBuilder.WriteString(variables[name])
		}
		insertCommandBuilder.WriteString(")")

		for _, name := range combineIntoID {
			insertCommandBuilder.WriteString(", `")
			insertCommandBuilder.WriteString(name)
			insertCommandBuilder.WriteString("`= @")
			insertCommandBuilder.WriteString(variables[name])
		}
	}
	insertCommandBuilder.WriteString(";")
	_, err := sd.DBConn.ExecContext(ctx, insertCommandBuilder.String())

	mysql.DeregisterReaderHandler("test")
	if err != nil {
		return fmt.Errorf("loadRecords: error during db exec insert exec make sure local_infle=ON,s [%v]: %v", insertCommandBuilder.String(), err)
	}

	return nil
}

func (sd *SibylDatabase) loadFile(ctx context.Context, filePathname string, databaseName, tableName string, combineIntoID []string, recordFieldNames []string, action LoadDupAction) error {
	if _, err := os.Stat(filePathname); os.IsNotExist(err) {
		return fmt.Errorf("loadFile: file must exist")
	}
	if len(recordFieldNames) == 0 {
		return fmt.Errorf("loadFile: no field names passed in, requires at least one")
	}

	mysql.RegisterLocalFile(filePathname)

	variables := make(map[string]string)
	for i, name := range combineIntoID {
		variables[name] = fmt.Sprintf("var%v", i)
	}

	insertCommandBuilder := strings.Builder{}
	insertCommandBuilder.WriteString("LOAD DATA LOCAL INFILE '")
	insertCommandBuilder.WriteString(filePathname)
	insertCommandBuilder.WriteString("' ")
	insertCommandBuilder.WriteString(string(action))
	insertCommandBuilder.WriteString(" INTO TABLE `")
	insertCommandBuilder.WriteString(databaseName)
	insertCommandBuilder.WriteString("`.`")
	insertCommandBuilder.WriteString(tableName)
	insertCommandBuilder.WriteString("` FIELDS TERMINATED BY ';' (")
	if inList(recordFieldNames[0], combineIntoID) {
		insertCommandBuilder.WriteString("@")
		insertCommandBuilder.WriteString(variables[recordFieldNames[0]])
	} else {
		insertCommandBuilder.WriteString("`")
		insertCommandBuilder.WriteString(recordFieldNames[0])
		insertCommandBuilder.WriteString("`")
	}
	for _, fieldName := range recordFieldNames[1:] {
		if inList(fieldName, combineIntoID) {
			insertCommandBuilder.WriteString(",@")
			insertCommandBuilder.WriteString(variables[fieldName])
		} else {
			insertCommandBuilder.WriteString(",`")
			insertCommandBuilder.WriteString(fieldName)
			insertCommandBuilder.WriteString("`")
		}
	}
	insertCommandBuilder.WriteString(")")

	if len(combineIntoID) > 0 {
		// we take these and concat into 'id'
		insertCommandBuilder.WriteString(" SET `id` = concat(")
		insertCommandBuilder.WriteString("@")
		insertCommandBuilder.WriteString(variables[combineIntoID[0]])
		for _, name := range combineIntoID[1:] {
			insertCommandBuilder.WriteString(",@")
			insertCommandBuilder.WriteString(variables[name])
		}
		insertCommandBuilder.WriteString(")")

		for _, name := range combineIntoID {
			insertCommandBuilder.WriteString(", `")
			insertCommandBuilder.WriteString(name)
			insertCommandBuilder.WriteString("`= @")
			insertCommandBuilder.WriteString(variables[name])
		}
	}
	insertCommandBuilder.WriteString(";")
	_, err := sd.DBConn.ExecContext(ctx, insertCommandBuilder.String())

	mysql.DeregisterLocalFile(filePathname)
	if err != nil {
		return fmt.Errorf("loadFile: error during db exec: %v", err)
	}
	return nil
}

func (sd *SibylDatabase) LoadCredsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()
	mysql.RegisterLocalFile(filePathname)

	var err error
	_, err = sd.DBConn.ExecContext(ctx, "LOAD DATA LOCAL INFILE '"+filePathname+"'  INTO TABLE `"+SibylDatabaseName+"`.`"+CredsTableName+"` FIELDS TERMINATED BY ';';")
	mysql.DeregisterLocalFile(filePathname)
	if err != nil {
		return fmt.Errorf("LoadCredsFromFile: error during db exec: %v", err)
	}

	logrus.Debugf("LoadCredsFromFile: Data saved to %v in %s", CredsTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadCreds(ctx context.Context, creds *core.SibylCreds) error {
	updateStr := fmt.Sprintf("REPLACE INTO `%v`.`%v` VALUES ('1', %v) ;",
		SibylDatabaseName,
		CredsTableName,
		creds.StringBlindWithDelimiter(",", "", true),
	)
	if _, err := sd.DBConn.ExecContext(ctx, updateStr); err != nil {
		return fmt.Errorf("LoadCreds: error while executing update: %v", err)
	}
	return nil
}

func (sd *SibylDatabase) GetCreds(ctx context.Context) (*core.SibylCreds, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v`", SibylDatabaseName, CredsTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("GetCreds: had error running query: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		if creds, err := scanners.ScanSibylCredsRow(rows); err != nil {
			return nil, fmt.Errorf("GetCreds: had error reading results: %v", err)
		} else {
			return creds, nil
		}
	}
	return nil, fmt.Errorf("GetCreds: no values")
}

func (sd *SibylDatabase) LoadStockQuoteRecords(ctx context.Context, quotes []*core.SibylStockQuoteRecord) error {
	startTime := time.Now()
	if len(quotes) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(quotes))
	for i, q := range quotes {
		records[i] = q
	}

	err := sd.loadRecords(ctx, records, SibylDatabaseName, StockQuotesTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"ask",
			"askTime",
			"askSize",
			"beta",
			"bid",
			"bidTime",
			"bidSize",
			"change",
			"highPrice",
			"lastTradePrice",
			"lastTradeTimestamp",
			"lastTradeVolume",
			"lowPrice",
			"symbol",
			"timestamp",
			"volume",
			"volWeightedAvgPrice",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadStockQuoteRecords: had an error while inserting records: %v", err)
	}

	logrus.Debugf("LoadStockQuoteRecords: Data saved to %v in %s", StockQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadStockQuoteRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()

	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, StockQuotesTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"ask",
			"askTime",
			"askSize",
			"beta",
			"bid",
			"bidTime",
			"bidSize",
			"change",
			"highPrice",
			"lastTradePrice",
			"lastTradeTimestamp",
			"lastTradeVolume",
			"lowPrice",
			"symbol",
			"timestamp",
			"volume",
			"volWeightedAvgPrice",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadStockQuoteRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadStockQuoteRecordsFromFile: Data saved to %v in %s", StockQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadOptionQuoteRecords(ctx context.Context, quotes []*core.SibylOptionQuoteRecord) error {
	startTime := time.Now()
	if len(quotes) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(quotes))
	for i, q := range quotes {
		records[i] = q
	}

	err := sd.loadRecords(ctx, records, SibylDatabaseName, OptionQuotesTableName,
		[]string{
			"symbol",
			"expiration",
			"equityType",
			"strikePrice",
			"timestamp",
		},
		[]string{
			"ask",
			"askTime",
			"askSize",
			"bid",
			"bidTime",
			"bidSize",
			"change",
			"delta",
			"equityType",
			"expiration",
			"gamma",
			"highPrice",
			"impliedVolatility",
			"lastTradePrice",
			"lastTradeTimestamp",
			"lastTradeVolume",
			"lowPrice",
			"openInterest",
			"rho",
			"strikePrice",
			"symbol",
			"theta",
			"timestamp",
			"vega",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadOptionQuoteRecords: had an error while inserting records: %v", err)
	}

	logrus.Debugf("LoadOptionQuoteRecords: Data saved to %v in %s", OptionQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadOptionQuoteRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()

	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, OptionQuotesTableName,
		[]string{
			"symbol",
			"expiration",
			"equityType",
			"strikePrice",
			"timestamp",
		},
		[]string{
			"ask",
			"askTime",
			"askSize",
			"bid",
			"bidTime",
			"bidSize",
			"change",
			"delta",
			"equityType",
			"expiration",
			"gamma",
			"highPrice",
			"impliedVolatility",
			"lastTradePrice",
			"lastTradeTimestamp",
			"lastTradeVolume",
			"lowPrice",
			"openInterest",
			"rho",
			"strikePrice",
			"symbol",
			"theta",
			"timestamp",
			"vega",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadOptionQuoteRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadOptionQuoteRecordsFromFile: Data saved to %v in %s", OptionQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadStableOptionQuoteRecords(ctx context.Context, quotes []*core.SibylStableOptionQuoteRecord) error {
	startTime := time.Now()
	if len(quotes) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(quotes))
	for i, q := range quotes {
		records[i] = q
	}
	err := sd.loadRecords(ctx, records, SibylDatabaseName, StableOptionQuotesTableName,
		[]string{
			"symbol",
			"expiration",
			"equityType",
			"strikePrice",
			"timestamp",
		},
		[]string{
			"closePrice",
			"contractSize",
			"equityType",
			"expiration",
			"highPrice52Wk",
			"highPrice52WkTimestamp",
			"lowPrice52Wk",
			"lowPrice52WkTimestamp",
			"multiplier",
			"openPrice",
			"strikePrice",
			"symbol",
			"timestamp",
		}, Replace)

	if err != nil {
		return fmt.Errorf("LoadStableOptionQuoteRecords: error while inserting records: %v", err)
	}

	logrus.Debugf("LoadStableOptionQuoteRecords: Data saved to %v in %s", StableOptionQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadStableOptionQuoteRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()
	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, StableOptionQuotesTableName,
		[]string{
			"symbol",
			"expiration",
			"equityType",
			"strikePrice",
			"timestamp",
		},
		[]string{
			"closePrice",
			"contractSize",
			"equityType",
			"expiration",
			"highPrice52Wk",
			"highPrice52WkTimestamp",
			"lowPrice52Wk",
			"lowPrice52WkTimestamp",
			"multiplier",
			"openPrice",
			"strikePrice",
			"symbol",
			"timestamp",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadStableOptionQuoteRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadStableOptionQuoteRecordsFromFile: Data saved to %v in %s", StableOptionQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadStableStockQuoteRecords(ctx context.Context, quotes []*core.SibylStableStockQuoteRecord) error {
	startTime := time.Now()
	if len(quotes) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(quotes))
	for i, q := range quotes {
		records[i] = q
	}

	err := sd.loadRecords(ctx, records, SibylDatabaseName, StableStockQuotesTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"annualDividend",
			"bookValue",
			"closePrice",
			"div",
			"divExTimestamp",
			"divFreq",
			"divPayTimestamp",
			"eps",
			"highPrice52Wk",
			"highPrice52WkTimestamp",
			"lowPrice52Wk",
			"lowPrice52WkTimestamp",
			"openPrice",
			"priceEarnings",
			"sharesOutstanding",
			"symbol",
			"timestamp",
			"volatility",
			"yield",
		}, Replace)

	if err != nil {
		return fmt.Errorf("LoadStableStockQuoteRecords: error while inserting records: %v", err)
	}

	logrus.Debugf("LoadStableStockQuoteRecords: Data saved to %v in %s", StableStockQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadStableStockQuoteRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()
	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, StableStockQuotesTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"annualDividend",
			"bookValue",
			"closePrice",
			"div",
			"divExTimestamp",
			"divFreq",
			"divPayTimestamp",
			"eps",
			"highPrice52Wk",
			"highPrice52WkTimestamp",
			"lowPrice52Wk",
			"lowPrice52WkTimestamp",
			"openPrice",
			"priceEarnings",
			"sharesOutstanding",
			"symbol",
			"timestamp",
			"volatility",
			"yield",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadStableStockQuoteRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadStableStockQuoteRecordsFromFile: Data saved to %v in %s", StableStockQuotesTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadHistoryRecords(ctx context.Context, histories []*core.SibylHistoryRecord) error {
	startTime := time.Now()
	if len(histories) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(histories))
	for i, q := range histories {
		records[i] = q
	}

	err := sd.loadRecords(ctx, records, SibylDatabaseName, HistoryTableName,
		[]string{"symbol", "timeStamp"},
		[]string{
			"closePrice",
			"highPrice",
			"lowPrice",
			"openPrice",
			"symbol",
			"timeStamp",
			"volume",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadHistoryRecords: had an error while inserting records: %v", err)
	}

	logrus.Debugf("LoadHistoryRecords: Data saved to %v in %s", HistoryTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadHistoryRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()

	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, HistoryTableName,
		[]string{"symbol", "timeStamp"},
		[]string{
			"closePrice",
			"highPrice",
			"lowPrice",
			"openPrice",
			"symbol",
			"timeStamp",
			"volume",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadHistoryRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadHistoryRecordsFromFile: Data saved to %v in %s", HistoryTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadIntradayRecords(ctx context.Context, intradays []*core.SibylIntradayRecord) error {
	startTime := time.Now()
	if len(intradays) == 0 {
		return nil
	}

	records := make([]DatabaseStringer, len(intradays))
	for i, q := range intradays {
		records[i] = q
	}

	err := sd.loadRecords(ctx, records, SibylDatabaseName, IntradayTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"highPrice",
			"lastPrice",
			"lowPrice",
			"openPrice",
			"symbol",
			"timestamp",
			"volume",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadIntradayRecords: had an error while inserting records: %v", err)
	}

	logrus.Debugf("LoadIntradayRecords: Data saved to %v in %s", IntradayTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) LoadIntradayRecordsFromFile(ctx context.Context, filePathname string) error {
	//this assumes the file was dumped by this struct's DumpToFile function
	startTime := time.Now()

	err := sd.loadFile(ctx, filePathname, SibylDatabaseName, IntradayTableName,
		[]string{"symbol", "timestamp"},
		[]string{
			"highPrice",
			"lastPrice",
			"lowPrice",
			"openPrice",
			"symbol",
			"timestamp",
			"volume",
		}, NoAction)

	if err != nil {
		return fmt.Errorf("LoadHistoryRecordsFromFile: had error while inserting into database: %v", err)
	}

	logrus.Debugf("LoadHistoryRecordsFromFile: Data saved to %v in %s", HistoryTableName, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) DumpCredsToFile(ctx context.Context, filePathname string) error {
	startTime := time.Now()

	if _, err := os.Stat(filePathname); os.IsExist(err) {
		return fmt.Errorf("DumpCredsToFile: had an error the file must already exist")
	}

	creds, err := sd.GetCreds(ctx)
	if err != nil {
		return fmt.Errorf("DumpCredsToFile: had an error while getting creds: %v", err)
	}

	file, err := os.Create(filePathname)
	if err != nil {
		return fmt.Errorf("DumpCredsToFile: could not create file: %v had error: %v", filePathname, err)
	}

	buf := bufio.NewWriter(file)
	stringToWrite := ";" + creds.StringBlindWithDelimiter(";", "\\N", false)

	count, err := buf.WriteString(stringToWrite + "\n")
	if err != nil {
		file.Close()
		os.Remove(filePathname)
		logrus.Errorf("DumpCredsToFile: failed to write out %v with error: %v", stringToWrite, err)
		return err
	}
	if count != len(stringToWrite)+1 {
		file.Close()
		os.Remove(filePathname)
		logrus.Errorf("DumpCredsToFile: failed to write out the expected number of bytes, expected %v found %v", len(stringToWrite), count)
		return err
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpCredsToFile: wrote creds to %v in %s", filePathname, time.Since(startTime))
	return nil
}

func (sd *SibylDatabase) DumpOptionQuoteRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeOptionQuoteRecordsToFile(ctx, filePathname, "", -1)
	return err
}

func (sd *SibylDatabase) DumpRangeOptionQuoteRecordsToFile(ctx context.Context, filePathname string, lastID string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, OptionQuotesTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id` > '%v' limit %v;", SibylDatabaseName, OptionQuotesTableName, lastID, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeOptionQuoteRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var quote *core.SibylOptionQuoteRecord
	for rows.Next() {
		nextLastID, quote, err = scanners.ScanSibylOptionQuoteRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeOptionQuoteRecordsToFile: failed to scan quote %v: %v", quote, err)
		}
		//IMPORTANT NOTE : we use the \N to denote NULL field values
		quoteStr := quote.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeOptionQuoteRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeOptionQuoteRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeOptionQuoteRecordsToFile: dumped all(%v) quotes to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}

func (sd *SibylDatabase) DumpStockQuoteRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeStockQuoteRecordsToFile(ctx, filePathname, "", -1)
	return err
}

func (sd *SibylDatabase) DumpRangeStockQuoteRecordsToFile(ctx context.Context, filePathname string, lastID string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, StockQuotesTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id` > '%v' limit %v;", SibylDatabaseName, StockQuotesTableName, lastID, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeStockQuoteRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var quote *core.SibylStockQuoteRecord
	for rows.Next() {
		nextLastID, quote, err = scanners.ScanSibylStockQuoteRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStockQuoteRecordsToFile: failed to scan quote %v: %v", quote, err)
		}
		//IMPORTANT NOTE : we use the \N to denote NULL field values
		quoteStr := quote.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStockQuoteRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStockQuoteRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeStockQuoteRecordsToFile: dumped all(%v) quotes to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}

func (sd *SibylDatabase) DumpStableOptionQuoteRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeStableStockQuoteRecordsToFile(ctx, filePathname, "", -1)
	return err
}

func (sd *SibylDatabase) DumpRangeStableOptionQuoteRecordsToFile(ctx context.Context, filePathname string, lastID string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, StableOptionQuotesTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id` > '%v' limit %v;", SibylDatabaseName, StableOptionQuotesTableName, lastID, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeStableOptionQuoteRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var stableQuote *core.SibylStableOptionQuoteRecord
	for rows.Next() {
		nextLastID, stableQuote, err = scanners.ScanSibylStableOptionQuoteRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableOptionQuoteRecordsToFile: failed to scan quote %v: %v", stableQuote, err)
		}
		//IMPORTANT NOTE : we use the \N to denote NULL field values
		quoteStr := stableQuote.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableOptionQuoteRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableOptionQuoteRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeStableOptionQuoteRecordsToFile: dumped all(%v) quotes to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}

func (sd *SibylDatabase) DumpStableStockQuoteRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeStableStockQuoteRecordsToFile(ctx, filePathname, "", -1)
	return err
}

func (sd *SibylDatabase) DumpRangeStableStockQuoteRecordsToFile(ctx context.Context, filePathname string, lastID string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, StableStockQuotesTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id` > '%v' limit %v;", SibylDatabaseName, StableStockQuotesTableName, lastID, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeStableStockQuoteRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var stableQuote *core.SibylStableStockQuoteRecord
	for rows.Next() {
		nextLastID, stableQuote, err = scanners.ScanSibylStableStockQuoteRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableStockQuoteRecordsToFile: failed to scan quote %v: %v", stableQuote, err)
		}
		//IMPORTANT NOTE : we use the \N to denote NULL field values
		quoteStr := stableQuote.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableStockQuoteRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeStableStockQuoteRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeStableStockQuoteRecordsToFile: dumped all(%v) quotes to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}

func (sd *SibylDatabase) DumpIntradayRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeIntradayRecordsToFile(ctx, filePathname, "", -1)
	return err
}

func (sd *SibylDatabase) DumpRangeIntradayRecordsToFile(ctx context.Context, filePathname string, lastId string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, IntradayTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id`> '%v' limit %v;", SibylDatabaseName, IntradayTableName, lastId, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("DumpRangeIntradayRecordsToFile: had an error: %v", err)
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeIntradayRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var intraday *core.SibylIntradayRecord
	for rows.Next() {
		nextLastID, intraday, err = scanners.ScanSibylIntradayRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeIntradayRecordsToFile: failed to scan quote %v: %v", intraday, err)
		}

		quoteStr := intraday.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeIntradayRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeIntradayRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeIntradayRecordsToFile: dumped all(%v) intradays to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}
func (sd *SibylDatabase) DumpHistoryRecordsToFile(ctx context.Context, filePathname string) error {
	_, err := sd.DumpRangeHistoryRecordsToFile(ctx, filePathname, "", -1)
	return err
}
func (sd *SibylDatabase) DumpRangeHistoryRecordsToFile(ctx context.Context, filePathname string, lastID string, count int) (string, error) {
	startTime := time.Now()
	var queryStr string
	if count < 0 {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v`;", SibylDatabaseName, HistoryTableName)
	} else {
		queryStr = fmt.Sprintf("SELECT * FROM `%v`.`%v` where `id` > '%v' limit %v;", SibylDatabaseName, HistoryTableName, lastID, count)
	}

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("DumpRangeHistoryRecordsToFile: had an error: %v", err)
	}
	defer rows.Close()

	file, err := os.Create(filePathname)
	if err != nil {
		return "", fmt.Errorf("DumpRangeHistoryRecordsToFile: could not create file: %v had error: %v", filePathname, err)
	}
	buf := bufio.NewWriter(file)
	rowCount := 0
	var nextLastID string
	var history *core.SibylHistoryRecord
	for rows.Next() {
		nextLastID, history, err = scanners.ScanSibylHistoryRecordRow(rows)
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeHistoryRecordsToFile: failed to scan quote %v: %v", history, err)
		}

		quoteStr := history.StringBlindWithDelimiter(";", "\\N", false)
		count, err := buf.WriteString(quoteStr + "\n")
		if err != nil {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeHistoryRecordsToFile: failed to write out %v with error: %v", quoteStr, err)
		}
		if count != len(quoteStr)+1 {
			file.Close()
			os.Remove(filePathname)
			return "", fmt.Errorf("DumpRangeHistoryRecordsToFile: failed to write out the expected number of bytes, expected %v found %v", len(quoteStr), count)
		}
		rowCount++
		if rowCount%10000 == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
	file.Close()
	logrus.Infof("DumpRangeHistoryRecordsToFile: dumped all(%v) histories to %v in %s", rowCount, filePathname, time.Since(startTime))
	return nextLastID, nil
}

func (gtd *SibylDatabase) connect(username, password string, address string) error {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/", username, password, address)
	var err error
	gtd.DBConn, err = sql.Open("mysql", dataSourceName)

	gtd.DBConn.SetMaxOpenConns(18)
	gtd.DBConn.Ping()
	return err
}

func (gtd *SibylDatabase) dbExists(ctx context.Context, databaseName string) bool {
	row := gtd.DBConn.QueryRowContext(ctx, "SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?;", databaseName)
	var count int
	if err := row.Scan(&count); err == nil {
		if count > 0 {
			return true
		}
	}
	return false
}

func (gtd *SibylDatabase) createDatabase(ctx context.Context, databaseName string) error {
	execStr := fmt.Sprintf("CREATE DATABASE `%v`;", databaseName)
	_, err := gtd.DBConn.ExecContext(ctx, execStr)
	return err
}

func (gtd *SibylDatabase) hasTable(ctx context.Context, databaseName, tableName string) bool {
	if gtd.dbExists(ctx, databaseName) {
		rows, err := gtd.DBConn.QueryContext(ctx, "SELECT * FROM information_schema.tables	WHERE table_schema = ? AND table_name = ? LIMIT 1;", databaseName, tableName)
		if err == nil {
			defer rows.Close()
			return rows.Next()
		}
	}
	return false
}

func (sd *SibylDatabase) GetAllStockRecords(ctx context.Context) ([]*core.SibylStockRecord, error) {
	startTime := time.Now()
	toReturn := make([]*core.SibylStockRecord, 0)
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v`", SibylDatabaseName, StocksTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return toReturn, err
	}
	defer rows.Close()

	errStrings := make([]string, 0)
	for rows.Next() {
		stock, err := scanners.ScanSibylStockRecordRow(rows)
		if err != nil {
			errStrings = append(errStrings, err.Error())
		} else {
			toReturn = append(toReturn, stock)
		}
	}

	logrus.Debugf("Found %v stocks in %s", len(toReturn), time.Since(startTime))

	if len(errStrings) != 0 {
		return toReturn, fmt.Errorf("GetAllStocks: had a error while getting stocks: %v", strings.Join(errStrings, ";"))
	}
	return toReturn, nil
}
func (sd *SibylDatabase) StockAdd(ctx context.Context, symbol core.StockSymbolType) error {
	insertStr := fmt.Sprintf("INSERT IGNORE into `%v`.`%v` (`symbol`) values ( '%v');", SibylDatabaseName, StocksTableName, strings.ToUpper(string(symbol)))

	if _, err := sd.DBConn.ExecContext(ctx, insertStr); err != nil {
		return fmt.Errorf("StockAdd: failed to add the stock had the following error: %v", err)
	}
	return nil
}

func (sd *SibylDatabase) stockOneElementChange(ctx context.Context, symbol core.StockSymbolType, fieldName string, status string) error {
	insertStr := fmt.Sprintf("UPDATE `%v`.`%v` SET `%v` ='%v' WHERE (`symbol` = '%v');", SibylDatabaseName, StocksTableName, fieldName, status, symbol)

	if _, err := sd.DBConn.ExecContext(ctx, insertStr); err != nil {
		return fmt.Errorf("stockOneElementChange: failed to update the stock for %v -> %v had the following error: %v", fieldName, status, err)
	}
	return nil
}

func (sd *SibylDatabase) StockDelete(ctx context.Context, symbol core.StockSymbolType) error {
	//In order to delete a stock it must be deleted from multiple tables basically all of
	// them but we'll start with the smallest tables first and work our way up
	errStrings := make([]string, 0)

	//we'll start with the stocks table
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, StocksTableName) {
		errStrings = append(errStrings, err.Error())
	}

	//next the stablequotes tables
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, StableStockQuotesTableName) {
		errStrings = append(errStrings, err.Error())
	}
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, StableOptionQuotesTableName) {
		errStrings = append(errStrings, err.Error())
	}

	//next the history table
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, HistoryTableName) {
		errStrings = append(errStrings, err.Error())
	}

	//next the intraday table
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, IntradayTableName) {
		errStrings = append(errStrings, err.Error())
	}

	//last is the quotes tables
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, StockQuotesTableName) {
		errStrings = append(errStrings, err.Error())
	}
	for _, err := range sd.deleteStockSymbolFromTable(ctx, symbol, OptionQuotesTableName) {
		errStrings = append(errStrings, err.Error())
	}

	if len(errStrings) != 0 {
		return fmt.Errorf("StockDelete: had erros while trying to delete %v : %v", symbol, strings.Join(errStrings, ";"))
	}
	return nil
}

func (sd *SibylDatabase) deleteStockSymbolFromTable(ctx context.Context, symbol core.StockSymbolType, tableName string) []error {
	errors := make([]error, 0)
	//we first check if it's there and if so then we do a delete
	hasItStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE `symbol` LIKE '%v';", SibylDatabaseName, tableName, symbol)

	rows, err := sd.DBConn.QueryContext(ctx, hasItStr)
	if rows.Next() {

		deleteStr := fmt.Sprintf("DELETE FROM `%v`.`%v` WHERE `symbol` LIKE '%v';", SibylDatabaseName, tableName, symbol)
		_, err = sd.DBConn.ExecContext(ctx, deleteStr)
		if err != nil {
			errors = append(errors, err)
		}
	}
	rows.Close()

	return errors
}
func (sd *SibylDatabase) StockEnableDownloading(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "downloadStatus", string(core.ActivityEnabled))
}

func (sd *SibylDatabase) StockDisableDownloading(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "downloadStatus", string(core.ActivityDisabled))
}

func (sd *SibylDatabase) StockEnableHistory(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "historyStatus", string(core.ActivityEnabled))
}

func (sd *SibylDatabase) StockDisableHistory(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "historyStatus", string(core.ActivityDisabled))
}

func (sd *SibylDatabase) StockDisableStableQuotes(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "stableQuotesStatus", string(core.ActivityDisabled))
}

func (sd *SibylDatabase) StockEnableStableQuotes(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "stableQuotesStatus", string(core.ActivityEnabled))
}

func (sd *SibylDatabase) StockDisableQuotes(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "quotesStatus", string(core.ActivityDisabled))
}

func (sd *SibylDatabase) StockEnableQuotes(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "quotesStatus", string(core.ActivityEnabled))
}

func (sd *SibylDatabase) StockEnableIntradayHistory(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "intradayHistoryStatus", string(core.ActivityEnabled))
}

func (sd *SibylDatabase) StockDisableIntradayHistory(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "intradayHistoryStatus", string(core.ActivityDisabled))
}

func (sd *SibylDatabase) StockRevalidate(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "downloadStatus", string(core.ValidationPending))
}

func (sd *SibylDatabase) StockValidate(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "validationStatus", string(core.ValidationValid))
}
func (sd *SibylDatabase) StockInvalidate(ctx context.Context, symbol core.StockSymbolType) error {
	return sd.stockOneElementChange(ctx, symbol, "validationStatus", string(core.ValidationInvalid))
}

func (sd *SibylDatabase) GetAgent(ctx context.Context) (core.SibylAgent, error) {
	creds, err := sd.GetCreds(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAgent: had an error while getting creds from database: %v", err)
	}

	switch creds.AgentSelection() {
	case core.AgentSelectionNone:
		return nil, fmt.Errorf("GetAgent: no agent is assinged")
	case core.AgentSelectionAlly:
		return ally.NewAllyAgent(
			creds.ConsumerKey(),
			creds.ConsumerSecret(),
			creds.Token(),
			creds.TokenSecret(),
		), nil
	case core.AgentSelectionTDAmeritrade:
		return nil, fmt.Errorf("GetAgent: TD Ameritrade isn't implemented yet")
	}
	return nil, fmt.Errorf("GetAgent: creds retrevied from database was corrupt")
}

func (sd *SibylDatabase) LastHistoryDate(ctx context.Context, symbol core.StockSymbolType) (core.DateType, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE `symbol` = '%v' ORDER by `timestamp` desc LIMIT 1;", SibylDatabaseName, HistoryTableName, symbol)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return core.NewDateTypeFromUnix(0), fmt.Errorf("LastHistoryDate: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if _, history, err := scanners.ScanSibylHistoryRecordRow(rows); err != nil {
			return core.NewDateTypeFromUnix(0), fmt.Errorf("LastHistoryDate: failed with error: %v", err)
		} else {
			return history.Timestamp, nil
		}
	}
	return core.NewDateTypeFromUnix(0), nil
}

func (sd *SibylDatabase) LastIntradayHistoryDate(ctx context.Context, symbol core.StockSymbolType) (core.TimestampType, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE `symbol` = '%v' ORDER by `timestamp` desc LIMIT 1;", SibylDatabaseName, IntradayTableName, symbol)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return core.NewTimestampTypeFromUnix(0), fmt.Errorf("LastIntradayHistoryDate: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if _, history, err := scanners.ScanSibylIntradayRecordRow(rows); err != nil {
			return core.NewTimestampTypeFromUnix(0), fmt.Errorf("LastIntradayHistoryDate: failed with error: %v", err)
		} else {
			return history.Timestamp, nil
		}
	}
	return core.NewTimestampTypeFromUnix(0), nil

}
func (sd *SibylDatabase) SetOptionsForStock(ctx context.Context, symbol core.StockSymbolType, optionsSymbols []*core.OptionSymbolType) error {
	deleteStr := fmt.Sprintf("DELETE FROM `%v`.`%v` WHERE (`symbol` = '%v');", SibylDatabaseName, OptionsTableName, symbol)

	if _, err := sd.DBConn.ExecContext(ctx, deleteStr); err != nil {
		return fmt.Errorf("SetOptionsForStock: failed to during deleting old data: %v", err)
	}

	records := make([]DatabaseStringer, len(optionsSymbols))
	for i, q := range optionsSymbols {
		tmp := core.SibylOptionRecord(*q)
		records[i] = &tmp
	}

	err := sd.loadRecords(ctx,
		records,
		SibylDatabaseName,
		OptionsTableName,
		[]string{},
		[]string{"expiration", "optionType", "strikePrice", "symbol"},
		NoAction,
	)

	if err != nil {
		return fmt.Errorf("SetOptionRecordsForStock: had the following error while loading options: %v", err)
	}

	return nil
}
func (sd *SibylDatabase) GetOptionsFor(ctx context.Context, symbols map[core.StockSymbolType]bool) ([]*core.OptionSymbolType, error) {
	if len(symbols) == 0 {
		return []*core.OptionSymbolType{}, nil
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE ", SibylDatabaseName, OptionsTableName))
	index := 0
	for symbol := range symbols {
		if index == 0 {
			queryBuilder.WriteString(fmt.Sprintf("`symbol` = '%v'", string(symbol)))
		} else {
			queryBuilder.WriteString(fmt.Sprintf("or `symbol` = '%v'", string(symbol)))
		}
		index++
	}
	queryBuilder.WriteString(";")

	rows, err := sd.DBConn.QueryContext(ctx, queryBuilder.String())
	if err != nil {
		return []*core.OptionSymbolType{}, fmt.Errorf("GetOptionsFor: failed to during getting option data: %v", err)
	}
	defer rows.Close()

	options := make([]*core.OptionSymbolType, 0)
	errString := make([]string, 0)

	for rows.Next() {
		if option, err := scanners.ScanSibylOptionRecordRow(rows); err != nil {
			errString = append(errString, err.Error())
		} else {
			tmp := core.OptionSymbolType(*option)
			options = append(options, &tmp)
		}
	}
	if len(errString) > 0 {
		return options, fmt.Errorf("GetOptionsFor: had some errors while scanning values form database: %v", strings.Join(errString, ","))
	}
	return options, nil
}
func (sd *SibylDatabase) Close() {
	sd.DBConn.Close()
}
func (sd *SibylDatabase) StockDisableAll(ctx context.Context) error {
	updateStr := fmt.Sprintf(
		"UPDATE `%v`.`%v` SET "+
			"`downloadStatus` = 'disabled', "+
			"`historyStatus` = 'disabled', "+
			"`intradayHistoryStatus` = 'disabled', "+
			"`quotesStatus` = 'disabled', "+
			"`stableQuotesStatus` = 'disabled' "+
			"WHERE (`id` > 0);", SibylDatabaseName, StocksTableName)

	if _, err := sd.DBConn.ExecContext(ctx, updateStr); err != nil {
		return fmt.Errorf("StockEnableAll: failed to while enabling all stocks error: %v", err)
	}
	return nil
}
func (sd *SibylDatabase) StockEnableAll(ctx context.Context) error {
	updateStr := fmt.Sprintf(
		"UPDATE `%v`.`%v` SET "+
			"`downloadStatus` = 'enabled', "+
			"`historyStatus` = 'enabled', "+
			"`intradayHistoryStatus` = 'enabled', "+
			"`quotesStatus` = 'enabled', "+
			"`stableQuotesStatus` = 'enabled' "+
			"WHERE (`id` > 0 );", SibylDatabaseName, StocksTableName)

	if _, err := sd.DBConn.ExecContext(ctx, updateStr); err != nil {
		return fmt.Errorf("StockEnableAll: failed to while enabling all stocks error: %v", err)
	}
	return nil
}
func (sd *SibylDatabase) StockSetExchangeInfoAndName(killCtx context.Context, symbol core.StockSymbolType, hasOptions bool, exchange string, exchangeDescription string, name string) error {
	if hasOptions {
		if err := sd.stockOneElementChange(killCtx, symbol, "hasOptions", "yes"); err != nil {
			return err
		}
	} else {
		if err := sd.stockOneElementChange(killCtx, symbol, "hasOptions", "no"); err != nil {
			return err
		}
	}
	if err := sd.stockOneElementChange(killCtx, symbol, "exchange", exchange); err != nil {
		return err
	}
	if err := sd.stockOneElementChange(killCtx, symbol, "exchangeDescription", exchangeDescription); err != nil {
		return err
	}

	escapeSingleQuote := strings.NewReplacer("'", "''")

	if err := sd.stockOneElementChange(killCtx, symbol, "name", escapeSingleQuote.Replace(name)); err != nil {
		return err
	}
	return nil
}
func (sd *SibylDatabase) GetHistory(ctx context.Context, stockSymbol core.StockSymbolType, startTimestamp core.DateType, endTimestamp core.DateType) ([]*core.SibylHistoryRecord, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE `symbol` = '%v' AND `timestamp` >= %v AND `timestamp` < %v", SibylDatabaseName, HistoryTableName, stockSymbol, startTimestamp.Unix(), endTimestamp.Unix())

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("GetHistory: had error running query: %v", err)
	}
	defer rows.Close()

	toReturn := make([]*core.SibylHistoryRecord, 0)
	for rows.Next() {
		if _, record, err := scanners.ScanSibylHistoryRecordRow(rows); err != nil {
			return nil, fmt.Errorf("GetHistory: had error reading results: %v", err)
		} else {
			toReturn = append(toReturn, record)
		}
	}
	return toReturn, nil
}

func (sd *SibylDatabase) GetIntraday(ctx context.Context, symbol core.StockSymbolType, startTimestamp core.TimestampType, endTimestamp core.TimestampType) ([]*core.SibylIntradayRecord, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` WHERE `symbol` = '%v' AND `timestamp` >= %v AND `timestamp` < %v", SibylDatabaseName, IntradayTableName, symbol, startTimestamp.Unix(), endTimestamp.Unix())

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return nil, fmt.Errorf("GetIntraday: had error running query: %v", err)
	}
	defer rows.Close()

	toReturn := make([]*core.SibylIntradayRecord, 0)
	for rows.Next() {
		if _, record, err := scanners.ScanSibylIntradayRecordRow(rows); err != nil {
			return nil, fmt.Errorf("GetIntraday: had error reading results: %v", err)
		} else {
			toReturn = append(toReturn, record)
		}
	}
	return toReturn, nil
}

func (sd *SibylDatabase) LastHistoryRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, HistoryTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastHistoryRecordID: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylHistoryRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastHistoryRecordID: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}

func (sd *SibylDatabase) LastIntradayRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, IntradayTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastIntradayRecord: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylIntradayRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastIntradayRecord: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}

func (sd *SibylDatabase) LastStableStockQuoteRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, StableStockQuotesTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastStableStockQuoteRecordID: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylStableStockQuoteRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastStableStockQuoteRecordID: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}

func (sd *SibylDatabase) LastStableOptionQuoteRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, StableOptionQuotesTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastStableOptionQuoteRecordID: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylStableOptionQuoteRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastStableOptionQuoteRecordID: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}

func (sd *SibylDatabase) LastStockQuoteRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, StockQuotesTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastStockQuoteRecordID: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylStockQuoteRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastStockQuoteRecordID: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}

func (sd *SibylDatabase) LastOptionQuoteRecordID(ctx context.Context) (string, error) {
	queryStr := fmt.Sprintf("SELECT * FROM `%v`.`%v` ORDER by `id` desc LIMIT 1;", SibylDatabaseName, OptionQuotesTableName)

	rows, err := sd.DBConn.QueryContext(ctx, queryStr)
	if err != nil {
		return "", fmt.Errorf("LastOptionQuoteRecordID: failed to execute query %v, had error: %v", queryStr, err)
	}
	defer rows.Close()

	if rows.Next() {
		if id, _, err := scanners.ScanSibylOptionQuoteRecordRow(rows); err != nil {
			return "", fmt.Errorf("LastOptionQuoteRecordID: failed with error: %v", err)
		} else {
			return id, nil
		}
	}
	return "", nil
}
