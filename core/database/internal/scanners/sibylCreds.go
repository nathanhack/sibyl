package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylCredsRow(rows *sql.Rows) (*core.SibylCreds, error) {
	var id string
	var AgentSelection core.AgentSelectionType
	var ConsumerKey string
	var ConsumerSecret string
	var Token string
	var TokenSecret string
	var UrlRedirect string
	var AccessToken string
	var RefreshToken string
	var ExpireTimestamp int64
	var RefreshExpireTimestamp int64

	err := rows.Scan(&id,
		&AgentSelection,
		&ConsumerKey,
		&ConsumerSecret,
		&Token,
		&TokenSecret,
		&UrlRedirect,
		&AccessToken,
		&RefreshToken,
		&ExpireTimestamp,
		&RefreshExpireTimestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("ScanSibylCredsRow: had error reading results: %v", err)
	}

	return core.NewSibylCreds(
		AgentSelection,
		ConsumerKey,
		ConsumerSecret,
		Token,
		TokenSecret,
		UrlRedirect,
		AccessToken,
		RefreshToken,
		ExpireTimestamp,
		RefreshExpireTimestamp,
	), nil
}
