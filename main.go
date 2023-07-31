package main

import (
	"ethbaas/cmd"
	"ethbaas/internal/db"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dbClient, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	if err := dbClient.Setup(); err != nil {
		log.Fatal(err)
	}

	c := cmd.NewCommand(dbClient)
	c.Execute()
}
