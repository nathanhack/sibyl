package core

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type DivFrequencyType string

const (
	AnnualDiv     DivFrequencyType = "A"
	SemiAnnualDiv DivFrequencyType = "S"
	QuarterlyDiv  DivFrequencyType = "Q"
	MonthlyDiv    DivFrequencyType = "M"
	NoDiv         DivFrequencyType = "N"
)

type NullDivFrequency struct {
	DivFreq DivFrequencyType
	Valid   bool
}

func (nsqt *NullDivFrequency) Scan(value interface{}) error {
	if value == nil {
		nsqt.DivFreq, nsqt.Valid = "", false
		return nil
	}

	tmpT := sql.NullString{}
	err := tmpT.Scan(value)
	nsqt.DivFreq, nsqt.Valid = DivFrequencyType(tmpT.String), tmpT.Valid
	if nsqt.Valid {
		switch DivFrequencyType(tmpT.String) {
		case AnnualDiv:
			nsqt.DivFreq = AnnualDiv
		case SemiAnnualDiv:
			nsqt.DivFreq = SemiAnnualDiv
		case QuarterlyDiv:
			nsqt.DivFreq = QuarterlyDiv
		case MonthlyDiv:
			nsqt.DivFreq = MonthlyDiv
		case NoDiv:
			nsqt.DivFreq = NoDiv
		default:
			nsqt.DivFreq = NoDiv
		}
	}
	return err
}

func (nsqt *NullDivFrequency) Value() (driver.Value, error) {
	if !nsqt.Valid {
		return nil, nil
	}
	return nsqt.DivFreq, nil
}

func nullDivFrequencyToString(v NullDivFrequency, nullString string, stringEscapes bool) string {
	if v.Valid {
		if stringEscapes {
			return fmt.Sprintf("'%v'", v.DivFreq)
		} else {
			return fmt.Sprintf("%v", v.DivFreq)
		}
	}
	return nullString
}
