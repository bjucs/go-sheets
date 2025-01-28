package go_sheets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddAssignment_SimpleAdd_Success(t *testing.T) {
	l := List{}

	taskName := "New Task"
	dueDate := "02/02/25"
	taskInfo := "Some new task"
	l.AddAssignment(taskName, dueDate, taskInfo)

	assert.Equal(t, l[0].Name, taskName)
	assert.Equal(t, getDateFromTime(l[0].DueAt), dueDate)
	assert.Equal(t, *l[0].Info, taskInfo)

}

func getDateFromTime(t time.Time) string {
	return t.Format(DateFormat)
}

// func TestAddAssignment_MultipleAdds_Success(t *testing.T) {
// 	l := List{}

// 	taskName1 := "New Task 1"
// 	taskName2 := "New Task 2"
// 	taskName3 := "New Task 3"
// 	l.AddAssignment(taskName1)
// 	l.AddAssignment(taskName2)
// 	l.AddAssignment(taskName3)

// 	expectedTasks := []string{taskName1, taskName2, taskName3}

// 	for i, task := range l {
// 		if task.Name != expectedTasks[i] {
// 			t.Errorf("Task at index %d: Expected %q, got %q instead", i, expectedTasks[i], task.Name)
// 		}
// 	}
// }
