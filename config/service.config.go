package config

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type GRPC struct {
	Port int `mapstructure:"port"`
}
