package config

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type GRPC struct {
	Port int `mapstructure:"port"`
}

type JWT struct {
	Issuer            string `mapstructure:"issuer"`
	WSAccessTokenKey  string `mapstructure:"ws_access_token_key"`
	WSRefreshTokenKey string `mapstructure:"ws_refresh_token_key"`
}
