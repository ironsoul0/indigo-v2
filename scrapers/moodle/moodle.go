package moodle

import (
	"errors"
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
	baseURL   = "https://moodle.nu.edu.kz"
	csrfToken = "logintoken"
)

type App struct {
	client *http.Client
}

type Grade struct {
	Name       string
	Grade      string
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

func (app *App) extractToken() (string, error) {
	response, _ := app.client.Get(baseURL)
	pageDoc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if loginToken, ok := pageDoc.Find(fmt.Sprintf("input[name|='%s']", csrfToken)).First().Attr("value"); ok {
		return loginToken, nil
	}

	return "", errors.New("no token was found")
}

func (app *App) login(username string, password string) LoginResponse {
	extractedToken, err := app.extractToken()

	if err != nil {
		return LoginResponse{RequestFailure: true}
	}

	loginURL := fmt.Sprintf("%s/login/index.php", baseURL)
	data := url.Values{
		"username": {username},
		"password": {password},
		csrfToken:  {extractedToken},
	}

	response, err := app.client.PostForm(loginURL, data)
	if err != nil {
		return LoginResponse{RequestFailure: true}
	}
	defer response.Body.Close()

	if !strings.Contains(response.Request.URL.String(), "my") {
		return LoginResponse{InvalidCredentials: true}
	}

	return LoginResponse{}
}

// TODO: Ban grades with names like Attendance, etc..
func (app *App) parseCourse(courseLink string, courseName string) []Grade {
	coursePage, err := app.client.Get(courseLink)
	if err != nil {
		return make([]Grade, 0)
	}
	defer coursePage.Body.Close()

	courseDoc, _ := goquery.NewDocumentFromReader(coursePage.Body)
	gradeRows := courseDoc.Find("tbody").First().Find("tr")
	grades := make([]Grade, 0)

	gradeRows.Each(func(i int, s *goquery.Selection) {
		grade := strings.TrimSpace(s.Find(".column-grade").First().Text())
		if len(grade) == 0 {
			return
		}

		name := strings.TrimSpace(s.Find(".column-itemname").First().Text())
		gradeRange := strings.TrimSpace(s.Find(".column-range").First().Text())
		percentage := strings.TrimSpace(s.Find(".column-percentage").First().Text())

		grades = append(grades, Grade{Name: name, Grade: grade, Range: gradeRange, Percentage: percentage})
	})

	return grades
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
			coursesChannel <- Course{Name: courseEntry.Name, Grades: app.parseCourse(courseEntry.Link, courseEntry.Name)}
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
	jar, _ := cookiejar.New(nil)

	return &App{
		client: &http.Client{Jar: jar, Timeout: 3 * time.Second},
	}
}
