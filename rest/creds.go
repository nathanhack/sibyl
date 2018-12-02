package rest

type Creds struct {
	AgentSelection         string     `json:"Agent"`
	ConsumerKey            string     `json:"ConsumerKey"`
	ConsumerSecret         string     `json:"ConsumerSecret"`
	Token                  string     `json:"Token"`
	TokenSecret            string     `json:"TokenSecret"`
	UrlRedirect            string     `json:"UrlRedirect"`
	AccessToken            string     `json:"AccessToken"`
	RefreshToken           string     `json:"RefreshToken"`
	ExpireTimestamp        int64      `json:"ExpireTimestamp"`
	RefreshExpireTimestamp int64      `json:"RefreshExpireTimestamp"`
	ErrorState             ErrorState `json:"ErrorState"`
}
