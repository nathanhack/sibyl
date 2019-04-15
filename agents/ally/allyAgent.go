package ally

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/nathanhack/sibyl/agents/internal/flowLimiter"
	"github.com/nathanhack/sibyl/agents/internal/ratelimiter"
	"net/http"
)

//we make some global rateLimiters to be used by all Ally agents
// see https://www.ally.com/api/invest/documentation/rate-limiting/ for more details

// From the website these are the actual rate limits
// 40 requests per min for Trade submission calls
// 60 requests per min for market/quotes calls
// 180 requests per min for all other calls
var rateLimitTrades = ratelimiter.New((40.0) * (1.0 / 60.0))
var rateLimitMarket = ratelimiter.New((60.0) * (1.0 / 60.0))
var rateLimitOthers = ratelimiter.New((180.0) * (1.0 / 60.0))

//Ally only allows 5 concurrent requests
var quoteLimits = flowLimiter.New(5)

// for everything that isn't quotes or stable quotes will additionally be rate limited
var rateLimitMarketLowPriority = ratelimiter.New(1.0)

type AllyAgent struct {
	consumerKey                string
	consumerSecret             string
	oAuthToken                 string
	oAuthTokenSecret           string
	httpClient                 *http.Client
	rateLimitTrades            *ratelimiter.RateLimiter
	rateLimitMarketCalls       *ratelimiter.RateLimiter
	rateLimitOtherCalls        *ratelimiter.RateLimiter
	rateLimitMarketLowPriority *ratelimiter.RateLimiter
	concurrentLimit            *flowLimiter.FlowLimiter
}

func NewAllyAgent(
	consumerKey string,
	consumerSecret string,
	oAuthToken string,
	oAuthTokenSecret string,
) *AllyAgent {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(oAuthToken, oAuthTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &AllyAgent{
		consumerKey:                consumerKey,
		consumerSecret:             consumerSecret,
		oAuthToken:                 oAuthToken,
		oAuthTokenSecret:           oAuthTokenSecret,
		httpClient:                 httpClient,
		rateLimitTrades:            rateLimitTrades,
		rateLimitMarketCalls:       rateLimitMarket,
		rateLimitOtherCalls:        rateLimitOthers,
		rateLimitMarketLowPriority: rateLimitMarketLowPriority,
		concurrentLimit:            quoteLimits,
	}
}

func (ag *AllyAgent) ConsumerKey() string {
	return ag.consumerKey
}

func (ag *AllyAgent) ConsumerSecret() string {
	return ag.consumerSecret
}

func (ag *AllyAgent) OAuthToken() string {
	return ag.oAuthToken
}

func (ag *AllyAgent) OAuthTokenSecret() string {
	return ag.oAuthTokenSecret
}

func (ac *AllyAgent) String() string {
	return fmt.Sprintf("{ConsumerKey: %v, ConsumerSecret: %v, OAuthToken: %v,OAuthTokenSecret: %v}", ac.consumerKey, ac.consumerSecret, ac.oAuthToken, ac.oAuthTokenSecret)
}
