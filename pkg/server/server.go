package server

import (
	"log"
	"net/http"
	"os"

	"final_project/pkg/api"
	"final_project/pkg/db"
)

const defaultPort = "7540"

func Start() error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return err
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Server started on port %s", port)

	return http.ListenAndServe(":"+port, nil)
}
