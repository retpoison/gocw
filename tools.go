package main

import (
	"fmt"
	"gocw/ocw"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/juju/ratelimit"
	"github.com/schollz/progressbar/v3"
)

func getIDs() []string {
	var ID string
	fmt.Println("hint: 65")
	fmt.Println("you can enter multiple IDs separated by comma. for example: 34,25,15")
	fmt.Print("enter course id: ")
	fmt.Scanln(&ID)

	return strings.Split(ID, ",")
}

func getCourses(ID []string) []ocw.Course {
	var course ocw.Course
	var courses []ocw.Course
	var checkCourse func(course ocw.Course)
	checkCourse = func(course ocw.Course) {
		if course.IsCourseExists() {
			courses = append(courses, course)
			fmt.Println(course.ID, course.Title)
		} else {
			var id int
			fmt.Println("course with id", course.ID, `doesn't exist.`)
			fmt.Printf("enter again (enter 0 to skip): ")
			fmt.Scanln(&id)
			if id != 0 {
				checkCourse(ocw.Course{ID: id})
			}
		}
	}

	for _, eachID := range ID {
		id, _ := strconv.Atoi(eachID)
		course = ocw.Course{ID: id}
		checkCourse(course)
	}
	return courses
}

func getDownloadConfirm() string {
	var confirm string
	for confirm != "y" && confirm != "n" {
		fmt.Printf("download all the courses? (y/n) ")
		fmt.Scanln(&confirm)
	}
	return confirm
}

func downloadCourses(courses []ocw.Course) {
	var unableToDownload = []ocw.Session{}

	for _, c := range courses {
		foldername, err := c.GetFolderName()
		if err != nil {
			fmt.Println(err)
			foldername = strconv.Itoa(c.ID)
		}

		err = makeFolder(foldername)
		if err != nil {
			fmt.Println("failed making folder.")
			fmt.Println(err)
			fmt.Println("skipping", c.ID, c.Title+".")
			continue
		}

		err = c.InitSessions()
		if err != nil {
			fmt.Println("an error occurred during getting sessions.")
			fmt.Println(err)
			return
		}

		var speedLimit int64 = getSpeedLimit() * 1024
		unable := downloadSessions(c.Sessions, foldername, speedLimit)
		if len(unable) > 0 {
			fmt.Println("\ntrying again to download the files that encountered a problem.")
		}
		unable = downloadSessions(unable, foldername, speedLimit)
		if len(unable) > 0 {
			unableToDownload = append(unableToDownload, unable...)
		}
	}

	if len(unableToDownload) > 0 {
		fmt.Println("a problem occurred when downloading these files:")
		for _, v := range unableToDownload {
			fmt.Printf("%s %s\nhttp://ocw.sharif.edu%s\n",
				v.Sort, v.Title, v.Link)
		}
	}
}

func downloadSessions(sessions []ocw.Session, foldername string, speedLimit int64) []ocw.Session {
	var unable = []ocw.Session{}
	for _, v := range sessions {
		fmt.Println("downloading", v.Sort, v.Title)
		fmt.Println("http://ocw.sharif.edu" + v.Link)

		err := downloadFile("http://ocw.sharif.edu"+v.Link,
			foldername+"/"+v.Sort+" - "+v.Title+path.Ext(v.Link),
			speedLimit)
		if err != nil {
			fmt.Printf("\nunable to download\n%s %s\nhttp://ocw.sharif.edu%s\n%s\n",
				v.Sort, v.Title, v.Link, err)
			unable = append(unable, v)
		}
	}
	return unable
}

func makeFolder(folderPath string) error {
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func downloadFile(url, filepath string, speedLimit int64) error {
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
	if speedLimit == 0 {
		_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	} else {
		bucket := ratelimit.NewBucketWithRate(float64(speedLimit), speedLimit)
		_, err = io.Copy(io.MultiWriter(out, bar),
			ratelimit.Reader(resp.Body, bucket))
	}
	if err != nil {
		return err
	}

	return nil
}

func getSpeedLimit() int64 {
	var limit int64 = 0
	for ok := true; ok; ok = (limit < 0) {
		fmt.Printf("enter the download speed limit in KB: (0 for default) ")
		fmt.Scanln(&limit)
	}
	return limit
}
