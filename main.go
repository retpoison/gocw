package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type session struct {
	title string
	sort  string
	link  string
}

func download(se session, folderName string) error {
	var url string = "http://ocw.sharif.edu" + se.link
	var filepath string = folderName + "/" + se.sort + " - " + se.title + path.Ext(se.link)
	var err error

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	bar := progressbar.DefaultBytes(resp.ContentLength)
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getData(courseID string) string {
	var err error
	var url string = "https://ocw.sharif.edu/api/v1/ocw/sessions"
	var values map[string]string = map[string]string{
		"limit":      "None",
		"order_type": "ASC",
		"course_id":  courseID}

	json_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header{
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
		"Referer":          {"https://ocw.sharif.edu/course/id/57"},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"T":                {"3yrg0o51usmccoow"},
		"U":                {"null"},
		"User-Agent":       {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"},
		"X-Requested-With": {"XMLHttpRequest"}}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func findAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func getTitleSort(tRe, sRe *regexp2.Regexp, link, data string) (string, string) {
	var endIndex int = strings.Index(data, link)
	var startIndex int = strings.LastIndex(data[:endIndex], `"title"`)

	title, _ := tRe.FindStringMatch(data[startIndex:endIndex])
	sort, _ := sRe.FindStringMatch(data[startIndex:endIndex])
	return title.String(), sort.String()
}

func extractData(data string) []session {
	var err error

	var linkPattern string = `(?<="link":").*?(?=["])`
	lRe, err := regexp2.Compile(linkPattern, 0)
	if err != nil {
		log.Fatal(err)
	}

	var tPattern string = `(?<="title":").*?(?=["])`
	tRe, err := regexp2.Compile(tPattern, 0)
	if err != nil {
		log.Fatal(err)
	}

	var sPattern string = `(?<="sort":).*?(?=[,"])`
	sRe, err := regexp2.Compile(sPattern, 0)
	if err != nil {
		log.Fatal(err)
	}

	var matchedLink []string = findAllString(lRe, data)

	var sessions []session
	var title, sort string
	for _, link := range matchedLink {
		title, sort = getTitleSort(tRe, sRe, link, data)
		title, err = strconv.Unquote(`"` + strings.Replace(title, `\t`, "", -1) + `"`)
		link = strings.Replace(link, "\\", "", -1)
		var v session = session{title: title, sort: sort, link: link}
		sessions = append(sessions, v)
	}

	return sessions
}

func getVideos(courseID string) []session {
	var stringData string = getData(courseID)
	var sessions []session = extractData(stringData)

	return sessions
}

func main() {
	var ID string
	fmt.Println("hint: 65")
	fmt.Println("you can enter multiple IDs separated by coma. for example: 34,25,15")
	fmt.Print("enter course id: ")
	fmt.Scanln(&ID)

	var sessions []session
	var confirmation string
	var err error
	for _, id := range strings.Split(ID, ",") {

		sessions = getVideos(id)

		fmt.Printf("enter Y/y for downloading the coures: ")
		fmt.Scanln(&confirmation)

		if confirmation == "y" || confirmation == "Y" {
			if err = os.MkdirAll(id, os.ModePerm); err != nil {
				log.Fatal(err)
			}

			for _, v := range sessions {
				fmt.Println("downloading", v.sort, v.title)
				fmt.Println("ocw.sharif.edu/" + v.link)
				err = download(v, id)
				if err != nil {
					fmt.Println("unable to download")
					fmt.Println(v.sort)
					fmt.Println(v.title)
					fmt.Println(v.link)
					fmt.Println(err, "\n")
				}
			}
		}
	}
}
