package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dm-vev/gowrite"
	"github.com/joho/godotenv"
)

const (
	dbID  = "example_database_id"
	colID = "example_collection_id"
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

	docs, err := databases.ListDocuments(dbID, colID, nil)
	if err != nil {
		log.Fatalf("failed to list documents: %v", err)
	}

	for i, d := range docs {
		fmt.Printf("%d) %+v\n", i+1, d)
	}
}
