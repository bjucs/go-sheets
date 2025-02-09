package courseapi

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

	assert.Equal(t, taskName, l[0].Name)
	assert.Equal(t, dueDate, getDateFromTime(l[0].DueAt))
	assert.Equal(t, taskInfo, *l[0].Info)
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
		if i < len(expectedTaskInfo) {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i], expectedTaskInfo[i])
		} else {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i])
		}
	}

	for i, task := range l {
		assert.Equal(t, task.Name, expectedTaskNames[i])
		assert.Equal(t, getDateFromTime(l[i].DueAt), expectedDueDates[i])

		if task.Info != nil {
			assert.Equal(t, expectedTaskInfo[i], *task.Info)
		}
	}

	assert.Equal(t, 3, len(l))
}

func TestAddAssignment_SortedInsertion_Success(t *testing.T) {
	l := AssignmentList{}

	expectedTaskNames := []string{"Task 1", "Task 2", "Task 3", "Task 4"}
	expectedDueDates := []string{"01/05/25", "02/20/25", "03/10/25", "05/15/25"}

	// Insert tasks in unsorted order
	l.AddAssignment("Task 3", "03/10/25")
	l.AddAssignment("Task 1", "01/05/25")
	l.AddAssignment("Task 4", "05/15/25")
	l.AddAssignment("Task 2", "02/20/25")

	for i, task := range l {
		assert.Equal(t, expectedTaskNames[i], task.Name)
		assert.Equal(t, expectedDueDates[i], getDateFromTime(l[i].DueAt))
	}

	assert.Equal(t, len(expectedTaskNames), len(l))
}

func TestRemoveAssignment_SimpleRemove_Success(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "02/02/25"
	taskInfo := "Some new task"

	l.AddAssignment(taskName, dueDate, taskInfo)
	l.RemoveAssignment(0)

	assert.Equal(t, 0, len(l))
}

func TestRemoveAssignment_EmptySlice_Failure(t *testing.T) {
	l := AssignmentList{}
	_, err := l.RemoveAssignment(0)

	assert.EqualErrorf(t, err, InvalidSliceRemoveErrMsg, "Error should be: %v, got: %v")
}

func TestRemoveAssignment_SliceOutOfBounds_Failure(t *testing.T) {
	l := AssignmentList{}

	taskName := "New Task"
	dueDate := "02/02/25"
	taskInfo := "Some new task"

	l.AddAssignment(taskName, dueDate, taskInfo)
	_, err := l.RemoveAssignment(3)

	assert.EqualErrorf(t, err, InvalidSliceRemoveErrMsg, "Error should be: %v, got: %v")
}

func TestRemoveAssignment_MultipleSimpleRemoves_Success(t *testing.T) {
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

	l.RemoveAssignment(0)
	l.RemoveAssignment(0)

	assert.Equal(t, 1, len(l))
	assert.Equal(t, "task3", l[0].Name)
}

func TestViewAssignment_MultipleSimpleViews_Success(t *testing.T) {
	l := AssignmentList{}

	expectedTaskNames := []string{"task1", "task2", "task3"}
	expectedDueDates := []string{"02/02/25", "03/17/25", "08/24/25"}
	expectedTaskInfo := []string{"Very difficult!", "Not very difficult!"}

	for i := 0; i < len(expectedTaskNames); i++ {
		if i < len(expectedTaskInfo) {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i], expectedTaskInfo[i])
		} else {
			l.AddAssignment(expectedTaskNames[i], expectedDueDates[i])
		}
	}

	for i := range l {
		task := l[i]
		assert.Equal(t, task.Name, expectedTaskNames[i])
		assert.Equal(t, getDateFromTime(l[i].DueAt), expectedDueDates[i])

		if task.Info != nil {
			assert.Equal(t, expectedTaskInfo[i], *task.Info)
		}
	}
}

func TestCourseItem_String_NoCourseInfo_NoAssignments_Success(t *testing.T) {
	course := CourseItem{
		Name: "Course 1",
	}

	expected := "Course: Course 1"
	assert.Equal(t, expected, course.String())
}

func TestCourseItem_String_WithCourseInfo_NoAssignments_Success(t *testing.T) {
	courseInfo := "Some course info"
	course := CourseItem{
		Name:        "Course 1",
		Course_Info: &courseInfo,
	}

	expected := "Course: Course 1\nSome course info"
	assert.Equal(t, expected, course.String())
}

func TestCourseItem_DetailedString_WithAssignments_Success(t *testing.T) {
	courseInfo := "Some course info"
	assignments := AssignmentList{
		{
			Name:  "Task 1",
			Info:  nil,
			DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	course := CourseItem{
		Name:        "Course 1",
		Course_Info: &courseInfo,
		Assignments: assignments,
	}

	expected := "Course: Course 1\nSome course info\n\nAssignments:\n1. Task 1\nDue: 02/02/25\n"
	assert.Equal(t, expected, course.DetailedString())
}

func TestCourseItem_DetailedString_WithMultipleAssignments_Success(t *testing.T) {
	courseInfo := "Course description"
	assignments := AssignmentList{
		{
			Name:  "Task 1",
			Info:  nil,
			DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:  "Task 2",
			Info:  nil,
			DueAt: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:  "Task 3",
			Info:  nil,
			DueAt: time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
		},
	}

	course := CourseItem{
		Name:        "Course 1",
		Course_Info: &courseInfo,
		Assignments: assignments,
	}

	expected := "Course: Course 1\nCourse description\n\nAssignments:\n1. Task 1\nDue: 02/02/25\n\n2. Task 2\nDue: 03/17/25\n\n3. Task 3\nDue: 08/24/25\n"
	assert.Equal(t, expected, course.DetailedString())
}

func TestCourseMap_String_EmptyMap_Success(t *testing.T) {
	cm := CourseMap{}
	expected := "No courses available.\n"
	assert.Equal(t, expected, cm.String())
}

func TestCourseMap_String_WithCourses_Success(t *testing.T) {
	courseInfo := "Some course info"
	assignments := AssignmentList{
		{
			Name:  "Task 1",
			Info:  nil,
			DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	courseMap := CourseMap{
		"Course 1": {
			Name:        "Course 1",
			Course_Info: &courseInfo,
			Assignments: assignments,
		},
	}

	expected := "Course: Course 1\nSome course info\n"
	print(courseMap.String())
	assert.Equal(t, expected, courseMap.String())
}

func TestAssignmentItem_String_NoInfo_Success(t *testing.T) {
	assignment := AssignmentItem{
		Name:  "Task 1",
		Info:  nil,
		DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
	}

	expected := "Task 1\nDue: 02/02/25"
	assert.Equal(t, expected, assignment.String())
}

func TestAssignmentItem_String_WithInfo_Success(t *testing.T) {
	taskInfo := "Some task info"
	assignment := AssignmentItem{
		Name:  "Task 1",
		Info:  &taskInfo,
		DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
	}

	expected := "Task 1\nSome task info\nDue: 02/02/25"
	assert.Equal(t, expected, assignment.String())
}

func TestAssignmentList_String_EmptyList_Success(t *testing.T) {
	assignmentList := AssignmentList{}
	expected := "No assignments available."
	assert.Equal(t, expected, assignmentList.String())
}

func TestAssignmentList_String_WithAssignments_Success(t *testing.T) {
	assignments := AssignmentList{
		{
			Name:  "Task 1",
			Info:  nil,
			DueAt: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:  "Task 2",
			Info:  nil,
			DueAt: time.Date(2025, 3, 17, 0, 0, 0, 0, time.UTC),
		},
	}

	expected := "1. Task 1\nDue: 02/02/25\n\n2. Task 2\nDue: 03/17/25\n"
	assert.Equal(t, expected, assignments.String())
}
