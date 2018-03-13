package main

import (
	"log"
	"net/http"
	"todo/src/restapi"

	"github.com/gorilla/mux"
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", restapi.GetUsers).Methods("GET")
	router.HandleFunc("/signin", restapi.SignIn).Methods("POST")

	router.HandleFunc("/user", restapi.GetUser).Methods("GET")
	router.HandleFunc("/user", restapi.CreateUser).Methods("POST")
	router.HandleFunc("/user", restapi.UpdateUser).Methods("PUT")
	router.HandleFunc("/user", restapi.DeleteUser).Methods("DELETE")

	router.HandleFunc("/project", restapi.GetProject).Methods("GET")
	router.Handle("/project", restapi.JwtMiddleware.Handler(restapi.CreateProject)).Methods("POST")
	router.Handle("/project", restapi.JwtMiddleware.Handler(restapi.UpdateProject)).Methods("PUT")
	router.Handle("/project", restapi.JwtMiddleware.Handler(restapi.DeleteProject)).Methods("DELETE")

	router.HandleFunc("/task", restapi.GetTask).Methods("GET")
	router.Handle("/task", restapi.JwtMiddleware.Handler(restapi.CreateTask)).Methods("POST")
	router.Handle("/task", restapi.JwtMiddleware.Handler(restapi.UpdateTask)).Methods("PUT")
	router.Handle("/task", restapi.JwtMiddleware.Handler(restapi.DeleteTask)).Methods("DELETE")

	router.HandleFunc("/tasks", restapi.GetTasks).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
