package courseapi

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// Format for inputted date strings is `MM/DD/YY`
	DateFormat               = "01/02/06"
	InvalidDateErrMsg        = "invalid date passed in (please use mm/dd/yy)"
	TooManyParamsErrMsg      = "excess info strings passed to addAssignment"
	InvalidSliceRemoveErrMsg = "tried to remove an out-of-bounds slice index"
)

type CourseMap map[string]*CourseItem

func (cm CourseMap) String() string {
	if len(cm) == 0 {
		return "No courses available.\n"
	}

	result := ""
	for _, course := range cm {
		result += fmt.Sprintf("%s\n\n", course.String())
	}

	return strings.TrimSuffix(result, "\n")
}

type CourseItem struct {
	Name        string         `json:"name"`
	Course_Info *string        `json:"course_info,omitempty"`
	Assignments AssignmentList `json:"assignments"`
}

func (c CourseItem) String() string {
	courseInfoStr := ""
	if c.Course_Info != nil {
		courseInfoStr = fmt.Sprintf("\n%s", *c.Course_Info)
	}

	return fmt.Sprintf("Course: %s%s", c.Name, courseInfoStr)
}

func (c CourseItem) DetailedString() string {
	return fmt.Sprintf("%s\n\nAssignments:\n%s", c.String(), c.Assignments.String())
}

type AssignmentItem struct {
	Name  string    `json:"name"`
	Info  *string   `json:"info,omitempty"`
	DueAt time.Time `json:"due_at"`
}

func (a AssignmentItem) String() string {
	infoStr := ""
	if a.Info != nil {
		infoStr = fmt.Sprintf("\n%s", *a.Info)
	}

	return fmt.Sprintf("%s%s\nDue: %s", a.Name, infoStr, a.DueAt.Format(DateFormat))
}

type AssignmentList []AssignmentItem

func (l AssignmentList) String() string {
	if len(l) == 0 {
		return "No assignments available."
	}

	result := ""
	for i, item := range l {
		// Use `i+1` since we want our printed list to be 1-indexed
		result += fmt.Sprintf("%d. %s\n\n", i+1, item.String())
	}

	return strings.TrimSuffix(result, "\n")
}

func (l *AssignmentList) AddAssignment(name string, due string, info ...string) (bool, error) {
	dueDate, err := time.Parse(DateFormat, due)
	if err != nil {
		return false, errors.New(InvalidDateErrMsg)
	}

	var infoPtr *string
	if len(info) == 1 {
		infoPtr = &info[0]
	} else if len(info) > 1 {
		return false, errors.New(TooManyParamsErrMsg)
	}

	t := AssignmentItem{
		Name:  name,
		Info:  infoPtr,
		DueAt: dueDate,
	}

	*l = append(*l, t)
	return true, nil
}

func (l *AssignmentList) RemoveAssignment(index int) (bool, error) {
	if index < 0 || index >= len(*l) {
		return false, errors.New(InvalidSliceRemoveErrMsg)
	}

	*l = append((*l)[:index], (*l)[index+1:]...)
	return true, nil
}

func (l *AssignmentList) ViewAssignment(index int) (AssignmentItem, error) {
	if index < 0 || index >= len(*l) {
		return AssignmentItem{}, errors.New(InvalidSliceRemoveErrMsg)
	}

	return (*l)[index], nil
}
