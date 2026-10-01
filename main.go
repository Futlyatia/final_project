package main

import (
	"log"

	"final_project/pkg/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
