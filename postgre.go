package main

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwe"
	dbname   = "test"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	//CREATE TABLE users (
	//	id SERIAL PRIMARY KEY,
	//	age INT,
	//	first_name TEXT,
	//	last_name TEXT,
	//	email TEXT UNIQUE NOT NULL
	//);
	fmt.Println("Successfully connected!")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + ` users(id SERIAL PRIMARY KEY,	age INT,first_name TEXT,last_name TEXT,email TEXT UNIQUE NOT NULL)`)
	if err != nil {
		panic(err)
	}

	sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, 30, RandStringBytes(7)+".oio", "Jonathan", "Calhoun").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	type User struct {
		ID        int
		Age       int
		Email     string
		FirstName string
		LastName  string
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		//var id int
		//var firstName string
		//var email string
		var user User
		err = rows.Scan(&user)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(user)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	//rows, err := db.Query("SELECT id, first_name FROM users LIMIT $1", 3)
	//sqlStatement2 := `SELECT * FROM users WHERE id=$1;`
	//var user User
	//row := db.QueryRow(sqlStatement2, 3)
	//err = row.Scan(&user.ID, &user.Age, &user.FirstName,
	//	&user.LastName, &user.Email)
	//switch err {
	//case sql.ErrNoRows:
	//	fmt.Println("No rows were returned!")
	//	return
	//case nil:
	//	fmt.Println(user)
	//default:
	//	panic(err)
	//}
}
