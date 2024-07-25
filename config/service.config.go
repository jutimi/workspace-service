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
