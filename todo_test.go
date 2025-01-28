package go_sheets

import "testing"

func TestAddAssignment_SimpleAdd_Success(t *testing.T) {
	l := List{}

	taskName := "New Task"
	l.AddAssignment(taskName)

	if l[0].Name != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Name)
	}
}

func TestAddAssignment_MultipleAdds_Success(t *testing.T) {
	l := List{}

	taskName1 := "New Task 1"
	taskName2 := "New Task 2"
	taskName3 := "New Task 3"
	l.AddAssignment(taskName1)
	l.AddAssignment(taskName2)
	l.AddAssignment(taskName3)

	expectedTasks := []string{taskName1, taskName2, taskName3}

	for i, task := range l {
		if task.Name != expectedTasks[i] {
			t.Errorf("Task at index %d: Expected %q, got %q instead", i, expectedTasks[i], task.Name)
		}
	}
}
