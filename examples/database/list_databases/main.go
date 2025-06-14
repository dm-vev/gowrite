package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dm-vev/gowrite"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("warning: failed to load .env file:", err)
	}

	endpoint := os.Getenv("APPWRITE_INSTANCE")
	project := os.Getenv("APPWRITE_PROJECT")
	token := os.Getenv("APPWRITE_TOKEN")

	if endpoint == "" || project == "" || token == "" {
		log.Fatal("missing required environment variables: APPWRITE_INSTANCE, APPWRITE_PROJECT, APPWRITE_TOKEN")
	}

	client := gowrite.NewClient(endpoint, project, token)
	databases := gowrite.NewDatabases(client)

	dbs, err := databases.ListDatabases()
	if err != nil {
		log.Fatalf("failed to list databases: %v", err)
	}

	for i, db := range dbs {
		fmt.Printf("%d) %+v\n", i+1, db)
	}
}
