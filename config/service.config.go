package config

type Server struct {
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	SentryUrl string `mapstructure:"sentry_url"`
}

type GRPC struct {
	OAuthPort     string `mapstructure:"oauth_port"`
	WorkspacePort string `mapstructure:"workspace_port"`
}

type JWT struct {
	Issuer            string `mapstructure:"issuer"`
	WSAccessTokenKey  string `mapstructure:"ws_access_token_key"`
	WSRefreshTokenKey string `mapstructure:"ws_refresh_token_key"`
}
