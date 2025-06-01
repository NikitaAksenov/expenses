package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/NikitaAksenov/expenses/commands"

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

	reader := bufio.NewReader(os.Stdin)

	switch os.Args[1] {
	case "add":
		commands.CommandAdd(reader, os.Args[2:], db)
	case "upd":
		commands.CommandUpdate(reader, os.Args[2:], db)
	case "list":
		commands.CommandList(reader, os.Args[2:], db)
	case "del":
		commands.CommandDelete(reader, os.Args[2:], db)
	}
}
