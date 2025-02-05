package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	courseapi "go-sheets/courseapi"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	WelcomeMsg                            = "Welcome to the Go-Sheets CLI! Type 'info' for a list of accepted commands, or 'exit' to quit."
	AssignmentInfoMsg                     = "Please input additional <due_date> (MM/DD/YY) and optional [<assignment_info>], space-delimited"
	ListAssignmentsCorrectUsageMsg        = "Usage: list-assignments <course_name>"
	CreateCourseCorrectUsageMsg           = "Usage: create-course <course_name> [<course_description>]"
	CreateAssignmentCorrectUsageMsg       = "Usage: create-assignment <course_name> <assignment_name>"
	RemoveCourseCorrectUsageMsg           = "Usage: remove-course <course_name>"
	RemoveAssignmentCorrectUsageMsg       = "Usage: remove-assignment <course_name> <assignment_number>"
	CreateAssignmentCorrectFieldsMsg      = "Fields: <due_date> (MM/DD/YY) [<assignment_info>]"
	AssignmentCourseDoesntExistMsg        = "Course for assignment doesn't exist"
	RemovalCourseDoesntExistMsg           = "Course to remove doesn't exist"
	ValidRemoveIndexMsg                   = "Removal index must be a valid integer"
	AssignmentRemovalCourseDoesntExistMsg = "Course for assignment removal doesn't exist"
	RemoveIndexOutOfBoundsMsg             = "Removal index out of bounds (check indices using `list-assignments <coursename>`)"
	CourseAlreadyExistsErrMsg             = "this course has already been added"

	UnsuccessfulLogSetupMsg      = "Unable to successfully setup logging file"
	UnsuccessfulSheetsSetupMsg   = "Unable to successfully connect to sheets service"
	UnsuccessfulCourseMapLoadMsg = "Unable to successfully load courseMap from sheets"

	UnsuccessfulCourseRemovalMsg      = "Unable to successfully remove course"
	UnsuccessfulAssignmentCreationMsg = "Unable to successfully create assignment for reason"
	UnsuccessfulAssignmentRemovalMsg  = "Unable to successfully remove assignment for reason"

	logFile   = "gosheets-cli.log"
	sheetName = "Sheet1"
)

type CourseMap = courseapi.CourseMap
type CourseItem = courseapi.CourseItem
type AssignmentList = courseapi.AssignmentList
type AssignmentItem = courseapi.AssignmentItem
type Service = sheets.Service

var (
	courseMap     CourseMap
	srv           *Service
	spreadsheetId string
)

func init() {
	_, err := initLog()
	if err != nil {
		log.Fatalf(UnsuccessfulLogSetupMsg+": %v", err)
	}

	srv, err = getSheetsService()
	if err != nil {
		log.Fatalf(UnsuccessfulSheetsSetupMsg+": %v", err)
	}

	spreadsheetId = getOrCreateSpreadsheet(srv, "Course Tracking Sheet")
	log.Printf("Using spreadsheet id: %s", spreadsheetId)

	courseMap, err = loadCourseMapFromSheets(srv, spreadsheetId)
	if err != nil {
		log.Fatalf(UnsuccessfulCourseMapLoadMsg+": %v", err)
	}
}

func main() {
	fmt.Println(WelcomeMsg)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
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
		case "list-courses":
			listCourses()
		case "list-assignments":
			if len(args) != 2 {
				fmt.Println(ListAssignmentsCorrectUsageMsg)
				continue
			}
			courseName := args[1]
			listAssignments(courseName)
		case "create-course":
			args = strings.SplitN(input, " ", 3)

			if len(args) < 2 || len(args) > 3 {
				fmt.Println(CreateAssignmentCorrectUsageMsg)
				continue
			}

			courseName := args[1]
			var courseDescription string
			if len(args) == 3 {
				courseDescription = args[2]
			}

			_, err := createCourse(courseName, courseDescription)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Course `%s` successfully created!\n", courseName)
			}
		case "create-assignment":
			if len(args) != 3 {
				fmt.Println(CreateAssignmentCorrectUsageMsg)
				continue
			}
			courseName := args[1]
			assignmentName := args[2]

			_, exists := courseMap[courseName]
			if !exists {
				fmt.Println(AssignmentCourseDoesntExistMsg)
				continue
			}

			fmt.Println(AssignmentInfoMsg)
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			args = strings.SplitN(input, " ", 2)
			if len(args) != 1 && len(args) != 2 {
				fmt.Println(CreateAssignmentCorrectFieldsMsg)
				continue
			} else {
				courseItem := courseMap[courseName]
				copy := courseItem.DeepCopy()

				var err error

				if len(args) == 1 {
					_, err = copy.Assignments.AddAssignment(assignmentName, args[0])
				} else {
					_, err = copy.Assignments.AddAssignment(assignmentName, args[0], args[1])
				}

				if err != nil {
					log.Printf(UnsuccessfulAssignmentCreationMsg+": %v", err)

					fmt.Printf("Unable to successfully add assignment `%s` to course `%s`\n", assignmentName, courseName)
				} else {
					err := updateCourseRow(srv, copy)

					if err != nil {
						log.Printf(UnsuccessfulAssignmentCreationMsg+": %v", err)

						fmt.Printf("Unable to successfully create assignment `%s`\n", assignmentName)
					} else {
						courseMap[courseName] = &copy

						fmt.Printf("Assignment `%s` successfully created!\n", assignmentName)
					}
				}
			}
		case "remove-course":
			if len(args) != 2 {
				fmt.Println(RemoveCourseCorrectUsageMsg)
				continue
			}

			courseName := args[1]
			_, exists := courseMap[courseName]
			if !exists {
				fmt.Println(RemovalCourseDoesntExistMsg)
				continue
			}

			err := removeCourseRow(courseName)
			if err != nil {
				log.Printf(UnsuccessfulCourseRemovalMsg+": %v", err)

				fmt.Printf("Unable to successfully remove course `%s`\n", courseName)
			} else {
				delete(courseMap, courseName)

				fmt.Printf("Course `%s` successfully removed!\n", courseName)
			}

		case "remove-assignment":
			if len(args) != 3 {
				fmt.Println(RemoveAssignmentCorrectUsageMsg)
				continue
			}

			courseName := args[1]
			removeIndex, err := strconv.Atoi(args[2])

			if err != nil {
				fmt.Println(ValidRemoveIndexMsg)
				continue
			}

			courseItem, exists := courseMap[courseName]
			if !exists {
				fmt.Println(AssignmentRemovalCourseDoesntExistMsg)
				continue
			}

			copy := courseItem.DeepCopy()

			_, err = copy.Assignments.RemoveAssignment(removeIndex - 1)
			if err != nil {
				fmt.Println(RemoveIndexOutOfBoundsMsg)
				continue
			}

			err = updateCourseRow(srv, copy)

			if err != nil {
				log.Printf(UnsuccessfulAssignmentRemovalMsg+": %v", err)

				fmt.Printf("Unable to successfully remove assignment number `%d`\n", removeIndex)
			} else {
				courseMap[courseName] = &copy

				fmt.Printf("Assignment number `%d` successfully removed!\n", removeIndex)
			}
		default:
			fmt.Println("Command not recognized")
		}

	}
}

func initLog() (*os.File, error) {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	return f, nil
}

func getSheetsService() (*Service, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("service-account.json"))

	return srv, err
}

func createSpreadsheet(srv *sheets.Service, title string) string {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}

	resp, err := srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		log.Fatalf("Unable to create spreadsheet: %v", err)
	}

	log.Println("New Spreadsheet Created!")
	log.Println("Spreadsheet Title:", resp.Properties.Title)
	log.Println("Spreadsheet ID:", resp.SpreadsheetId)

	saveToEnv("SPREADSHEET_ID", resp.SpreadsheetId)

	return resp.SpreadsheetId
}

func saveToEnv(key, value string) {
	envFile := ".env"
	f, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to write to .env file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		log.Fatalf("Failed to write .env variable: %v", err)
	}
}

func getOrCreateSpreadsheet(srv *sheets.Service, title string) string {
	err := godotenv.Load(".env")
	if err == nil {
		if sheetID, exists := os.LookupEnv("SPREADSHEET_ID"); exists {
			log.Println("Using existing Spreadsheet ID:", sheetID)
			return sheetID
		}
	}

	log.Println("No existing sheet found. Creating a new one...")
	return createSpreadsheet(srv, title)
}

func loadCourseMapFromSheets(srv *sheets.Service, spreadsheetId string) (CourseMap, error) {
	readRange := fmt.Sprintf("%s!A:B", sheetName) // A: Course Name, B: Course JSON
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to read data: %v", err)
	}

	courseMap := make(CourseMap)

	for _, row := range resp.Values {
		if len(row) < 2 {
			continue // Skip incomplete rows
		}

		courseName := row[0].(string)
		jsonData := row[1].(string)

		var course CourseItem
		err := json.Unmarshal([]byte(jsonData), &course)
		if err != nil {
			log.Printf("Skipping invalid JSON for course %s: %v\n", courseName, err)
			continue
		}

		courseMap[courseName] = &course
	}

	log.Println("CourseMap loaded from Google Sheets!")
	return courseMap, nil
}

func findCourseRow(srv *sheets.Service, courseName string) (string, error) {
	rangeToSearch := fmt.Sprintf("%s!A:A", sheetName)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, rangeToSearch).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve data: %v", err)
	}

	for i, row := range resp.Values {
		if len(row) > 0 && row[0] == courseName {
			return fmt.Sprintf("%s!A%d:B%d", sheetName, i+1, i+1), nil
		}
	}
	return "", fmt.Errorf("course not found")
}

func updateCourseRow(srv *sheets.Service, updatedCourse CourseItem) error {
	row, err := findCourseRow(srv, updatedCourse.Name)
	if err != nil {
		return fmt.Errorf("unable to find course to update")
	}

	jsonData, err := json.Marshal(updatedCourse)
	if err != nil {
		return fmt.Errorf("failed to encode CourseItem to JSON: %v", err)
	}

	values := [][]interface{}{
		{updatedCourse.Name, string(jsonData)}, // Update A (name) and B (JSON data)
	}
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, row, valueRange).
		ValueInputOption("RAW").
		Do()
	if err != nil {
		return fmt.Errorf("failed to update course: %v", err)
	}
	return nil
}

func removeCourseRow(courseName string) error {
	row, err := findCourseRow(srv, courseName)
	if err != nil {
		log.Printf("Failed to find course with name %s in sheet: %v\n", courseName, err)

		return fmt.Errorf("failed to find course to remove within sheet: %v", err)
	}

	_, err = srv.Spreadsheets.Values.Clear(spreadsheetId, row, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		log.Printf("Failed to update course with name %s in sheet: %v\n", courseName, err)

		return fmt.Errorf("failed to update course: %v", err)
	}
	return nil
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

func listCourses() {
	fmt.Print(courseMap.String())
}

func listAssignments(courseName string) {
	courseItem, exists := courseMap[courseName]

	if exists && hasAssignments(*courseItem) {
		fmt.Print(courseItem.DetailedString())
	} else if exists {
		fmt.Println(courseItem.String())
	} else {
		fmt.Println("Course for assignment doesn't exist")
	}
}

func addCourseToSheets(srv *sheets.Service, spreadsheetId string, course *CourseItem) (bool, error) {
	// Serialize the CourseItem to JSON
	jsonData, err := json.Marshal(course)
	if err != nil {
		return false, fmt.Errorf("failed to encode CourseItem to JSON: %v", err)
	}

	// Prepare the new row (A: Course Name, B: Course JSON)
	values := [][]interface{}{
		{course.Name, string(jsonData)},
	}

	writeRange := fmt.Sprintf("%s!A:B", sheetName)
	resp, err := srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").Do()

	if err != nil {
		return false, fmt.Errorf("failed to append new course: %v", err)
	}

	if resp.Updates.UpdatedRows < 1 {
		return false, fmt.Errorf("no rows were updated, course addition failed")
	}

	return true, nil
}

func createCourse(courseName string, courseDescription string) (bool, error) {
	_, exists := courseMap[courseName]

	if exists {
		return false, errors.New(CourseAlreadyExistsErrMsg)
	} else {
		if emptyDescription(courseDescription) {
			newCourse := CourseItem{Name: courseName, Course_Info: nil, Assignments: AssignmentList{}}

			// Try adding course to Google Sheets first
			success, err := addCourseToSheets(srv, spreadsheetId, &newCourse)
			if err != nil || !success {
				return false, fmt.Errorf("failed to add course `%s` to Sheets: %v", courseName, err)
			}

			courseMap[courseName] = &newCourse
		} else {
			newCourse := CourseItem{Name: courseName, Course_Info: &courseDescription, Assignments: AssignmentList{}}

			success, err := addCourseToSheets(srv, spreadsheetId, &newCourse)
			if err != nil || !success {
				return false, fmt.Errorf("failed to add course `%s` to Sheets: %v", courseName, err)
			}

			courseMap[courseName] = &newCourse
		}
		return true, nil
	}

}

func emptyDescription(name string) bool {
	return name == ""
}

func hasAssignments(course CourseItem) bool {
	return len(course.Assignments) > 0
}
