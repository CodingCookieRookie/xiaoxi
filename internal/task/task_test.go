package task

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	title := "Buy groceries"
	id := 1
	tk := NewTask(title, id)

	if tk.Title != title {
		t.Errorf("expected title %q, got %q", title, tk.Title)
	}

	if tk.ID != id {
		t.Errorf("expected id %d, got %d", id, tk.ID)
	}

	if tk.Completed {
		t.Error("new task should not be completed")
	}

	if tk.CreatedAt.IsZero() {
		t.Error("task CreatedAt should not be zero")
	}
}

func TestTaskList_Add(t *testing.T) {
	tl := NewTaskList()

	t1 := tl.Add("Task 1")
	t2 := tl.Add("Task 2")

	if len(tl.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tl.Tasks))
	}

	if t1.ID != 1 {
		t.Errorf("expected task ID 1, got %d", t1.ID)
	}

	if t2.ID != 2 {
		t.Errorf("expected task ID 2, got %d", t2.ID)
	}
}

func TestTaskList_Complete(t *testing.T) {
	tl := NewTaskList()
	tk := tl.Add("Test task")

	found, ok := tl.Complete(tk.ID)
	if !ok {
		t.Error("expected to find task")
	}

	if !found.Completed {
		t.Error("task should be marked as completed")
	}

	if !tl.Tasks[0].Completed {
		t.Error("task in list should be marked as completed")
	}
}

func TestTaskList_Complete_NotFound(t *testing.T) {
	tl := NewTaskList()

	_, ok := tl.Complete(999)
	if ok {
		t.Error("should not find non-existent task")
	}
}

func TestTaskList_Delete(t *testing.T) {
	tl := NewTaskList()
	tk := tl.Add("Test task")

	if !tl.Delete(tk.ID) {
		t.Error("expected to delete task")
	}

	if len(tl.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tl.Tasks))
	}
}

func TestTaskList_Delete_NotFound(t *testing.T) {
	tl := NewTaskList()

	if tl.Delete(999) {
		t.Error("should not delete non-existent task")
	}
}

func TestTaskList_ListPending(t *testing.T) {
	tl := NewTaskList()
	t1 := tl.Add("Task 1")
	tl.Add("Task 2")
	tl.Complete(t1.ID)

	pending := tl.ListPending()
	if len(pending) != 1 {
		t.Errorf("expected 1 pending task, got %d", len(pending))
	}

	if pending[0].Title != "Task 2" {
		t.Errorf("expected Task 2, got %s", pending[0].Title)
	}
}

func TestTaskList_ListCompleted(t *testing.T) {
	tl := NewTaskList()
	t1 := tl.Add("Task 1")
	tl.Add("Task 2")
	tl.Complete(t1.ID)

	completed := tl.ListCompleted()
	if len(completed) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(completed))
	}

	if completed[0].Title != "Task 1" {
		t.Errorf("expected Task 1, got %s", completed[0].Title)
	}
}

func TestTaskList_GetByID(t *testing.T) {
	tl := NewTaskList()
	tk := tl.Add("Test task")

	found, ok := tl.GetByID(tk.ID)
	if !ok {
		t.Error("expected to find task")
	}

	if found.Title != "Test task" {
		t.Errorf("expected title %q, got %q", "Test task", found.Title)
	}
}

func TestTaskList_GetByID_NotFound(t *testing.T) {
	tl := NewTaskList()

	_, ok := tl.GetByID(999)
	if ok {
		t.Error("should not find non-existent task")
	}
}

func TestNextID(t *testing.T) {
	tl := NewTaskList()

	if id := tl.NextID(); id != 1 {
		t.Errorf("expected first ID 1, got %d", id)
	}

	if id := tl.NextID(); id != 2 {
		t.Errorf("expected second ID 2, got %d", id)
	}
}
