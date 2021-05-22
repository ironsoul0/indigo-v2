package main

import (
	"fmt"

	"github.com/ironsoul0/indigo-v2/scrapers/moodle"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.ReadInConfig()

	app := moodle.Init()
	grades := app.GetGrades(viper.Get("USER.NAME").(string), viper.Get("USER.PASSWORD").(string))

	if grades.Success {
		fmt.Println(grades.Courses)
	}
}
