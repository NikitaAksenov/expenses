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

func CommandAdd(r *bufio.Reader, args []string, db *sql.DB) {
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

	fs.Parse(args)

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
		name, _ := r.ReadString('\n')
		name = strings.TrimSuffix(name, "\r\n")

		if name == "" {
			break
		}

		fmt.Print("Category: ")
		category, _ := r.ReadString('\n')
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
				dateStr, _ = r.ReadString('\n')
				dateStr = strings.TrimSuffix(dateStr, "\r\n")
			}

			date, err = time.Parse(internal.TimeFormatDateOnly, dateStr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		dateStr := date.Format(internal.TimeFormatDateOnly)

		_, err = db.Exec("insert into expenses (Name, Category, Value, Date, Created, Updated) values ($1, $2, $3, $4, $5, $5)",
			name, category, value, dateStr, time.Now().UTC().Format(internal.TimeFormatDateOnly))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
