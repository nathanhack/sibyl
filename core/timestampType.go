package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type TimestampType struct {
	int64
}

func (tt TimestampType) String() string {
	return fmt.Sprintf("%v", tt.int64)
}

func NewTimestampTypeFromTime(time time.Time) TimestampType {
	return TimestampType{time.Unix()}
}

func NewTimestampTypeFromUnix(timestamp int64) TimestampType {
	return TimestampType{timestamp}
}

func (tt *TimestampType) Scan(value interface{}) error {
	if value == nil {
		*tt = TimestampType{0}
		return nil
	}
	switch value.(type) {
	case int, int8, int16, int32, int64:
		*tt = TimestampType{value.(int64)}
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows TimestampType", value, value)
		}
		*tt = TimestampType{int64(u64)}
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to TimestampType", value, value)
		}
		*tt = TimestampType{i}
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to TimestampType", value, value)
		}
		*tt = TimestampType{i}
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to TimestampType", value, value)
		}
		*tt = TimestampType{i}
		return nil
	}

	return fmt.Errorf("sql/driver: unsupported value %v (type %T) converting to TimestampType", value, value)

}

func (tt TimestampType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is int64
	return int64(tt.int64), nil
}

func (tt TimestampType) AddDate(years int, months int, days int) TimestampType {
	return NewTimestampTypeFromTime(tt.Time().AddDate(years, months, days))
}

func (tt TimestampType) Time() time.Time {
	return time.Unix(tt.int64, 0)
}

func (tt TimestampType) Unix() int64 {
	return tt.int64
}

func (tt TimestampType) Before(date TimestampType) bool {
	return tt.int64 < date.int64
}
func (tt TimestampType) IsWeekDay() bool {
	t := time.Unix(tt.int64, 0).Weekday()
	return time.Sunday < t && t < time.Saturday
}
