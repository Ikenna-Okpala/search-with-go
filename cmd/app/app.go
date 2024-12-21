package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webcrawler/internal/database"
	"webcrawler/web"

	"github.com/joho/godotenv"
)



func main() {

	//scraped 1000 wiki articles in 10 minutes

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err:= godotenv.Load("../../.env")

	if err != nil{
		log.Fatal("Env set up error")
	}

	serverIP:= os.Getenv("SERVER_IP")
	serverPort:= os.Getenv("SERVER_PORT")


	mux:= http.NewServeMux()

	db:= database.DB()

	handler:= web.Handler{DB: db}


	mux.HandleFunc("GET /", handler.Search)

	

	http.ListenAndServe(fmt.Sprintf("%s:%s", serverIP, serverPort), mux)


}