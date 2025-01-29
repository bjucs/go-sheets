package go_sheets

import (
	"errors"
	"time"
)

const (
	// Format for inputted date strings is `MM/DD/YY`
	DateFormat               = "01/02/06"
	InvalidDateErrMsg        = "invalid date passed in (please use mm/dd/yy)"
	TooManyParamsErrMsg      = "too many additional info strings in addAssignment (optional: 1 max)"
	InvalidSliceRemoveErrMsg = "tried to remove an out-of-bounds slice index"
)

type CourseItem struct {
	Name               string
	Course_Description *string
	Assignments        AssignmentList
}

type AssignmentItem struct {
	Name  string
	Info  *string
	DueAt time.Time
}

type AssignmentList []AssignmentItem

func (l *AssignmentList) AddAssignment(name string, due string, info ...string) (bool, error) {
	dueDate, err := time.Parse(DateFormat, due)
	if err != nil {
		return false, errors.New(InvalidDateErrMsg)
	}

	var infoPtr *string
	if len(info) == 1 {
		infoPtr = &info[0]
	} else {
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
