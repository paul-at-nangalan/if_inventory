package main

import (
	"database/sql"
	"fmt"
	"if_inventory/errorhandler"
	"if_inventory/services"
	"net/http"
	"os"
)

func runServer(){
	addr := os.ExpandEnv("PORT")
	cert := os.ExpandEnv("CERTFILE")
	keyfile := os.ExpandEnv("KEYFILE")
	err := http.ListenAndServeTLS(addr, cert, keyfile, nil)
	if err != nil{
		fmt.Println("Failed to start server: ", err)
	}
}

func main(){
	defer errorhandler.HandlePanic()

	var db *sql.DB

	////TODO ... connect to the database

	services.NewSpacesrafts(db)
	runServer()
}
