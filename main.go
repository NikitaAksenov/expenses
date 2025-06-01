package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const TimeFormatDateOnly = "2006-01-02"

type Expense struct {
	ID       int
	Name     string
	Category string
	Value    float64
	Date     time.Time
	Created  time.Time
	Updated  time.Time
}

func (e Expense) String() string {
	return fmt.Sprintf("#%03d | %15s | %15s | %7.2f | %s | %s | %s", e.ID, e.Name, e.Category, e.Value, e.Date.Format(TimeFormatDateOnly), e.Created.Format(TimeFormatDateOnly), e.Updated.Format(TimeFormatDateOnly))
}

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
		fs := flag.NewFlagSet("add", flag.ExitOnError)

		var multipleFlag bool
		multipleFlagDesc := "If set user can enter expenses one after another until empty row is passed"
		fs.BoolVar(&multipleFlag, "multiple", false, multipleFlagDesc)
		fs.BoolVar(&multipleFlag, "m", false, multipleFlagDesc)

		var dateFlag bool
		dateFlagDesc := "If set user will be prompt to enter date (in case expense is not of today)"
		fs.BoolVar(&dateFlag, "date", false, dateFlagDesc)
		fs.BoolVar(&dateFlag, "d", false, dateFlagDesc)

		var dateFilledFlag string
		dateFilledFlagDesc := "Sets expense date to this flag's value"
		fs.StringVar(&dateFilledFlag, "datefilled", "", dateFilledFlagDesc)
		fs.StringVar(&dateFilledFlag, "df", "", dateFilledFlagDesc)

		fs.Parse(os.Args[2:])

		if dateFlag && dateFilledFlag != "" {
			fmt.Println("Can't have `date` and `datefilled` flags both present")
			return
		}

		if multipleFlag {
			fmt.Println("Adding expense entries")
		} else {
			fmt.Println("Adding expense entry")
		}

		for ok := true; ok; ok = multipleFlag {
			fmt.Print("Name:     ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSuffix(name, "\r\n")

			if name == "" {
				break
			}

			fmt.Print("Category: ")
			category, _ := reader.ReadString('\n')
			category = strings.TrimSuffix(category, "\r\n")

			fmt.Print("Value:    ")
			var value float64
			_, err := fmt.Scanf("%f\r\n", &value)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			date := time.Now().UTC()
			if dateFlag || dateFilledFlag != "" {
				dateStr := dateFilledFlag

				if dateFlag {
					fmt.Print("Date:     ")
					dateStr, _ = reader.ReadString('\n')
					dateStr = strings.TrimSuffix(dateStr, "\r\n")
				}

				date, err = time.Parse(TimeFormatDateOnly, dateStr)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}

			dateStr := date.Format(TimeFormatDateOnly)

			_, err = db.Exec("insert into expenses (Name, Category, Value, Date, Created, Updated) values ($1, $2, $3, $4, $5, $5)",
				name, category, value, dateStr, time.Now().UTC().Format(TimeFormatDateOnly))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case "upd":
		fs := flag.NewFlagSet("upd", flag.ExitOnError)

		var idFlag int
		fs.IntVar(&idFlag, "id", -1, "ID of entry to update")

		var nameFlag bool
		nameFlagDesc := "If set user will be prompt to enter Name to update"
		fs.BoolVar(&nameFlag, "name", false, nameFlagDesc)
		fs.BoolVar(&nameFlag, "n", false, nameFlagDesc)

		var categoryFlag bool
		categoryFlagDesc := "If set user will be prompt to enter Category to update"
		fs.BoolVar(&categoryFlag, "category", false, categoryFlagDesc)
		fs.BoolVar(&categoryFlag, "c", false, categoryFlagDesc)

		var valueFlag bool
		valueFlagDesc := "If set user will be prompt to enter Value to update"
		fs.BoolVar(&valueFlag, "value", false, valueFlagDesc)
		fs.BoolVar(&valueFlag, "v", false, valueFlagDesc)

		var dateFlag bool
		dateFlagDesc := "If set user will be prompt to enter Date to update"
		fs.BoolVar(&dateFlag, "date", false, dateFlagDesc)
		fs.BoolVar(&dateFlag, "d", false, dateFlagDesc)

		fs.Parse(os.Args[2:])

		fmt.Println("Updating entry with id", idFlag)

		if !nameFlag && !categoryFlag && !valueFlag && !dateFlag {
			fmt.Println("No columns were passed")
			fs.Usage()
			return
		}

		var args []string

		if nameFlag {
			fmt.Print("Name:     ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSuffix(name, "\r\n")

			args = append(args, fmt.Sprintf("Name = '%s'", name))
		}

		if categoryFlag {
			fmt.Print("Category: ")
			category, _ := reader.ReadString('\n')
			category = strings.TrimSuffix(category, "\r\n")

			args = append(args, fmt.Sprintf("Category = '%s'", category))
		}

		if valueFlag {
			fmt.Print("Value:    ")
			var value float64
			_, err := fmt.Scanf("%f\r\n", &value)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			args = append(args, fmt.Sprintf("Value = %f", value))
		}

		if dateFlag {
			fmt.Print("Date:     ")
			dateStr, _ := reader.ReadString('\n')
			dateStr = strings.TrimSuffix(dateStr, "\r\n")

			_, err := time.Parse(TimeFormatDateOnly, dateStr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			args = append(args, fmt.Sprintf("Date = '%s'", dateStr))
		}

		updatedStr := time.Now().UTC().Format(TimeFormatDateOnly)
		args = append(args, fmt.Sprintf("Updated = '%s'", updatedStr))

		query := fmt.Sprintf(`
		update expenses
			set %s
			where ID = %d
		`, strings.Join(args, ","), idFlag)

		_, err := db.Exec(query)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Query:", query)
			return
		}
	case "list":
		fmt.Println("Listing expense entries")

		fs := flag.NewFlagSet("list", flag.ExitOnError)

		var dateFlag bool
		dateFlagDesc := "List expenses for date"
		fs.BoolVar(&dateFlag, "date", false, dateFlagDesc)
		fs.BoolVar(&dateFlag, "d", false, dateFlagDesc)

		var dateRangeFlag bool
		dateRangeFlagDesc := "List expenses for date range"
		fs.BoolVar(&dateRangeFlag, "daterange", false, dateRangeFlagDesc)
		fs.BoolVar(&dateRangeFlag, "dr", false, dateRangeFlagDesc)

		var categoryFlag bool
		categoryFlagDesc := "List expenses for category"
		fs.BoolVar(&categoryFlag, "category", false, categoryFlagDesc)
		fs.BoolVar(&categoryFlag, "c", false, categoryFlagDesc)

		var categoriesFlag bool
		categoriesFlagDesc := "List expenses for categories"
		fs.BoolVar(&categoriesFlag, "categories", false, categoriesFlagDesc)
		fs.BoolVar(&categoriesFlag, "cs", false, categoriesFlagDesc)

		var limitFlag int
		limitFlagDesc := "Limit expenses count"
		fs.IntVar(&limitFlag, "limit", 100, limitFlagDesc)
		fs.IntVar(&limitFlag, "l", 100, limitFlagDesc)

		fs.Parse(os.Args[2:])

		if dateFlag && dateRangeFlag {
			fmt.Println("Can't have `date` and `daterange` flags both present")
			return
		}

		if categoryFlag && categoriesFlag {
			fmt.Println("Can't have `category` and `categories` flags both present")
			return
		}

		if limitFlag <= 0 {
			fmt.Println("Limit flag must be > 0, but entered", limitFlag)
			return
		}

		var whereArgs []string

		if dateFlag {
			fmt.Print("Date:     ")
			dateStr, _ := reader.ReadString('\n')
			dateStr = strings.TrimSuffix(dateStr, "\r\n")

			_, err := time.Parse(TimeFormatDateOnly, dateStr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			whereArgs = append(whereArgs, fmt.Sprintf("Date = '%s'", dateStr))
		}

		if dateRangeFlag {
			{
				fmt.Print("Date from:     ")
				dateStr, _ := reader.ReadString('\n')
				dateStr = strings.TrimSuffix(dateStr, "\r\n")

				_, err := time.Parse(TimeFormatDateOnly, dateStr)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				whereArgs = append(whereArgs, fmt.Sprintf("Date >= '%s'", dateStr))
			}
			{
				fmt.Print("Date to  :     ")
				dateStr, _ := reader.ReadString('\n')
				dateStr = strings.TrimSuffix(dateStr, "\r\n")

				_, err := time.Parse(TimeFormatDateOnly, dateStr)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				whereArgs = append(whereArgs, fmt.Sprintf("Date <= '%s'", dateStr))
			}
		}

		if categoryFlag || categoriesFlag {
			var categories []string

			for ok := true; ok; ok = categoriesFlag {
				fmt.Print("Category: ")
				category, _ := reader.ReadString('\n')
				category = strings.TrimSuffix(category, "\r\n")

				if category == "" {
					break
				}

				categories = append(categories, "'"+category+"'")
			}

			if len(categories) > 0 {
				whereArgs = append(whereArgs, fmt.Sprintf("Category IN (%s)", strings.Join(categories, ",")))
			}
		}

		// Make query
		query := "SELECT * FROM expenses"

		if len(whereArgs) > 0 {
			whereStmt := strings.Join(whereArgs, " AND ")
			query += "\nWHERE " + whereStmt
		}

		query += "\nORDER BY Date"
		query += fmt.Sprintf("\nLIMIT %d", limitFlag)

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()

		expenseEntries := []Expense{}

		for rows.Next() {
			expenseEntry := Expense{}
			var dateStr, createdStr, updatedStr string
			err := rows.Scan(&expenseEntry.ID, &expenseEntry.Name, &expenseEntry.Category, &expenseEntry.Value, &dateStr, &createdStr, &updatedStr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			expenseEntry.Date, err = time.Parse(TimeFormatDateOnly, dateStr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			expenseEntry.Created, err = time.Parse(TimeFormatDateOnly, createdStr)
			if err != nil {
				fmt.Println(err)
				continue
			}
			expenseEntry.Updated, err = time.Parse(TimeFormatDateOnly, updatedStr)
			if err != nil {
				fmt.Println(err)
				continue
			}

			expenseEntries = append(expenseEntries, expenseEntry)
		}

		fmt.Printf(" %3s | %15s | %15s | %7s | %10s | %10s | %10s\n", "ID", "Name", "Category", "Value", "Date", "Created", "Updated")
		for _, expenseEntry := range expenseEntries {
			fmt.Println(expenseEntry)
		}
	case "del":
		fs := flag.NewFlagSet("del", flag.ExitOnError)

		var idFlag int
		fs.IntVar(&idFlag, "id", -1, "ID of entry to delete")

		fs.Parse(os.Args[2:])

		fmt.Println("Deleting entry with id", idFlag)

		_, err := db.Exec("delete from expenses where ID = $1", idFlag)
		if err != nil {
			panic(err)
		}
	}
}
