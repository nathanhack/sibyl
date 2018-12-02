package core

import (
	"database/sql/driver"
	"fmt"
)

type AgentSelectionType string

const (
	AgentSelectionNone         AgentSelectionType = "none"
	AgentSelectionAlly         AgentSelectionType = "ally_invest"
	AgentSelectionTDAmeritrade AgentSelectionType = "td_ameritrade"
)

func (ast *AgentSelectionType) Scan(value interface{}) error {
	if value == nil {
		*ast = AgentSelectionNone
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err == nil {

		if v, ok := sv.(string); ok {
			switch AgentSelectionType(v) {
			case AgentSelectionAlly:
				*ast = AgentSelectionAlly
			case AgentSelectionTDAmeritrade:
				*ast = AgentSelectionTDAmeritrade
			case AgentSelectionNone:
				*ast = AgentSelectionNone
			default:
				*ast = AgentSelectionNone
			}
			return nil
		}
		if v, ok := sv.([]byte); ok {
			switch AgentSelectionType(string(v)) {
			case AgentSelectionAlly:
				*ast = AgentSelectionAlly
			case AgentSelectionTDAmeritrade:
				*ast = AgentSelectionTDAmeritrade
			case AgentSelectionNone:
				*ast = AgentSelectionNone
			default:
				*ast = AgentSelectionNone
			}

			return nil
		}
	}
	// otherwise, return an error
	return fmt.Errorf("failed to scan AgentSelectionType with value:%v", value)
}

func (ast AgentSelectionType) Value() (driver.Value, error) {
	//the driver requires a standard type for
	//which for this type is string
	return string(ast), nil
}
