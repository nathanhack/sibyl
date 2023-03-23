package config

type ConfigFileType struct {
	Agents struct {
		Alpaca struct {
			Live struct {
				ApiKeyID  string `json:"apiKeyid"`
				SecretKey string `json:"secretkey"`
				Url       string `json:"url"`
			} `json:"live"`
			Paper struct {
				ApiKeyID  string `json:"apikeyid"`
				SecretKey string `json:"secretkey"`
				Url       string `json:"url"`
			} `json:"paper"`
			Endpoint string `json:"endpoint`
			Plan     string `json:"plan"`
		} `json:"alpaca"`
		PolygonIO struct {
			ApiKey string `json:"apikey"`
			Plan   string `json:"plan"`
		} `json:"polygonio"`
	} `json:"agents"`
	Database struct {
		Dialect string `json:"dialect"`
		DSN     string `json:"dsn"`
	} `json:"database"`
	Address string `json:"address"`
	Logging struct {
		Directory string `json:"directory"`
		Level     string `json:"level"`
	} `json:"logging"`
}

var State ConfigFileType
