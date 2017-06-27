package modelPostgre

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwe"
	dbname   = "test"
)

var db *sql.DB
var err error

func SetDBParams(host, port, user, password, dbname string) error {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Successfully connected!")
	return nil
}
