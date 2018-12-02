package core

import (
	"database/sql/driver"
	"fmt"
)

type ValidationStatusType string

const (
	ValidationPending ValidationStatusType = "pending"
	ValidationValid   ValidationStatusType = "valid"
	ValidationInvalid ValidationStatusType = "invalid"
)

func (vst *ValidationStatusType) Scan(value interface{}) error {
	if value == nil {
		*vst = ValidationPending
		return nil
	}

	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			switch ValidationStatusType(v) {
			case ValidationValid:
				*vst = ValidationValid
			case ValidationInvalid:
				*vst = ValidationInvalid
			case ValidationPending:
				*vst = ValidationPending
			default:
				*vst = ValidationPending
			}
			return nil
		}
		if v, ok := sv.([]byte); ok {
			switch ValidationStatusType(string(v)) {
			case ValidationValid:
				*vst = ValidationValid
			case ValidationInvalid:
				*vst = ValidationInvalid
			case ValidationPending:
				*vst = ValidationPending
			default:
				*vst = ValidationPending
			}
			return nil
		}
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan AgentSelectionType with value of :%v", value)
}

func (vst ValidationStatusType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return string(vst), nil
}
