package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NikitaAksenov/expenses/commands"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	appName string
	dbName  string

	db *sql.DB
}

func (app *application) init() error {
	fmt.Println("Init app", app.appName)

	dbPath, err := app.getDatabasePath()
	if err != nil {
		return err
	}

	// Check if db exists
	_, err = os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	if !dbExists {
		// If db didn't existed, then we need to create "expenses" table

		fmt.Println("Database not found, creating new one...")

		schemaFile := "db/create_table_expenses.sql"

		createTableSQLBytes, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		_, err = db.Exec(string(createTableSQLBytes))
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	app.db = db

	fmt.Println("Successfully inited app")

	return nil
}

func (app *application) close() {
	fmt.Println("Closing app", app.appName)
	app.db.Close()
}

func (app *application) getDatabasePath() (string, error) {
	appDataDir, err := os.UserConfigDir() // Typically resolves to AppData\Roaming on Windows
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(appDataDir, app.appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	dbPath := filepath.Join(appDir, app.dbName)
	return dbPath, nil
}

func main() {
	app := &application{
		appName: "expenses",
		dbName:  "expenses.db",
	}

	err := app.init()
	if err != nil {
		fmt.Println("failed to init application: %w", err)
		return
	}
	defer app.close()

	if len(os.Args) < 2 {
		return
	}

	reader := bufio.NewReader(os.Stdin)

	args := os.Args[2:]

	switch os.Args[1] {
	case "add":
		commands.CommandAdd(reader, args, app.db)
	case "upd":
		commands.CommandUpdate(reader, args, app.db)
	case "list":
		commands.CommandList(reader, args, app.db)
	case "del":
		commands.CommandDelete(reader, args, app.db)
	case "rep":
		commands.CommandReport(reader, args, app.db)
	}
}
