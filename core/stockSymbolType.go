package core

import (
	"database/sql/driver"
	"fmt"
)

type StockSymbolType string

func (vst *StockSymbolType) Scan(value interface{}) error {
	if value == nil {
		*vst = ""
		return nil
	}

	if sv, err := driver.String.ConvertValue(value); err == nil {
		if bs, ok := sv.([]byte); ok {
			*vst = StockSymbolType(string(bs))
			return nil
		}
		if v, ok := sv.(string); ok {
			*vst = StockSymbolType(v)
			return nil
		}
	}

	return fmt.Errorf("sql/driver: unsupported value %v (type %T) converting to StockSymbolType", value, value)

}

func (vst StockSymbolType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is int64
	return string(vst), nil
}
