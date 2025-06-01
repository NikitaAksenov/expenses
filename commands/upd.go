package commands

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/NikitaAksenov/expenses/internal"
)

func CommandUpdate(r *bufio.Reader, args []string, db *sql.DB) {
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

	fs.Parse(args)

	fmt.Println("Updating entry with id", idFlag)

	if !nameFlag && !categoryFlag && !valueFlag && !dateFlag {
		fmt.Println("No columns were passed")
		fs.Usage()
		return
	}

	var setArgs []string

	if nameFlag {
		fmt.Print("Name:     ")
		name, _ := r.ReadString('\n')
		name = strings.TrimSuffix(name, "\r\n")

		setArgs = append(setArgs, fmt.Sprintf("Name = '%s'", name))
	}

	if categoryFlag {
		fmt.Print("Category: ")
		category, _ := r.ReadString('\n')
		category = strings.TrimSuffix(category, "\r\n")

		setArgs = append(setArgs, fmt.Sprintf("Category = '%s'", category))
	}

	if valueFlag {
		fmt.Print("Value:    ")
		var value float64
		_, err := fmt.Scanf("%f\r\n", &value)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		setArgs = append(setArgs, fmt.Sprintf("Value = %f", value))
	}

	if dateFlag {
		fmt.Print("Date:     ")
		dateStr, _ := r.ReadString('\n')
		dateStr = strings.TrimSuffix(dateStr, "\r\n")

		_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		setArgs = append(setArgs, fmt.Sprintf("Date = '%s'", dateStr))
	}

	updatedStr := time.Now().UTC().Format(internal.TimeFormatDateOnly)
	setArgs = append(setArgs, fmt.Sprintf("Updated = '%s'", updatedStr))

	query := fmt.Sprintf(`
		update expenses
			set %s
			where ID = %d
		`, strings.Join(setArgs, ","), idFlag)

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Query:", query)
		return
	}
}
