package go_sheets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getDateFromTime(t time.Time) string {
	return t.Format(DateFormat)
}

func TestAddAssignment_SimpleAdd_Success(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "02/02/25"
	taskInfo := "Some new task"
	l.AddAssignment(taskName, dueDate, taskInfo)

	assert.Equal(t, l[0].Name, taskName)
	assert.Equal(t, getDateFromTime(l[0].DueAt), dueDate)
	assert.Equal(t, *l[0].Info, taskInfo)
}

func TestAddAssignment_InvalidDateFormat_Failure(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "21/02/25"
	taskInfo := "Some new task"

	_, err := l.AddAssignment(taskName, dueDate, taskInfo)
	assert.EqualErrorf(t, err, InvalidDateErrMsg, "Error should be: %v, got: %v")
}

func TestAddAssignment_TooManyInfoParams_Failure(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "01/12/25"
	taskInfo := "Some new task"
	invalidExtraTaskInfo := "Another task descriptor"

	_, err := l.AddAssignment(taskName, dueDate, taskInfo, invalidExtraTaskInfo)
	assert.EqualErrorf(t, err, TooManyParamsErrMsg, "Error should be: %v, got: %v")
}

func TestAddAssignment_MultipleSimpleAdds_Success(t *testing.T) {
	l := AssignmentList{}

	expectedTaskNames := []string{"task1", "task2", "task3"}
	expectedDueDates := []string{"02/02/25", "03/17/25", "08/24/25"}
	expectedTaskInfo := []string{"Very difficult!", "Not very difficult!"}

	for i := 0; i < len(expectedTaskNames); i++ {
		if i < len(expectedTaskNames)-1 {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i], expectedTaskInfo[i])
		} else {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i])
		}
	}

	for i, task := range l {
		assert.Equal(t, task.Name, expectedTaskNames[i])
		assert.Equal(t, getDateFromTime(l[i].DueAt), expectedDueDates[i])

		if task.Info != nil {
			assert.Equal(t, *task.Info, expectedTaskInfo[i])
		}
	}
}

func TestRemoveAssignment_SimpleRemove_Success(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "02/02/25"
	taskInfo := "Some new task"
	println(len(l))
	l.AddAssignment(taskName, dueDate, taskInfo)
	println(len(l))
	l.RemoveAssignment(0)
	println(len(l))

	assert.Equal(t, 0, len(l))
}
