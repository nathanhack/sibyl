package core

import (
	"database/sql/driver"
	"fmt"
)

type ActivityStatusType string

const (
	ActivityEnabled  ActivityStatusType = "enabled"
	ActivityDisabled ActivityStatusType = "disabled"
)

func (vst *ActivityStatusType) Scan(value interface{}) error {
	if value == nil {
		*vst = ActivityDisabled
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err == nil {
		if v, ok := sv.(string); ok {
			switch ActivityStatusType(v) {
			case ActivityEnabled:
				*vst = ActivityEnabled
			default:
				*vst = ActivityDisabled
			}
			return nil
		}

		if v, ok := sv.([]byte); ok {
			switch ActivityStatusType(string(v)) {
			case ActivityEnabled:
				*vst = ActivityEnabled
			default:
				*vst = ActivityDisabled
			}
			return nil
		}
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan ActivityStatusType: %v", value)
}

func (vst ActivityStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return string(vst), nil
}
