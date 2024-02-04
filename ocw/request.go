package ocw

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var header map[string][]string

func init() {
	header = http.Header{
		"Accept":           {"application/json, text/plain, */*"},
		"Accept-Encoding":  {"gzip, deflate, br"},
		"Accept-Language":  {"en-US,en;q=0.5"},
		"Connection":       {"keep-alive"},
		"Content-Length":   {"77"},
		"Content-Type":     {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Cookie":           {"_T=3yrg0o51usmccoow"},
		"DNT":              {"1"},
		"Host":             {"ocw.sharif.edu"},
		"Origin":           {"https://ocw.sharif.edu"},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"U":                {"null"},
		"User-Agent":       {"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"},
		"X-Requested-With": {"XMLHttpRequest"}}
}

func request(method, url string, body io.Reader) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getSessionsData(id int) (string, error) {
	var url string = "https://ocw.sharif.edu/api/v1/ocw/sessions"
	values := map[string]any{
		"limit":      "None",
		"order_type": "ASC",
		"course_id":  id}

	body, err := getBody(values)
	if err != nil {
		return "", err
	}

	data, err := request("POST", url, body)
	if err != nil {
		return "", err
	}
	return data, nil
}

func getCourseData(id int) (string, error) {
	var url string = "https://ocw.sharif.edu/api/v1/ocw/courses/users"
	values := map[string]any{"course_id": id, "role": []string{"teacher"}}

	body, err := getBody(values)
	if err != nil {
		return "", err
	}

	data, err := request("POST", url, body)
	if err != nil {
		return "", err
	}
	return data, nil
}

func getBody(data interface{}) (io.Reader, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(json_data), nil
}

func toString(str string) (string, error) {
	return strconv.Unquote(`"` + strings.Replace(str, `\t`, "", -1) + `"`)
}
