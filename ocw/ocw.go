package ocw

import (
	"fmt"
	"strings"
)

type Course struct {
	ID       int
	Title    string
	Sessions []Session
}

type Session struct {
	Title string
	Sort  string
	Link  string
}

func (c *Course) IsCourseExists() bool {
	var url string = "https://ocw.sharif.edu/api/v1/ocw/course/get"
	var values map[string]int = map[string]int{"id": c.ID}

	body, err := getBody(values)
	if err != nil {
		return false
	}

	data, err := request("POST", url, body)
	if strings.Contains(data, `"http_code":404`) || err != nil {
		return false
	}

	c.Title, _ = toString(getTitle(data))
	return true
}

func (c *Course) GetFolderName() (string, error) {
	data, err := getCourseData(c.ID)
	if err != nil {
		return "", err
	}

	var title, teacher string = getTitle(data), getTeacher(data)
	title, _ = toString(title)
	teacher, _ = toString(teacher)
	return fmt.Sprintf("%s - %s", title, teacher), nil
}

func (c *Course) InitSessions() error {
	data, err := getSessionsData(c.ID)
	if err != nil {
		return err
	}

	sessions, err := extractSessions(data)
	if err != nil {
		return err
	}
	c.Sessions = sessions
	return nil
}
