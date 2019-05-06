package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type ValidationStatusType int

const (
	ValidationPending ValidationStatusType = 0 //default
	ValidationValid   ValidationStatusType = 1
	ValidationInvalid ValidationStatusType = 2
)

func (vst ValidationStatusType) String() string {
	switch vst {
	case ValidationValid:
		return "valid"
	case ValidationInvalid:
		return "invalid"
	default:
		return "pending"
	}
}

func (vst *ValidationStatusType) Scan(value interface{}) error {
	if value == nil {
		*vst = ValidationPending
		return nil
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		*vst = ValidationStatusType(value.(int64))
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows ValidationStatusType", value, value)
		}
		*vst = ValidationStatusType(int64(u64))
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ValidationStatusType", value, value)
		}
		*vst = ValidationStatusType(i)
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ValidationStatusType", value, value)
		}
		*vst = ValidationStatusType(i)
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to ValidationStatusType", value, value)
		}
		*vst = ValidationStatusType(i)
		return nil
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan ValidationStatusType with value of: %v", value)
}

func (vst ValidationStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return int(vst), nil
}
