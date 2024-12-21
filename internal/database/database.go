package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func DB() *sql.DB{

	host:= os.Getenv("DB_HOST")
	port:= os.Getenv("DB_PORT")
	user:= os.Getenv("DB_USER")
	password:= os.Getenv("DB_PASSWORD")
	dbname:= os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB ....")

	return db

}