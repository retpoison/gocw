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
		"Content-Type":     {"application/x-www-form-urlencoded; charset=UTF-8"},
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

func getTeacherData(id int) (string, error) {
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
