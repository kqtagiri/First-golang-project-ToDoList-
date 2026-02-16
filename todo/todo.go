package todo

import (
	"ToDo/database"
	"time"
)

type Task struct {
	Title      string     `json:"title"`
	Desc       string     `json:"description"`
	Priority   int        `json:"priority"`
	Status     bool       `json:"status"`
	DateCreate time.Time  `json:"dateCreate"`
	DateCompl  *time.Time `json:"dateComplete"`
}

type TaskDTO struct {
	Title    string `json:"title"`
	Desc     string `json:"description"`
	Priority int    `json:"priority"`
}

func CreateTask(title, desc string, priority int) Task {

	return Task{
		Title:      title,
		Desc:       desc,
		Priority:   priority,
		Status:     false,
		DateCreate: time.Now(),
		DateCompl:  nil,
	}
}

func TaskConvert(task database.TaskModel) Task {

	return Task{
		Title:      task.Title,
		Desc:       task.Desc,
		Priority:   task.Priority,
		Status:     task.Status,
		DateCreate: task.DateCreate,
		DateCompl:  task.DateCompl,
	}

}

func CreateTaskDTO(title, desc string, priority int) TaskDTO {
	return TaskDTO{
		Title:    title,
		Desc:     desc,
		Priority: priority,
	}
}

func (t *Task) CompleteTask() {
	t.Status = true
	temp := time.Now()
	t.DateCompl = &temp
}

func (t *Task) UncompleteTask() {
	t.Status = false
	t.DateCompl = nil
}
