package core

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
)

type EquityType string

const (
	UnknownEquity EquityType = "UNKNOWN"
	CallEquity    EquityType = "CALL"
	PutEquity     EquityType = "PUT"
)

func (vst *EquityType) Scan(value interface{}) error {
	if value == nil {
		*vst = UnknownEquity
		return nil
	}

	if sv, err := driver.String.ConvertValue(value); err == nil {
		if bs, ok := sv.([]byte); ok {
			switch EquityType(string(bs)) {
			case CallEquity:
				*vst = CallEquity
			case PutEquity:
				*vst = PutEquity
			case UnknownEquity:
				*vst = UnknownEquity
			default:
				*vst = UnknownEquity
			}
			return nil
		}
		if v, ok := sv.(string); ok {
			switch EquityType(v) {
			case CallEquity:
				*vst = CallEquity
			case PutEquity:
				*vst = PutEquity
			case UnknownEquity:
				*vst = UnknownEquity
			default:
				*vst = UnknownEquity
			}
			return nil
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan EquityType")
}

func (vst EquityType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return string(vst), nil
}

type NullEquityType struct {
	Type  EquityType
	Valid bool
}

func (net *NullEquityType) Scan(value interface{}) error {
	if value == nil {
		net.Type, net.Valid = UnknownEquity, false
		return nil
	}
	tmpT := sql.NullString{}
	err := tmpT.Scan(value)

	net.Type, net.Valid = EquityType(tmpT.String), tmpT.Valid
	if net.Valid {
		switch net.Type {
		case CallEquity:
			return nil
		case PutEquity:
			return nil
		case UnknownEquity:
			return nil
		default:
			net.Type = UnknownEquity
			net.Valid = false
		}
	} else {
		net.Type = UnknownEquity
	}
	return err
}

func (nsqt *NullEquityType) Value() (driver.Value, error) {
	if !nsqt.Valid {
		return nil, nil
	}
	return nsqt.Type, nil
}

func nullEquityTypeToString(v NullEquityType, nullString string, stringEscapes bool) string {
	if v.Valid {
		if stringEscapes {
			return fmt.Sprintf("'%v'", v.Type)
		} else {
			return fmt.Sprintf("%v", v.Type)
		}
	}
	return nullString
}
