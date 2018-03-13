package models

import (
	"errors"
	"fmt"
	"todo/src/db"
)

// Project - to group tasks
type Project struct {
	ID     string `db:"id"      json:"id"`
	Name   string `db:"name"    json:"name"`
	UserID string `db:"user_id" json:"user_id"`
}

// Update syncs the struct instance changes into the database
func (p *Project) Get() error {
	// could potentially get other ways as well, keeping this simple
	if p.Name != "" {
		err := db.Todo.Get(p,
			`Select * from projects where name = $1`,
			p.Name)
		if err != nil {
			return fmt.Errorf("No project with given name %s, err = %v", p.Name, err)
		}
	} else if p.ID != "" {
		err := db.Todo.Get(p,
			`Select * from projects where id = $1`,
			p.ID)
		if err != nil {
			return errors.New("No project with specified id to get")
		}
	} else {
		return errors.New("Need name or id for project to get")
	}
	return nil
}

func (p *Project) Create() error {
	err := db.Todo.Get(&p.ID,
		`INSERT INTO projects (name, user_id) VALUES ($1, $2) RETURNING id`,
		p.Name, p.UserID,
	)
	return err
}

// Update syncs the struct instance changes into the database
func (p *Project) Update() error {
	var prevProject Project
	err := db.Todo.Get(&prevProject,
		`SELECT * FROM projects WHERE id = $1 OR name = $2`,
		p.ID, p.Name)
	if err != nil {
		return errors.New("No project with specified ID to update")
	}

	_, err = db.Todo.Exec(
		`UPDATE project SET
		name=$2,
		user_id=$3
	   WHERE id=$1`,
		prevProject.ID,
		p.Name,
		p.UserID)
	if err != nil {
		return err
	}
	return nil
}

// Delete the struct project from the database
func (p *Project) Delete() error {
	_, err := db.Todo.Exec(
		`DELETE FROM projects where id = $1 or name = $2`,
		p.ID, p.Name)
	if err != nil {
		return errors.New("No Project to delete")
	}
	return nil
}
