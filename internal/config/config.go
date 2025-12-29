package config

import (
	"github.com/spf13/viper"
)

var AppVersion = "dev"

type Config struct {
	DBName     string `mapstructure:"DB_NAME"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`
	NewsURL    string `mapstructure:"NEWS_URL"`

	RateLimitPerMinute int `mapstructure:"RATE_LIMIT_PER_MINUTE"`

	CacheTTLSeconds int `mapstructure:"CACHE_TTL_SECONDS"`

	APIKey string `mapstructure:"API_KEY"`

	ProxyMode string `mapstructure:"PROXY_MODE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetDefault("RATE_LIMIT_PER_MINUTE", 60)
	viper.SetDefault("CACHE_TTL_SECONDS", 60)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("NEWS_URL", "https://gita.cherkasyoblenergo.com/obl-main-controller/api/news2?size=18&category=1&page=0")
	viper.SetDefault("PROXY_MODE", "none")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
