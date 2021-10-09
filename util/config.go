package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	UserName            string        `mapstructure:"USER_NAME"`
	UserPassword        string        `mapstructure:"USER_PASSWORD"`
	BotToken            string        `mapstructure:"BOT_TOKEN"`
	DBPath              string        `mapstructure:"DB_PATH"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("could not read config", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("could not read config", err)
	}

	return config
}
