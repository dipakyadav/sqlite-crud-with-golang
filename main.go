package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// todo struct, represent single task
type todo struct {
	id      int
	task    string
	owner   string
	checked int
}

func main() {

	// Remove the todo database file if exists.
	// Comment out the below line if you don't want to remove the database.
	os.Remove("./todo.db")

	// Open database connection
	db, err := sql.Open("sqlite3", "./todo.db")

	// Check if database connection was opened successfully
	if err != nil {
		// Print error and exit if there was problem opening connection.
		log.Fatal(err)
	}
	// close database connection before exiting program.
	defer db.Close()

	{ // Create table Block
		// SQL statement to create a task table, with no records in it.
		sqlStmt := `
		CREATE TABLE task (id INTEGER NOT NULL PRIMARY KEY, task TEXT, owner TEXT, checked INTEGER);
		DELETE FROM task;
		`
		// Execute the SQL statement
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}
	}

	{ // Create records Block
		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		// Prepare prepared statement that can be reused.
		stmt, err := tx.Prepare("INSERT INTO task(id, task, owner, checked) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		// close statement before exiting program.
		defer stmt.Close()

		// Create empty slice of todo struct pointers.
		tasks := []*todo{}
		// Create task, and append to tasks.
		tasks = append(tasks, &todo{id: 1, task: "Learn Go", owner: "Dipak", checked: 0})
		tasks = append(tasks, &todo{id: 2, task: "Learn Sqlite", owner: "Dipak", checked: 0})
		tasks = append(tasks, &todo{id: 3, task: "Learn sql driver specification", owner: "Dipak", checked: 0})
		tasks = append(tasks, &todo{id: 4, task: "Write simple sqlite driver for go", owner: "Dipak", checked: 0})

		for i := range tasks {
			// Insert records
			// Execute statements for each tasks
			_, err = stmt.Exec(tasks[i].id, tasks[i].task, tasks[i].owner, tasks[i].checked)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Commit the transaction, so that inserts are permanent.
		tx.Commit()
	}

	{ // Read records Block
		// Start reading records
		// Hint: try changing query to meet your needs.
		stmt, err := db.Prepare("SELECT id, task, owner from task where checked = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(0)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var task string
			var owner string
			err = rows.Scan(&id, &task, &owner)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, task, owner)
		}
		// To just check if any error was occured during iteration.
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	{ // Update records Block
		// Updating record(s)
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("UPDATE task SET owner = ? where id = 4")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec("GoLang Community")
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}

	{ // Delete records block
		// Delete record(s)s
		// _, err = db.Exec("DELETE from task")
		// if err != nil {
		//	log.Fatal(err)
		// }
	}

}
