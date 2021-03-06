package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type OptionStatusType int

const (
	OptionsDisabled OptionStatusType = 0 //default
	OptionsEnabled  OptionStatusType = 1
)

func (ost OptionStatusType) String() string {
	if ost == OptionsEnabled {
		return "enabled"
	}
	return "disabled"
}

func (ost *OptionStatusType) Scan(value interface{}) error {
	if value == nil {
		*ost = OptionsDisabled
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*ost = OptionStatusType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows OptionStatusType", value, value)
		}
		*ost = OptionStatusType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to OptionStatusType", value, value)
		}
		*ost = OptionStatusType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to OptionStatusType", value, value)
		}
		*ost = OptionStatusType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to OptionStatusType", value, value)
		}
		*ost = OptionStatusType(i)
		return nil
	}

	// otherwise, return an error
	return fmt.Errorf("failed to scan OptionStatusType: %v", value)
}

func (ost OptionStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(ost), nil
}
