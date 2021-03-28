package moodle

import (
	"fmt"
	_ "github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	baseURL = "https://moodle.nu.edu.kz"
)

type App struct {
	Client *http.Client
}

var app App

func (app *App) login(username string, password string) {
	fmt.Println(username, password)

	client := app.Client
	loginURL := fmt.Sprintf("%s/login/index.php", baseURL)
	data := url.Values{
		"username": {username},
		"password": {password},
	}
	response, _ := client.PostForm(loginURL, data)
	defer response.Body.Close()

	myResponse, _ := client.Get(fmt.Sprintf("%s/grade/report/overview/index.php", baseURL))
	defer myResponse.Body.Close()
}

func Init() {
	jar, _ := cookiejar.New(nil)
	app = App{
		Client: &http.Client{Jar: jar},
	}
}

func Login(username string, password string) {
	app.login(username, password)
}
