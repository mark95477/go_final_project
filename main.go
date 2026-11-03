package main

import (
	"TODO_List/pkg/db"
	"TODO_List/pkg/server"
	"log"
	"os"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.DB.Close()

	logger := log.New(os.Stdout, "INFO:", log.LstdFlags)
	serv := server.Server(logger)
	err = serv.ListenAndServe()
	logger.Fatal(err)
}
