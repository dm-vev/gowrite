package main

import (
	"log"
	"os"

	"github.com/dm-vev/gowrite"
	"github.com/joho/godotenv"
)

const (
	dbID  = "example_database_id"
	colID = "example_collection_id"
	docID = "example_document_id"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file:", err)
	}

	endpoint := os.Getenv("APPWRITE_INSTANCE")
	project := os.Getenv("APPWRITE_PROJECT")
	token := os.Getenv("APPWRITE_TOKEN")

	if endpoint == "" || project == "" || token == "" {
		log.Fatal("missing required environment variables: APPWRITE_INSTANCE, APPWRITE_PROJECT, APPWRITE_TOKEN")
	}

	client := gowrite.NewClient(endpoint, project, token)
	databases := gowrite.NewDatabases(client)

	if err := databases.DeleteDocument(dbID, colID, docID); err != nil {
		log.Fatalf("failed to delete document: %v", err)
	}

	log.Println("document deleted")
}
