package ocw

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
)

var lRe, tRe, sRe, fnRe, lnRe *regexp2.Regexp

func init() {
	lRe, _ = regexp2.Compile(`(?<="link":").*?(?=["])`, 0)
	tRe, _ = regexp2.Compile(`(?<="title":").*?(?=["])`, 0)
	sRe, _ = regexp2.Compile(`(?<="sort":).*?(?=[,"])`, 0)
	fnRe, _ = regexp2.Compile(`(?<="first_name":").*?(?=[,"])`, 0)
	lnRe, _ = regexp2.Compile(`(?<="last_name":").*?(?=[,"])`, 0)
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

func extractSessions(data string) ([]Session, error) {
	var sessions []Session = []Session{}

	var matchedLink []string = findAllString(lRe, data)

	var title, sort string
	for _, link := range matchedLink {
		title, sort = getTitleSort(tRe, sRe, link, data)
		title, _ = toString(title)
		link = strings.Replace(link, "\\", "", -1)
		var v Session = Session{Title: title, Sort: sort, Link: link}
		sessions = append(sessions, v)
	}

	return sessions, nil
}

func getTitleSort(tRe, sRe *regexp2.Regexp, link, data string) (string, string) {
	var endIndex int = strings.Index(data, link)
	var startIndex int = strings.LastIndex(data[:endIndex], `"title"`)

	title, _ := tRe.FindStringMatch(data[startIndex:endIndex])
	sort, _ := sRe.FindStringMatch(data[startIndex:endIndex])
	return title.String(), sort.String()
}

func getTitle(data string) string {
	var title []string = findAllString(tRe, data)
	if len(title) > 0 {
		return title[0]
	}
	return ""
}

func getTeacher(data string) string {
	var firstName []string = findAllString(fnRe, data)
	var lastName []string = findAllString(lnRe, data)
	if len(firstName) > 0 && len(lastName) > 0 {
		return fmt.Sprintf("%s %s", firstName[0], lastName[0])
	} else if len(firstName) > 0 {
		return firstName[0]
	} else if len(lastName) > 0 {
		return lastName[0]
	}
	return ""
}
