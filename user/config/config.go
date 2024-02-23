package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	Port          int           `mapstructure:"port"`
}
