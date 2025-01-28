package go_sheets

import (
	"fmt"
	"time"
)

const (
	// Format for inputted date strings is `MM/DD/YY`
	DateFormat = "01/02/06"
)

type AssignmentItem struct {
	Name  string
	Info  *string
	DueAt time.Time
}

type List []AssignmentItem

// Callable methods for go-sheets
func (l *List) AddAssignment(name string, due string, info ...string) {
	dueDate, err := time.Parse(DateFormat, due)
	if err != nil {
		fmt.Println("Invalid date passed in (please use mm/dd/yy):", err)
		return
	}

	var infoPtr *string
	if len(info) == 0 || len(info) > 1 {
		// Invalid amount of info passed in
	} else if len(info) == 1 {
		infoPtr = &info[0]
	}

	t := AssignmentItem{
		Name:  name,
		Info:  infoPtr,
		DueAt: dueDate,
	}

	*l = append(*l, t)
}
