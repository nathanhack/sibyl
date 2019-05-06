package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

const (
	IntradayStateActive IntradayStateType = 1
	IntradayStateDaily  IntradayStateType = 0 //Default
)

type IntradayStateType int

func (ihst *IntradayStateType) IsActive() bool {
	return IntradayStateType(*ihst) == IntradayStateActive
}

func (ihst *IntradayStateType) IsDaily() bool {
	return IntradayStateType(*ihst) == IntradayStateDaily
}

func (ihst IntradayStateType) String() string {
	if int(ihst) == 0 {
		return "daily"
	}
	tmp := int(ihst)

	if tmp&int(IntradayStateActive) > 0 {
		return "active"
	}

	return "daily"
}

func (ihst *IntradayStateType) Scan(value interface{}) error {
	if value == nil {
		*ihst = 0
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*ihst = IntradayStateType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows IntradayStateType", value, value)
		}
		*ihst = IntradayStateType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStateType", value, value)
		}
		*ihst = IntradayStateType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStateType", value, value)
		}
		*ihst = IntradayStateType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to IntradayStateType", value, value)
		}
		*ihst = IntradayStateType(i)
		return nil
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan IntradayStateType: %v", value)
}

func (ihst IntradayStateType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(ihst), nil
}
