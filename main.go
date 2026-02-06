package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func CreateTask(Title, Desc string, Priority int) Task {
	return Task{
		Title:      Title,
		Desc:       Desc,
		Priority:   Priority,
		Status:     false,
		DateCreate: time.Now(),
		DateCompl:  nil,
	}
}

var list = []Task{}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if len(list) == 0 {
			w.WriteHeader(400)
			return
		}
		URL := r.URL.Path
		if params := strings.SplitN(URL, "/", 3); len(params) == 3 && params[2] != "" {
			for i := 0; i < len(list); i++ {
				if list[i].Title == params[2] {
					data, err := json.Marshal(list[i])
					if err != nil {
						w.WriteHeader(400)
						return
					}
					w.WriteHeader(200)
					w.Write(data)
					return
				} else if i == len(list)-1 {
					w.WriteHeader(400)
					return
				}
			}
		} else {
			for i := 0; i < len(list); i++ {
				data, err := json.Marshal(list[i])
				if err != nil {
					w.WriteHeader(400)
					return
				}
				w.WriteHeader(200)
				w.Write(data)
			}
		}
	} else if r.Method == "POST" {
		var temp Task
		if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
			w.WriteHeader(400)
			return
		}
		for i := 0; i < len(list); i++ {
			if list[i].Title == temp.Title {
				w.WriteHeader(409)
				return
			}
		}
		if temp.Priority > 3 || temp.Priority < 1 {
			w.WriteHeader(400)
			return
		}
		task := CreateTask(temp.Title, temp.Desc, temp.Priority)
		list = append(list, task)
		data, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	} else if r.Method == "PATCH" {
		if len(list) == 0 {
			w.WriteHeader(400)
			return
		}
		URL := r.URL.Path
		httpRequestBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		temp := Task{}
		err = json.Unmarshal(httpRequestBody, &temp)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if params := strings.SplitN(URL, "/", 3); len(params) == 3 && params[2] != "" {
			pos := -1
			for i := 0; i < len(list); i++ {
				if list[i].Title == params[2] {
					pos = i
					break
				} else if i == len(list)-1 {
					w.WriteHeader(400)
					return
				}
			}
			if temp.Title != "" {
				for i := 0; i < len(list); i++ {
					if list[i].Title == temp.Title && i != pos {
						w.WriteHeader(409)
						return
					}
				}
				list[pos].Title = temp.Title
			}
			if temp.Desc != "" {
				list[pos].Desc = temp.Desc
			}
			if temp.Priority != 0 {
				if temp.Priority < 1 || temp.Priority > 3 {
					w.WriteHeader(400)
					return
				}
				list[pos].Priority = temp.Priority
			}
			if temp.Status == true {
				list[pos].Status = true
				x := time.Now()
				list[pos].DateCompl = &x
			}
			data, err := json.Marshal(list[pos])
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
			w.Write(data)
		} else {
			w.WriteHeader(400)
			return
		}
	} else if r.Method == "DELETE" {
		if len(list) == 0 {
			w.WriteHeader(400)
			return
		}
		URL := r.URL.Path
		if params := strings.SplitN(URL, "/", 3); len(params) == 3 && params[2] != "" {
			for i := 1; i <= len(list); i++ {
				if list[i-1].Title == params[2] {
					list = append(list[:i-1], list[i:]...)
					break
				} else if i == len(list) {
					w.WriteHeader(404)
					return
				}
			}
		} else {
			for i := 1; i <= len(list); {
				if list[i-1].Status == true {
					list = append(list[:i-1], list[i:]...)
					continue
				}
				i++
			}
		}
	} else {
		fmt.Println("ХУЕ")
		w.WriteHeader(400)
		return
	}

}

func main() {

	http.HandleFunc("/tasks/", handler)
	if err := http.ListenAndServe(":9111", nil); err != nil {
		fmt.Println(err)
		return
	}

}
