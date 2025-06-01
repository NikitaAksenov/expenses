package commands

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
)

func CommandDelete(r *bufio.Reader, args []string, db *sql.DB) {
	fs := flag.NewFlagSet("del", flag.ExitOnError)

	var idFlag int
	fs.IntVar(&idFlag, "id", -1, "ID of entry to delete")

	fs.Parse(args)

	fmt.Println("Deleting entry with id", idFlag)

	_, err := db.Exec("delete from expenses where ID = $1", idFlag)
	if err != nil {
		panic(err)
	}
}
