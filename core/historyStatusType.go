package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

const (
	HistoryStatusDisabled HistoryStatusType = 0 //default value
	HistoryStatusDaily    HistoryStatusType = 1 //bit flag
	HistoryStatusWeekly   HistoryStatusType = 2 //bit flag
	HistoryStatusMonthly  HistoryStatusType = 4 //bit flag
	HistoryStatusYearly   HistoryStatusType = 8 //bit flag
)

//HistorystatusType is a big flag value
type HistoryStatusType int

func (hst *HistoryStatusType) HasDaily() bool {
	return int(*hst)&int(HistoryStatusDaily) > 0
}

func (hst *HistoryStatusType) HasWeekly() bool {
	return int(*hst)&int(HistoryStatusWeekly) > 0
}

func (hst *HistoryStatusType) HasMonthly() bool {
	return int(*hst)&int(HistoryStatusMonthly) > 0
}

func (hst *HistoryStatusType) HasYearly() bool {
	return int(*hst)&int(HistoryStatusYearly) > 0
}

func (hst HistoryStatusType) String() string {
	if int(hst) == 0 {
		return "disabled"
	}
	tmp := int(hst)
	toReturn := make([]string, 0)

	if tmp&int(HistoryStatusDaily) > 0 {
		toReturn = append(toReturn, "daily")
	}
	if tmp&int(HistoryStatusWeekly) > 0 {
		toReturn = append(toReturn, "weekly")
	}
	if tmp&int(HistoryStatusMonthly) > 0 {
		toReturn = append(toReturn, "monthly")
	}
	if tmp&int(HistoryStatusYearly) > 0 {
		toReturn = append(toReturn, "yearly")
	}

	return strings.Join(toReturn, ",")
}

func (hst *HistoryStatusType) Scan(value interface{}) error {
	if value == nil {
		*hst = 0
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*hst = HistoryStatusType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows HistoryStatusType", value, value)
		}
		*hst = HistoryStatusType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to HistoryStatusType", value, value)
		}
		*hst = HistoryStatusType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to HistoryStatusType", value, value)
		}
		*hst = HistoryStatusType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to HistoryStatusType", value, value)
		}
		*hst = HistoryStatusType(i)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("failed to scan HistoryStatusType: %v", value)
}

func (hst HistoryStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(hst), nil
}
