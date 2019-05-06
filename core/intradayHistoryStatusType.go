package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

const (
	IntradaystatusDisabled IntradayStatusType = 0 //default value
	IntradayStatus1Min     IntradayStatusType = 1 //bit flag
	IntradayStatus5Min     IntradayStatusType = 2 //bit flag
	IntradayStatusTicks    IntradayStatusType = 4 //bit flag
)

//IntradayStatusType is a bitflag value
type IntradayStatusType int

func (ist *IntradayStatusType) HasTicks() bool {
	return int(*ist)&int(IntradayStatusTicks) > 0
}

func (ist *IntradayStatusType) Has1Min() bool {
	return int(*ist)&int(IntradayStatus1Min) > 0
}

func (ist *IntradayStatusType) Has5Min() bool {
	return int(*ist)&int(IntradayStatus5Min) > 0
}

func (ist IntradayStatusType) String() string {
	if int(ist) == 0 {
		return "disabled"
	}
	tmp := int(ist)
	toReturn := make([]string, 0)

	if tmp&int(IntradayStatusTicks) > 0 {
		toReturn = append(toReturn, "ticks")
	}
	if tmp&int(IntradayStatus1Min) > 0 {
		toReturn = append(toReturn, "1min")
	}
	if tmp&int(IntradayStatus5Min) > 0 {
		toReturn = append(toReturn, "5min")
	}

	return strings.Join(toReturn, ",")
}

func (ist *IntradayStatusType) Scan(value interface{}) error {
	if value == nil {
		*ist = 0
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*ist = IntradayStatusType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows IntradayStatusType", value, value)
		}
		*ist = IntradayStatusType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStatusType", value, value)
		}
		*ist = IntradayStatusType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStatusType", value, value)
		}
		*ist = IntradayStatusType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStatusType", value, value)
		}
		*ist = IntradayStatusType(i)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("failed to scan IntradayStatusType: %v", value)
}

func (ist IntradayStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(ist), nil
}
