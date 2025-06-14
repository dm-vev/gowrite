package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dm-vev/gowrite"
	"github.com/joho/godotenv"
)

const (
	dbID   = "example_database_id"
	dbName = "Example Database"
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

	db, err := databases.CreateDatabase(dbID, dbName, true)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	fmt.Printf("created: %+v\n", db)
}
