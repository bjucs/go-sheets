package go_sheets

import (
	"fmt"
	"time"
)

const (
	// Format for inputted date strings is `MM/DD/YY`
	DateFormat          = "01/02/06"
	InvalidDateErrMsg   = "invalid date passed in (please use mm/dd/yy)"
	TooManyParamsErrMsg = "too many additional info strings in addAssignment (optional: 1 max)"
)

type AssignmentItem struct {
	Name  string
	Info  *string
	DueAt time.Time
}

type List []AssignmentItem

// Callable methods for go-sheets
func (l *List) AddAssignment(name string, due string, info ...string) (bool, error) {
	dueDate, err := time.Parse(DateFormat, due)
	if err != nil {
		return false, fmt.Errorf(InvalidDateErrMsg)
	}

	var infoPtr *string
	if len(info) == 1 {
		infoPtr = &info[0]
	} else {
		return false, fmt.Errorf(TooManyParamsErrMsg)
	}

	t := AssignmentItem{
		Name:  name,
		Info:  infoPtr,
		DueAt: dueDate,
	}

	*l = append(*l, t)
	return true, nil
}
