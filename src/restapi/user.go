package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/src/db"
	"todo/src/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}
	err := db.Todo.Select(&users,
		"SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("here1, %s", err.Error()), 500)
		return
	}
	err = user.Get()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't get user, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(user)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = user.Update()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't update user, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(user)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = user.Create()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't create user, err %s", err.Error()), 500)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&user)
	err = user.Delete()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't delete user, err %s", err.Error()), 500)
		return
	}
	w.Write([]byte("User Deleted"))
}

func getUserFromRequest(r *http.Request) (models.User, error) {
	email := ""
	id := ""
	emails, emailOk := r.URL.Query()["email"]
	if emailOk && len(emails) > 0 {
		email = emails[0]
	}
	ids, idOk := r.URL.Query()["id"]
	if idOk && len(ids) > 0 {
		id = ids[0]
	}
	user := models.User{
		ID:    id,
		Email: email,
	}
	if !idOk && !emailOk {
		return user, fmt.Errorf("No identifier to get user, expecting ID or email")
	}
	return user, nil
}
