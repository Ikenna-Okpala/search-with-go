package main

import (
	"net/http"
	"webcrawler/api/handler"
	"webcrawler/internal/database"
)



func main() {

	db:= database.DB()

	handler:= handler.Handler{DB: db}


	http.HandleFunc("GET /", handler.Search)

	http.ListenAndServe("localhost:8080", nil)


}