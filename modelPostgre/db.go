package modelPostgre

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"log"
)

//
//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "postgres"
//	password = "qwe"
//	dbname   = "test"
//)

var PostgreDB *sql.DB
var err error

func SetDBParams(host string, port int, user, password, dbname string) error {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	PostgreDB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("PostgreSQL SetDBParams sql.Open error: ", err)
		return err
	}
	//defer PostgreDB.Close()

	err = PostgreDB.Ping()
	if err != nil {
		log.Println("PostgreSQL SetDBParams db.Ping() error: ", err)
		return err
	}
	fmt.Println("Successfully connected!")
	return nil
}
