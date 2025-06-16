package commands

import (
	"bufio"
	"database/sql"
)

func CommandReport(r *bufio.Reader, args []string, db *sql.DB) {
	// var dateFromStr, dateToStr string
	// {
	// 	fmt.Print("Date from:     ")
	// 	dateStr, _ := r.ReadString('\n')
	// 	dateStr = strings.TrimSuffix(dateStr, "\r\n")

	// 	_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// 	}

	// 	dateFromStr = dateStr
	// }
	// {
	// 	fmt.Print("Date to  :     ")
	// 	dateStr, _ := r.ReadString('\n')
	// 	dateStr = strings.TrimSuffix(dateStr, "\r\n")

	// 	_, err := time.Parse(internal.TimeFormatDateOnly, dateStr)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return
	// 	}

	// 	dateToStr = dateStr
	// }

	// query := "SELECT * FROM expenses"
	// query += fmt.Sprintf("\nWHERE Date >= '%s' AND Date <= '%s'", dateFromStr, dateToStr)

	// fmt.Println(query)

	// rows, err := db.Query(query)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer rows.Close()

	// expenseEntries := []models.Expense{}

	// for rows.Next() {
	// 	expenseEntry := models.Expense{}
	// 	var dateStr, createdStr, updatedStr string
	// 	err := rows.Scan(&expenseEntry.ID, &expenseEntry.Name, &expenseEntry.Category, &expenseEntry.Value, &dateStr, &createdStr, &updatedStr)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	expenseEntry.Date, err = time.Parse(internal.TimeFormatDateOnly, dateStr)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	expenseEntry.Created, err = time.Parse(internal.TimeFormatDateOnly, createdStr)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	expenseEntry.Updated, err = time.Parse(internal.TimeFormatDateOnly, updatedStr)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	expenseEntries = append(expenseEntries, expenseEntry)
	// }

	// // fmt.Printf(" %3s | %15s | %15s | %7s | %10s | %10s | %10s\n", "ID", "Name", "Category", "Value", "Date", "Created", "Updated")
	// // for _, expenseEntry := range expenseEntries {
	// // 	fmt.Println(expenseEntry)
	// // }

	// if len(expenseEntries) <= 0 {
	// 	fmt.Println("No expense entries")
	// 	return
	// }

	// for
}
