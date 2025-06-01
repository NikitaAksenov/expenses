package commands

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/NikitaAksenov/expenses/internal"
	"github.com/NikitaAksenov/expenses/internal/models"
)

func CommandList(r *bufio.Reader, args []string, db *sql.DB) {
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

	fs.Parse(args)

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
		dateStr, _ := r.ReadString('\n')
		dateStr = strings.TrimSuffix(dateStr, "\r\n")

		_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		whereArgs = append(whereArgs, fmt.Sprintf("Date = '%s'", dateStr))
	}

	if dateRangeFlag {
		{
			fmt.Print("Date from:     ")
			dateStr, _ := r.ReadString('\n')
			dateStr = strings.TrimSuffix(dateStr, "\r\n")

			_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			whereArgs = append(whereArgs, fmt.Sprintf("Date >= '%s'", dateStr))
		}
		{
			fmt.Print("Date to  :     ")
			dateStr, _ := r.ReadString('\n')
			dateStr = strings.TrimSuffix(dateStr, "\r\n")

			_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
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
			category, _ := r.ReadString('\n')
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

	expenseEntries := []models.Expense{}

	for rows.Next() {
		expenseEntry := models.Expense{}
		var dateStr, createdStr, updatedStr string
		err := rows.Scan(&expenseEntry.ID, &expenseEntry.Name, &expenseEntry.Category, &expenseEntry.Value, &dateStr, &createdStr, &updatedStr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		expenseEntry.Date, err = time.Parse(internal.TimeFormatDateOnly, dateStr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		expenseEntry.Created, err = time.Parse(internal.TimeFormatDateOnly, createdStr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		expenseEntry.Updated, err = time.Parse(internal.TimeFormatDateOnly, updatedStr)
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
}
