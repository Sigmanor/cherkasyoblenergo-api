package config

import (
	"github.com/spf13/viper"
)

const AppVersion = "dev"

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`
	NewsURL    string `mapstructure:"NEWS_URL"`

	// Rate limiting
	RateLimitPerMinute int `mapstructure:"RATE_LIMIT_PER_MINUTE"`

	// Cache
	CacheTTLSeconds int `mapstructure:"CACHE_TTL_SECONDS"`

	// Authentication (optional)
	APIKey string `mapstructure:"API_KEY"`
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

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
