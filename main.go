package main

import (
	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.ReadInConfig()

	moodle.Init()
	moodle.Login(viper.Get("USER.NAME").(string), viper.Get("USER.PASSWORD").(string))
}
