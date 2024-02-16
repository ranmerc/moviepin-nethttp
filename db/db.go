// This package provides a global db connection variable.
package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("postgres", "user=ranmerc dbname=moviepin sslmode=disable")

	if err != nil {
		fmt.Println("failed to open connection with the Database")
		fmt.Println(err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("failed to connect to Database")
		fmt.Println(err)
		os.Exit(1)
	}
}
