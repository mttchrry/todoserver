package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/src/models"
)

func GetProject(w http.ResponseWriter, r *http.Request) {
	project, err := getProjectFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = project.Get()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get project, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(project)
}

var CreateProject = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	proj := models.Project{}
	err := json.NewDecoder(r.Body).Decode(&proj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, err := GetUserFromToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get user from token, err : %s", err.Error()), 500)
	}
	proj.UserID = user.ID
	err = proj.Create()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't create Project, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(proj)
})

var UpdateProject = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	proj := models.Project{}
	err := json.NewDecoder(r.Body).Decode(&proj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	user, err := GetUserFromToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get user from token, err : %s", err.Error()), 500)
	}
	proj.UserID = user.ID
	err = proj.Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't update Project, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(proj)
})

var DeleteProject = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func getProjectFromRequest(r *http.Request) (models.Project, error) {
	name := ""
	id := ""
	names, nameOk := r.URL.Query()["name"]
	if nameOk && len(names) > 0 {
		name = names[0]
	}
	ids, idOk := r.URL.Query()["id"]
	if idOk && len(ids) > 0 {
		id = ids[0]
	}
	proj := models.Project{
		ID:   id,
		Name: name,
	}
	if !idOk && !nameOk {
		return proj, fmt.Errorf("No identifier to get project, expecting ID or name")
	}
	return proj, nil
}
