package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("expenses")

	if len(os.Args) < 2 {
		return
	}

	dbPath := "./db/expenses.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "add":
	case "upd":
	case "list":
	case "del":
	}
}
