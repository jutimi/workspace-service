package config

type Server struct {
	Port        int    `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	SentryUrl   string `mapstructure:"sentry_url"`
	ServiceName string `mapstructure:"service_name"`
	UptraceDNS  string `mapstructure:"uptrace_dns"`
}

type GRPC struct {
	OAuthUrl     string `mapstructure:"oauth_url"`
	WorkspaceUrl string `mapstructure:"workspace_url"`
}
