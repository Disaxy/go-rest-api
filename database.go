package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// DB ..
type DB interface {
	create()
	read()
	update()
	delete()
}

// Category ...
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsLimit     bool   `json:"is_limit"`
}

func connection(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (cat Category) create(db *sql.DB, name, description string, isLimit bool) (int64, error) {
	result, err := db.Exec("insert into `category` (name, description, is_limit) values ($1, $2, $3)", name, description, isLimit)
	if err != nil {
		log.Fatal(err)
	}

	return result.LastInsertId()
}

func (cat Category) read(db *sql.DB) []Category {
	rows, err := db.Query("select * from `category`")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	categories := []Category{}

	for rows.Next() {
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.IsLimit)
		if err != nil {
			log.Fatal(err)
			continue
		}
		categories = append(categories, cat)
	}

	return categories
}

func (cat Category) update(db *sql.DB, id string, description string) (int64, error) {
	result, err := db.Exec("update `category` set description = $2 where id = $1", id, description)
	if err != nil {
		log.Fatal(err)
	}

	return result.RowsAffected()
}

func (cat Category) delete(db *sql.DB, id string) (int64, error) {
	result, err := db.Exec("delete from `category` where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	return result.RowsAffected()
}
