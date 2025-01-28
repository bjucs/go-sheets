package go_sheets

import (
	"time"
)

type AssignmentItem struct {
	Name        string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []AssignmentItem

// Callable methods for go-sheets
func (l *List) AddAssignment(name string) {
	t := AssignmentItem{
		Name:        name,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}
