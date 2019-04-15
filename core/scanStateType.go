package core

import (
	"database/sql/driver"
	"fmt"
)

type ScanStateType string

const (
	ScanUnknown  = ScanStateType("unknown")
	ScanScanning = ScanStateType("scanning")
	ScanScanned  = ScanStateType("scanned")
)

func (vst *ScanStateType) Scan(value interface{}) error {
	if value == nil {
		*vst = ScanUnknown
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err == nil {
		if v, ok := sv.(string); ok {
			switch ScanStateType(v) {
			case ScanScanned:
				*vst = ScanScanned
			case ScanScanning:
				*vst = ScanScanning
			default:
				*vst = ScanUnknown
			}
			return nil
		}

		if v, ok := sv.([]byte); ok {
			switch ScanStateType(string(v)) {
			case ScanScanned:
				*vst = ScanScanned
			case ScanScanning:
				*vst = ScanScanning
			default:
				*vst = ScanUnknown
			}
			return nil
		}
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan ActivityStatusType: %v", value)
}

func (vst ScanStateType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return string(vst), nil
}
