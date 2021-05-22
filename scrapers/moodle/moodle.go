package moodle

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://moodle.nu.edu.kz"
)

type App struct {
	client *http.Client
}

type Grade struct {
	Name       string
	Value      string
	Range      string
	Percentage string
}

type Course struct {
	Name   string
	Grades []Grade
}

type LoginResponse struct {
	RequestFailure     bool
	InvalidCredentials bool
}

type MoodleResponse struct {
	Success bool
	Courses []Course
	LoginResponse
}

func (app *App) login(username string, password string) LoginResponse {
	loginURL := fmt.Sprintf("%s/login/index.php", baseURL)
	data := url.Values{
		"username": {username},
		"password": {password},
	}

	response, err := app.client.PostForm(loginURL, data)
	if err != nil {
		return LoginResponse{RequestFailure: true}
	}
	defer response.Body.Close()

	if !strings.Contains(response.Header.Get("Location"), "testsession") {
		return LoginResponse{InvalidCredentials: true}
	}

	return LoginResponse{}
}

func (app *App) parseCourse(courseLink string) []Grade {
	return make([]Grade, 0)
}

func (app *App) GetGrades(username string, password string) MoodleResponse {
	loginResponse := app.login(username, password)
	if loginResponse.InvalidCredentials || loginResponse.RequestFailure {
		return MoodleResponse{Success: false, LoginResponse: loginResponse}
	}

	gradesPage, err := app.client.Get(fmt.Sprintf("%s/grade/report/overview/index.php", baseURL))
	if err != nil {
		return MoodleResponse{Success: false, LoginResponse: loginResponse}
	}
	defer gradesPage.Body.Close()

	gradesDoc, _ := goquery.NewDocumentFromReader(gradesPage.Body)
	gradeRows := gradesDoc.Find("tbody").First().Find("tr")

	type CourseEntry struct {
		Name string
		Link string
	}

	courseEntries := make([]CourseEntry, 0)

	gradeRows.Each(func(i int, s *goquery.Selection) {
		courseLink := s.Find("a").First()
		if value, ok := courseLink.Attr("href"); ok && len(value) > 0 {
			courseEntries = append(courseEntries, CourseEntry{Link: value, Name: courseLink.Text()})
		}
	})

	var wg sync.WaitGroup
	coursesChannel := make(chan Course)

	for _, courseEntry := range courseEntries {
		wg.Add(1)

		go func(courseEntry CourseEntry) {
			defer wg.Done()
			coursesChannel <- Course{Name: courseEntry.Name, Grades: app.parseCourse(courseEntry.Link)}
		}(courseEntry)
	}

	go func() {
		wg.Wait()
		close(coursesChannel)
	}()

	courses := make([]Course, 0)
	for course := range coursesChannel {
		courses = append(courses, course)
	}

	return MoodleResponse{Success: true, Courses: courses}
}

func Init() *App {
	var app App
	jar, _ := cookiejar.New(nil)
	app = App{
		client: &http.Client{Jar: jar, Timeout: 3 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}},
	}
	return &app
}
