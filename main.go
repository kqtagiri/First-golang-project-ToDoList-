package main

import (
	"ToDo/database"
	"ToDo/todo"

	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var mtx sync.Mutex
var RWmtx sync.RWMutex
var db database.Database

var list = map[string]*todo.Task{}

func HandlerCreateTask(w http.ResponseWriter, r *http.Request) {

	mtx.Lock()
	defer mtx.Unlock()
	var temp todo.TaskDTO
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
	task := todo.CreateTask(temp.Title, temp.Desc, temp.Priority)
	if err := db.Insert(db.Ctx, db.Conn, task.Title, task.Desc, task.Priority, task.Status, task.DateCreate); err != nil {
		fmt.Println(err)
		return
	}
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
	titleStart := title
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
	temp := todo.Task{}
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

	taskmodel := database.TaskModel{}
	taskmodel.Title = list[title].Title
	taskmodel.Desc = list[title].Desc
	taskmodel.Priority = list[title].Priority
	taskmodel.Status = list[title].Status
	taskmodel.DateCompl = list[title].DateCompl
	db.Update(db.Ctx, db.Conn, taskmodel, titleStart)
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
	db.Delete(db.Ctx, db.Conn, title)
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
		db.DeleteCompletedTasks(db.Ctx, db.Conn)
	}

}

func main() {

	ctx := context.Background()
	conn, err := database.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	db.InsertConn(ctx, conn)

	if err := database.CreateTable(db.Ctx, db.Conn); err != nil {
		fmt.Println(err)
		return
	}

	ListTaskModel := []database.TaskModel{}
	err, ListTaskModel = db.GetAllTasks(db.Ctx, db.Conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(ListTaskModel); i++ {
		task := todo.TaskConvert(ListTaskModel[i])
		list[ListTaskModel[i].Title] = &task
	}

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
