package ally

import (
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func toAllySymbol(symbol core.OptionSymbolType) string {
	// for ally the symbol has the following format
	//(1) - Underlying symbol
	//(2) – 2-digit expiration year
	//(3) – 2-digit expiration month
	//(4) – 2-digit expiration day
	//(5) – "C" for Call or "P" for Put
	//(6) – 8-digit strike price
	// ex.  CAT181024C00150000

	return fmt.Sprintf("%v%v%v%08v",
		symbol.Symbol,
		symbol.Expiration.Time().Local().Format("060102"),
		symbol.OptionType[:1],
		int(symbol.StrikePrice*1000),
	)
}

func toOptionSymbol(allySymbol string) (*core.OptionSymbolType, error) {
	// for ally the symbol has the following format
	//(1) - Underlying symbol
	//(2) – 2-digit expiration year
	//(3) – 2-digit expiration month
	//(4) – 2-digit expiration day
	//(5) – "C" for Call or "P" for Put
	//(6) – 8-digit strike price
	// ex.  CAT181024C00150000

	r := regexp.MustCompile(`(?P<stock>.*?)(?P<date>\d{6})(?P<type>.)(?P<strike>\d{8})`)
	match := r.FindStringSubmatch(allySymbol)
	if match == nil {
		return nil, fmt.Errorf("toOptionSymbol: was unable to find symbol in: %v", allySymbol)
	}
	stock := match[1]
	if len(stock) == 0 {
		return nil, fmt.Errorf("toOptionSymbol: was unable to extract symbol from: %v", allySymbol)
	}

	var err error
	//IMPORTANT: parse this timestamp with local time otherwise the time stamp is incorrect
	timestamp, err := time.ParseInLocation("060102", match[2], time.Local)
	if err != nil {
		return nil, fmt.Errorf("toOptionSymbol: was unable to extract symbol from: %v", allySymbol)
	}

	var equityType core.EquityType
	switch strings.ToLower(match[3]) {
	case "c":
		equityType = core.CallEquity
	case "p":
		equityType = core.PutEquity
	default:
		return nil, fmt.Errorf("toOptionSymbol: there was a problem extracting equity information from: %v", allySymbol)
	}

	strikeInt, err := strconv.Atoi(match[4])
	if err != nil {
		return nil, fmt.Errorf("toOptionSymbol: was unable to extract symbol from: %v", allySymbol)
	}

	return &core.OptionSymbolType{
		Expiration:  core.NewDateTypeFromTime(timestamp.Local()), //number of seconds since epoch
		OptionType:  equityType,
		Symbol:      core.StockSymbolType(stock),
		StrikePrice: float64(strikeInt) / 1000,
	}, nil
}
