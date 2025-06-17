package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NikitaAksenov/expenses/commands"

	_ "modernc.org/sqlite"
)

type application struct {
	appName string
	dbName  string

	db *sql.DB

	infoLogger  *log.Logger
	infoLogFile *os.File
}

func (app *application) init() error {
	appDir, err := app.getApplicationDir()
	if err != nil {
		return err
	}

	err = app.setupLogger(appDir)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	app.infoLogger.Println("Init app", app.appName)

	err = app.setupDatabase(appDir)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	return nil
}

func (app *application) close() {
	app.infoLogger.Println("Closing app", app.appName)

	app.db.Close() // Close db

	app.infoLogFile.Close() // Close infoLog file
}

func (app *application) getApplicationDir() (string, error) {
	appDataDir, err := os.UserConfigDir() // Typically resolves to AppData\Roaming on Windows
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(appDataDir, app.appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return appDir, nil
}

func (app *application) setupDatabase(appDir string) error {
	dbPath := filepath.Join(appDir, app.dbName)

	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	if !dbExists {
		// If db didn't existed, then we need to create "expenses" table

		app.infoLogger.Println("Database not found, creating new one...")

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
	return nil
}

func (app *application) setupLogger(appDir string) error {
	// Check log directory
	logDir := filepath.Join(appDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Create log file
	timestamp := time.Now().Format("2006-01-02_15-04-05") // Format to YYYY-MM-DD_HH-MM-SS
	logFileName := fmt.Sprintf("log_%s_%s.txt", app.appName, timestamp)
	logFilePath := filepath.Join(logDir, logFileName)
	var err error
	app.infoLogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Init logger to file
	multiWriter := io.MultiWriter(app.infoLogFile)
	app.infoLogger = log.New(multiWriter, "", log.LstdFlags)

	return nil
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

	app.infoLogger.Println("Args:", os.Args[1:])
	if len(os.Args) < 2 {
		return
	}

	args := os.Args[2:]

	reader := bufio.NewReader(os.Stdin)

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
