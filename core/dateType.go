package core

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type DateType struct {
	int64
}

func (tt DateType) String() string {
	return fmt.Sprintf("%v", tt.int64)
}

func NewDateType(year, month, day int) DateType {
	return DateType{
		time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix(),
	}
}

func NewDateTypeFromTime(toConvert time.Time) DateType {
	year, month, day := toConvert.Date()
	return NewDateType(year, int(month), day)

}

func NewDateTypeFromUnix(timestamp int64) DateType {
	return NewDateTypeFromTime(time.Unix(timestamp, 0).Local())
}

func (tt DateType) Time() time.Time {
	return time.Unix(tt.int64, 0).Local()
}

func (tt DateType) Unix() int64 {
	return tt.int64
}

func (tt *DateType) Scan(value interface{}) error {
	if value == nil {
		*tt = DateType{0}
		return nil
	}
	switch value.(type) {
	case int, int8, int16, int32, int64:
		*tt = DateType{value.(int64)}
		return nil
	case uint, uint8, uint16, uint32, uint64:
		u64 := value.(uint64)
		if u64 > (1<<63)-1 {
			return fmt.Errorf("sql/driver: value %v (type %T) overflows DateType", value, value)
		}
		*tt = DateType{int64(u64)}
		return nil
	case string:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to DateType", value, value)
		}
		*tt = DateType{i}
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to DateType", value, value)
		}
		*tt = DateType{i}
		return nil

	case []int8:
		tmp := []byte{}
		for _, i := range value.([]int8) {
			tmp = append(tmp, byte(i))
		}
		i, err := strconv.ParseInt(string(tmp), 10, 64)
		if err != nil {
			return fmt.Errorf("sql/driver: value %v (type %T) can't be converted to DateType", value, value)
		}
		*tt = DateType{i}
		return nil
	}

	return fmt.Errorf("sql/driver: unsupported value %v (type %T) converting to DateType", value, value)
}

func (tt DateType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is int64
	return int64(tt.int64), nil
}
func (tt DateType) AddDate(years int, months int, days int) DateType {
	return NewDateTypeFromTime(tt.Time().AddDate(years, months, days))
}
func (tt DateType) Before(date DateType) bool {
	return tt.int64 < date.int64
}
func (tt DateType) IsWeekDay() bool {
	t := time.Unix(tt.int64, 0).Weekday()
	return time.Sunday < t && t < time.Saturday
}

func (tt DateType) Day() time.Weekday {
	return time.Unix(tt.int64, 0).Weekday()
}
