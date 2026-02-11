package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var mtx sync.Mutex
var RWmtx sync.RWMutex

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

var list = map[string]*Task{}

func HandlerCreateTask(w http.ResponseWriter, r *http.Request) {

	mtx.Lock()
	defer mtx.Unlock()
	var temp TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if _, ok := list[temp.Title]; ok {
		w.WriteHeader(409)
		w.Write([]byte("Another task have same title"))
		return
	}
	if temp.Priority > 3 || temp.Priority < 1 {
		w.WriteHeader(400)
		w.Write([]byte("Priority would be for 1 to 3"))
		return
	}
	task := CreateTask(temp.Title, temp.Desc, temp.Priority)
	list[temp.Title] = &task
	data, err := json.MarshalIndent(task, "", "\t")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(data)

}

func HandlerGetTask(w http.ResponseWriter, r *http.Request) {

	RWmtx.Lock()
	defer RWmtx.Unlock()
	title := mux.Vars(r)["title"]
	if _, ok := list[title]; ok {
		data, err := json.MarshalIndent(list[title], "", "\t")
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(data)
		return
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Not found task with this title"))
		return
	}
}

func HandlerGetAllTasks(w http.ResponseWriter, r *http.Request) {

	RWmtx.Lock()
	defer RWmtx.Unlock()
	if len(list) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No have any task"))
		return
	}
	for _, task := range list {
		data, err := json.MarshalIndent(task, "", "\t")
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	}

}

func HandlerChangeTask(w http.ResponseWriter, r *http.Request) {

	mtx.Lock()
	defer mtx.Unlock()
	title := mux.Vars(r)["title"]
	if len(list) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No have any task"))
		return
	}
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	temp := Task{}
	err = json.Unmarshal(httpRequestBody, &temp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	if _, ok := list[title]; !ok {
		w.WriteHeader(404)
		w.Write([]byte("Not found task with this title"))
		return
	}
	if temp.Title != "" {
		if _, ok := list[temp.Title]; ok {
			w.WriteHeader(409)
			w.Write([]byte("Another task have same title"))
			return
		}
		list[temp.Title] = list[title]
		list[temp.Title].Title = temp.Title
		delete(list, title)
		title = temp.Title
	}
	if temp.Desc != "" {
		list[title].Desc = temp.Desc
	}
	if temp.Priority != 0 {
		if temp.Priority < 1 || temp.Priority > 3 {
			w.WriteHeader(400)
			w.Write([]byte("Priority would be for 1 to 3"))
			return
		}
		list[title].Priority = temp.Priority
	}
	if temp.Status != list[title].Status {
		list[title].Status = temp.Status
		if temp.Status == true {
			list[title].CompleteTask()
		} else {
			list[title].UncompleteTask()
		}
	}
	data, err := json.MarshalIndent(list[title], "", "\t")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write(data)

}

func HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {

	title := mux.Vars(r)["title"]
	if len(list) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No have any task"))
		return
	}
	if _, ok := list[title]; !ok {
		w.WriteHeader(404)
		w.Write([]byte("Not found task with this title"))
		return
	}
	delete(list, title)
	w.WriteHeader(200)

}

func HandlerDeleteCompletedTasks(w http.ResponseWriter, r *http.Request) {

	if len(list) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No have any task"))
		return
	}
	for title, task := range list {
		if task.Status == true {
			delete(list, title)
		}
	}

}

func main() {

	router := mux.NewRouter()
	router.Path("/tasks").Methods("POST").HandlerFunc(HandlerCreateTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(HandlerGetAllTasks)
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(HandlerGetTask)
	router.Path("/tasks/{title}").Methods("PATCH").HandlerFunc(HandlerChangeTask)
	router.Path("/tasks/{title}").Methods("DELETE").HandlerFunc(HandlerDeleteTask)
	router.Path("/tasks").Methods("DELETE").Queries("Status", "true").HandlerFunc(HandlerDeleteCompletedTasks)

	if err := http.ListenAndServe(":9111", router); err != nil {
		fmt.Println(err)
		return
	}

}
