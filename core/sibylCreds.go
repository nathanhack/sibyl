package core

import "fmt"

type SibylCreds struct {
	agentSelection         AgentSelectionType
	consumerKey            string
	consumerSecret         string
	token                  string
	tokenSecret            string
	urlRedirect            string
	accessToken            string
	refreshToken           string
	expireTimestamp        int64
	refreshExpireTimestamp int64
}

func DefaultSibylCreds() *SibylCreds {
	return &SibylCreds{
		AgentSelectionNone,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		0,
		0,
	}
}

func NewSibylCreds(
	agentSelection AgentSelectionType,
	consumerKey string,
	consumerSecret string,
	token string,
	tokenSecret string,
	urlRedirect string,
	accessToken string,
	refreshToken string,
	expireTimestamp int64,
	refreshExpireTimestamp int64,
) *SibylCreds {
	return &SibylCreds{
		agentSelection,
		consumerKey,
		consumerSecret,
		token,
		tokenSecret,
		urlRedirect,
		accessToken,
		refreshToken,
		expireTimestamp,
		refreshExpireTimestamp,
	}
}

func (sc *SibylCreds) AgentSelection() AgentSelectionType {
	return sc.agentSelection
}
func (sc *SibylCreds) ConsumerKey() string {
	return sc.consumerKey
}
func (sc *SibylCreds) ConsumerSecret() string {
	return sc.consumerSecret
}
func (sc *SibylCreds) Token() string {
	return sc.token
}
func (sc *SibylCreds) TokenSecret() string {
	return sc.tokenSecret
}
func (sc *SibylCreds) UrlRedirect() string {
	return sc.urlRedirect
}
func (sc *SibylCreds) AccessToken() string {
	return sc.accessToken
}
func (sc *SibylCreds) RefreshToken() string {
	return sc.refreshToken
}
func (sc *SibylCreds) ExpireTimestamp() int64 {
	return sc.expireTimestamp
}
func (sc *SibylCreds) RefreshExpireTimestamp() int64 {
	return sc.refreshExpireTimestamp
}

func (sc *SibylCreds) String() string {
	return sc.StringBlindWithDelimiter(",", "", true)
}

func (sc *SibylCreds) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
		esc, sc.agentSelection, esc, delimiter,
		esc, sc.consumerKey, esc, delimiter,
		esc, sc.consumerSecret, esc, delimiter,
		esc, sc.token, esc, delimiter,
		esc, sc.tokenSecret, esc, delimiter,
		esc, sc.urlRedirect, esc, delimiter,
		esc, sc.accessToken, esc, delimiter,
		esc, sc.refreshToken, esc, delimiter,
		sc.expireTimestamp, delimiter,
		sc.refreshExpireTimestamp,
	)
}
