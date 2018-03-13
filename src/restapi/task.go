package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/src/models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]string)
	params := r.URL.Query()

	for key, val := range params {
		if len(val) > 0 {
			// only accepting one val per query type at this time
			filters[key] = val[0]
		}
	}
	tasks, err := models.GetTasksByFilters(filters)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := ""
	ids, idOk := r.URL.Query()["id"]
	if idOk && len(ids) > 0 {
		id = ids[0]
	}
	task := models.Task{
		ID: id,
	}
	err := task.Get()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get task, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(task)
}

var CreateTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, err := GetUserFromToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get user from token, err : %s", err.Error()), 500)
	}
	task.UserID = user.ID
	err = task.Create()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't create task, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(task)
})

var UpdateTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, err := GetUserFromToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get user from token, err : %s", err.Error()), 500)
	}
	task.UserID = user.ID
	err = task.Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't update task, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(task)
})

var DeleteTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	proj := models.Project{}
	err := json.NewDecoder(r.Body).Decode(&proj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = proj.Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't Delete Project, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(proj)
})
