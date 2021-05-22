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

	moodleClient := moodle.Init()
	grades := moodleClient.GetGrades(viper.Get("USER.NAME").(string), viper.Get("USER.PASSWORD").(string))

	if grades.Success {
		for _, course := range grades.Courses {
			fmt.Println(course.Name)
			fmt.Println()
			fmt.Println(course.Grades)
			fmt.Println()
			fmt.Println()
		}
	}
}
