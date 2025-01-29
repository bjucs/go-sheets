package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	WelcomeMsg                = "Welcome to the Go-Sheets CLI! Type 'info' for a list of accepted commands, or 'exit' to quit."
	AssignmentInfoMsg         = "Please input additional <due_date> (MM/DD/YY) and optional [<assignment_info>], space-delimited"
	CourseAlreadyExistsErrMsg = "this course is has already been added"
)

var (
	courseMap CourseMap
)

func init() {
	courseMap = make(CourseMap)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(WelcomeMsg)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "info":
			showInfo()
		case "create-course":
			if len(args) < 2 || len(args) > 3 {
				fmt.Println("Usage: create-course <course_name> [<course_description>]")
				continue
			}
			courseName := args[1]
			var courseDescription string
			if len(args) == 3 {
				courseDescription = args[2]
			}

			_, err := createCourse(courseName, courseDescription)
			if err != nil {
				fmt.Println("Course already exists")
			}
		case "create-assignment":
			if len(args) != 3 {
				fmt.Println("Usage: create-assignment <course_name> <assignment_name>")
				continue
			}
			courseName := args[1]
			assignmentName := args[2]

			_, exists := courseMap[courseName]
			if !exists {
				fmt.Println("Course for assignment doesn't exist")
				continue
			}

			fmt.Println(AssignmentInfoMsg)
			fmt.Print("> ")
			scanner.Scan()
			input := scanner.Text()
			input = strings.TrimSpace(input)

			args := strings.Fields(input)
			if len(args) != 1 && len(args) != 2 {
				fmt.Println("Fields: <due_date> (MM/DD/YY) [<assignment_info>]")
				continue
			} else {
				courseItem := courseMap[courseName]

				var err error

				if len(args) == 1 {
					_, err = courseItem.Assignments.AddAssignment(assignmentName, args[0])
				} else {
					_, err = courseItem.Assignments.AddAssignment(assignmentName, args[0], args[1])
				}

				if err != nil {
					println(err)
				}
			}
		}

	}
}

func showInfo() {
	info := `Available Commands:
    
create-course <course_name>
    - User can optionally include an additional <class_description> parameter
    
create-assignment <course_name> <assignment_name>
    - User will be prompted for other info, such as:
        - due_date (required)
        - info (optional notes)
    
list-courses
    - Lists all available courses
    
list-assignments <course_name>
    - Lists all assignments for the specified course
    
remove-course <course_number>
    - Removes the course at the specified 1-based index
    
remove-assignment <assignment_number>
    - Removes the assignment at the specified 1-based index
    - Use the indices provided in the list-courses or list-assignments commands`

	fmt.Println(info)
}

func createCourse(courseName string, courseDescription string) (bool, error) {
	_, exists := courseMap[courseName]

	if exists {
		return false, errors.New(CourseAlreadyExistsErrMsg)
	} else {
		if emptyDescription(courseDescription) {
			courseMap[courseName] = CourseItem{courseName, nil, AssignmentList{}}
		} else {
			courseMap[courseName] = CourseItem{courseName, &courseDescription, AssignmentList{}}
		}
		return true, nil
	}

}

func emptyDescription(name string) bool {
	return name == ""
}
