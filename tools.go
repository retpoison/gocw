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

		for _, v := range c.Sessions {
			fmt.Println("downloading", v.Sort, v.Title)
			fmt.Println("ocw.sharif.edu/" + v.Link)

			err = downloadFile("http://ocw.sharif.edu"+v.Link,
				foldername+"/"+v.Sort+" - "+v.Title+path.Ext(v.Link))
			if err != nil {
				fmt.Printf("unable to download\n%s\n%s\n%s\n%s", v.Sort, v.Title, v.Link, err)
			}
		}
	}
}

func makeFolder(folderPath string) error {
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func downloadFile(url, filepath string) error {
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
