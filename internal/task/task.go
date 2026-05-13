package task

import (
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskList struct {
	Tasks         []Task `json:"tasks"`
	NextIDCounter int    `json:"next_id"`
}

func NewTaskList() *TaskList {
	return &TaskList{
		Tasks:         []Task{},
		NextIDCounter: 1,
	}
}

func (tl *TaskList) NextID() int {
	id := tl.NextIDCounter
	tl.NextIDCounter++
	return id
}

func (tl *TaskList) SyncNextID() {
	for _, t := range tl.Tasks {
		if t.ID >= tl.NextIDCounter {
			tl.NextIDCounter = t.ID + 1
		}
	}
}

func NewTask(title string, id int) Task {
	return Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

func (tl *TaskList) Add(title string) Task {
	t := NewTask(title, tl.NextID())
	tl.Tasks = append(tl.Tasks, t)
	return t
}

func (tl *TaskList) Complete(id int) (Task, bool) {
	for i, t := range tl.Tasks {
		if t.ID == id {
			tl.Tasks[i].Completed = true
			return tl.Tasks[i], true
		}
	}
	return Task{}, false
}

func (tl *TaskList) Delete(id int) bool {
	for i, t := range tl.Tasks {
		if t.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (tl *TaskList) GetByID(id int) (Task, bool) {
	for _, t := range tl.Tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}

func (tl *TaskList) ListPending() []Task {
	var pending []Task
	for _, t := range tl.Tasks {
		if !t.Completed {
			pending = append(pending, t)
		}
	}
	return pending
}

func (tl *TaskList) ListCompleted() []Task {
	var completed []Task
	for _, t := range tl.Tasks {
		if t.Completed {
			completed = append(completed, t)
		}
	}
	return completed
}
