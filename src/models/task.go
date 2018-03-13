package models

import (
	"errors"
	"fmt"
	"time"
	"todo/src/db"
)

// Task - the thing to do.
type Task struct {
	ID          string     `db:"id"          json:"id"`
	Summary     string     `db:"summary"     json:"summary"`
	ProjectID   string     `db:"project_id"  json:"project_id"`
	UserID      string     `db:"user_id"     json:"user_id"`
	Description *string    `db:"description" json:"description"`
	DueDate     *time.Time `db:"due_date"    json:"due_date"`
	Priority    *int       `db:"priority"    json:"priority"`
}

// Update syncs the struct instance changes into the database
func (t *Task) Get() error {
	// could potentially get other ways as well, keeping this simple
	if t.ID != "" {
		err := db.Todo.Get(t,
			`Select * from tasks where id = $1`,
			t.ID)
		if err != nil {
			return errors.New("No task with given id")
		}
	} else {
		return errors.New("Need id for task to get")
	}
	return nil
}

func (t *Task) Create() error {
	err := db.Todo.Get(&t.ID,
		`INSERT INTO tasks (
			summary,
			project_id,
			user_id,
			description,
			due_date,
			priority
		) VALUES ($1, $2,  $3, $4, $5, $6) RETURNING id`,
		t.Summary, t.ProjectID, t.UserID, t.Description,
		t.DueDate, t.Priority)
	return err
}

// Update syncs the struct instance changes into the database
func (t *Task) Update() error {
	_, err := db.Todo.Exec(
		`UPDATE tasks SET
		summary=$2,
		project_id=$3,
		user_id=$4,
		description=$5,
		due_date=$6,
		priority=$7
	   WHERE id=$1`,
		t.ID, t.Summary, t.ProjectID, t.UserID, t.Description,
		t.DueDate, t.Priority)
	if err != nil {
		return err
	}
	return nil
}

// Delete the struct user from the database
func (t *Task) Delete() error {
	_, err := db.Todo.Exec(
		`DELETE FROM tasks where id = $1`,
		t.ID)
	if err != nil {
		return errors.New("No Project to delete")
	}
	return nil
}

func GetTasksByFilters(filters map[string]string) ([]Task, error) {
	query := `Select * from tasks`

	isFirst := true

	for filterType, value := range filters {
		if value == "" {
			continue
		}
		queryModifier := ""
		if isFirst {
			queryModifier = "WHERE"
		} else {
			queryModifier = "AND"
		}
		switch filterType {
		case "user_id":
			queryModifier = fmt.Sprintf("%s user_id = '%s'", queryModifier, value)
		case "project_id":
			queryModifier = fmt.Sprintf("%s project_id = '%s'", queryModifier, value)
		case "priority":
			fmt.Println(value)
			queryModifier = fmt.Sprintf("%s priority = %+v", queryModifier, value)
		case "due_date":
			queryModifier = fmt.Sprintf("%s due_date = '%s'", queryModifier, value)
		default:
			// log, can ignore
			fmt.Printf("Invalid key value for filters, %s: %s\n", filterType, value)
			continue
		}
		isFirst = false
		query = fmt.Sprintf("%s %s", query, queryModifier)
	}
	query += ";"
	tasks := []Task{}
	fmt.Println(query)
	err := db.Todo.Select(&tasks, query)
	return tasks, err
}
