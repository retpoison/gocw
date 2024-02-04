package main

import (
	"fmt"
	"gocw/ocw"
)

func menu() {
	var option int
	fmt.Println("1 - download course.")
	fmt.Println("0 - exit.")
	fmt.Printf("Choose an option:")
	fmt.Scanln(&option)

	switch option {
	case 1:
		downloadCourse()
	}
}

func downloadCourse() {
	var ID []string = getIDs()
	var courses []ocw.Course = getCourses(ID)

	var confirm string = getDownloadConfirm()
	if confirm == "n" {
		goto menu
	}

	downloadCourses(courses)
	fmt.Println("download done.")

menu:
	menu()
}
