package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type ActivityStatusType int

const (
	ActivityEnabled  ActivityStatusType = 1
	ActivityDisabled ActivityStatusType = 0 //default
)

func (ast *ActivityStatusType) String() string {
	if *ast == ActivityEnabled {
		return "enabled"
	}
	return "disabled"
}

func (ast *ActivityStatusType) Scan(value interface{}) error {
	if value == nil {
		*ast = ActivityDisabled
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*ast = ActivityStatusType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows ActivityStatusType", value, value)
		}
		*ast = ActivityStatusType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ActivityStatusType", value, value)
		}
		*ast = ActivityStatusType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ActivityStatusType", value, value)
		}
		*ast = ActivityStatusType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ActivityStatusType", value, value)
		}
		*ast = ActivityStatusType(i)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("failed to scan ActivityStatusType: %v", value)
}

func (ast ActivityStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(ast), nil
}
