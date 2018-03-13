package models

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetTaskByFilters(t *testing.T) {

	// Ideally we move this boilerplate into factories to generate test data
	// and handle cleanup, as well as mock out the DB itself
	pswd := "test123"
	user1 := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "tester@testing.com",
		Password:  &pswd,
	}
	err := user1.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer user1.Delete()

	pswd = "123test"
	user2 := User{
		FirstName: "Test2",
		LastName:  "User2",
		Email:     "tester2@testing.com",
		Password:  &pswd,
	}
	err = user2.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer user2.Delete()

	proj1 := Project{
		Name:   "Project1",
		UserID: user1.ID,
	}
	err = proj1.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer proj1.Delete()

	proj2 := Project{
		Name:   "Project2",
		UserID: user2.ID,
	}
	err = proj2.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer proj2.Delete()

	priority1 := 1
	priority2 := 2
	dueDate1, _ := time.Parse("2006-01-02", "2019-01-01")
	dueDate2, _ := time.Parse("2006-01-02", "2019-01-02")

	task1 := Task{
		Summary:   "Task 1",
		ProjectID: proj1.ID,
		DueDate:   &dueDate1,
		Priority:  &priority1,
		UserID:    user1.ID,
	}
	err = task1.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer task1.Delete()

	task2 := Task{
		Summary:   "Task 2",
		ProjectID: proj1.ID,
		DueDate:   &dueDate2,
		Priority:  &priority2,
		UserID:    user2.ID,
	}
	err = task2.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer task2.Delete()

	task3 := Task{
		Summary:   "Task 3",
		ProjectID: proj2.ID,
		DueDate:   &dueDate1,
		Priority:  &priority2,
		UserID:    user1.ID,
	}
	err = task3.Create()
	if err != nil {
		t.Fatal(err)
	}
	//defer task3.Delete()

	filters := make(map[string]string)
	filters["user_id"] = user1.ID
	tasks, err := GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected two tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.UserID != user1.ID {
			t.Fatal("Not correctly filtered by user")
		}
	}

	filters = make(map[string]string)
	filters["project_id"] = proj2.ID
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.Summary != task3.Summary {
			t.Fatal("Not correctly filtered by projectID")
		}
	}

	filters = make(map[string]string)
	fmt.Println("priority = ", priority2, string(priority2))
	filters["priority"] = strconv.Itoa(priority2)
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.Summary != task3.Summary && task.Summary != task2.Summary {
			t.Fatal("Not correctly filtered by priority")
		}
	}

	filters = make(map[string]string)
	// get just the date part, 2019-01-02
	filters["due_date"] = strings.Split(dueDate2.String(), " ")[0]
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.Summary != task2.Summary {
			t.Fatal("Not correctly filtered by Due Date")
		}
	}

	filters = make(map[string]string)
	// 2 for user1, 2 for project1, only one for both
	filters["project_id"] = proj1.ID
	filters["user_id"] = user1.ID
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.Summary != task1.Summary {
			t.Fatal("Not correctly filtered by user + project")
		}
	}

	filters = make(map[string]string)
	// 2 for duedate 1, 2 for priority2, only one for both
	filters["due_date"] = strings.Split(dueDate1.String(), " ")[0]
	filters["priority"] = strconv.Itoa(priority2)
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 tasks, got %+v", len(tasks))
	}
	for _, task := range tasks {
		if task.Summary != task3.Summary {
			t.Fatal("Not correctly filtered by Due Date + priority")
		}
	}

	filters = make(map[string]string)
	// 2 for duedate 1, 2 for priority2, 2 for Project2, None for all
	filters["due_date"] = strings.Split(dueDate1.String(), " ")[0]
	filters["priority"] = strconv.Itoa(priority2)
	filters["project_id"] = proj1.ID
	tasks, err = GetTasksByFilters(filters)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %+v", len(tasks))
	}
}
