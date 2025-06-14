package gowrite_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dm-vev/gowrite"
	"github.com/dm-vev/gowrite/id"
)

func getClient(t *testing.T) *gowrite.AppwriteClient {
	endpoint := os.Getenv("APPWRITE_ENDPOINT")
	projectID := os.Getenv("APPWRITE_PROJECT_ID")
	apiKey := os.Getenv("APPWRITE_API_KEY")

	if endpoint == "" || projectID == "" || apiKey == "" {
		t.Skip("Appwrite credentials are not set")
	}

	return gowrite.NewClient(endpoint, projectID, apiKey)
}

func TestDatabaseIntegration(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	client := getClient(t)
	db := gowrite.NewDatabases(client)

	dbCount := rand.Intn(3) + 1
	type colInfo struct {
		id       string
		name     string
		docCount int
	}
	databases := make([]struct {
		id          string
		name        string
		newName     string
		collections []colInfo
	}, dbCount)

	for i := 0; i < dbCount; i++ {
		dbID := id.Unique()
		dbName := fmt.Sprintf("db_%d", i)
		created, err := db.CreateDatabase(dbID, dbName, true)
		if err != nil {
			t.Fatalf("CreateDatabase: %v", err)
		}
		log.Printf("Created database - ID: %s Name: %s", created.ID, created.Name)
		t.Cleanup(func(id string) func() {
			return func() {
				if err := db.DeleteDatabase(id); err != nil {
					log.Printf("Cleanup DeleteDatabase error: %v", err)
				} else {
					log.Printf("Deleted database - ID: %s", id)
				}
			}
		}(dbID))

		got, err := db.GetDatabase(dbID)
		if err != nil || got.ID != dbID || got.Name != dbName {
			t.Fatalf("GetDatabase mismatch: %v", err)
		}

		newName := fmt.Sprintf("db_%d_updated", i)
		upd, err := db.UpdateDatabase(dbID, newName, true)
		if err != nil || upd.Name != newName {
			t.Fatalf("UpdateDatabase mismatch: %v", err)
		}
		log.Printf("Updated database - ID: %s Name: %s", upd.ID, upd.Name)

		colCount := rand.Intn(3) + 1
		cols := make([]colInfo, colCount)
		for j := 0; j < colCount; j++ {
			colID := id.Unique()
			colName := fmt.Sprintf("col_%d_%d", i, j)
			col, err := db.CreateCollection(dbID, colID, colName, []string{gowrite.ReadAny, gowrite.WriteAny}, true, true)
			if err != nil {
				t.Fatalf("CreateCollection: %v", err)
			}
			log.Printf("Created collection - ID: %s Name: %s", col.ID, col.Name)

			t.Cleanup(func(databaseID, collectionID string) func() {
				return func() {
					if err := db.DeleteCollection(databaseID, collectionID); err != nil {
						log.Printf("Cleanup DeleteCollection error: %v", err)
					} else {
						log.Printf("Deleted collection - ID: %s", collectionID)
					}
				}
			}(dbID, colID))

			gcol, err := db.GetCollection(dbID, colID)
			if err != nil || gcol.ID != colID || gcol.Name != colName {
				t.Fatalf("GetCollection mismatch: %v", err)
			}

			updColName := fmt.Sprintf("col_%d_%d_upd", i, j)
			updCol, err := db.UpdateCollection(dbID, colID, updColName, []string{gowrite.ReadAny, gowrite.WriteAny}, true, true)
			if err != nil || updCol.Name != updColName {
				t.Fatalf("UpdateCollection mismatch: %v", err)
			}
			log.Printf("Updated collection - ID: %s Name: %s", updCol.ID, updCol.Name)

			docCount := rand.Intn(3) + 1
			for k := 0; k < docCount; k++ {
				docID := id.Unique()
				docData := map[string]interface{}{
					"num":  rand.Intn(1000),
					"text": fmt.Sprintf("doc_%d_%d_%d", i, j, k),
				}
				doc, err := db.CreateDocument(dbID, colID, docID, docData, []string{gowrite.ReadAny, gowrite.WriteAny})
				if err != nil {
					t.Fatalf("CreateDocument: %v", err)
				}
				log.Printf("Created document - ID: %s Data: %v", doc.ID, doc.Data)

				t.Cleanup(func(databaseID, collectionID, documentID string) func() {
					return func() {
						if err := db.DeleteDocument(databaseID, collectionID, documentID); err != nil {
							log.Printf("Cleanup DeleteDocument error: %v", err)
						} else {
							log.Printf("Deleted document - ID: %s", documentID)
						}
					}
				}(dbID, colID, docID))

				gdoc, err := db.GetDocument(dbID, colID, docID)
				if err != nil || gdoc.ID != docID {
					t.Fatalf("GetDocument mismatch: %v", err)
				}

				updData := map[string]interface{}{
					"num":  rand.Intn(1000),
					"text": fmt.Sprintf("upd_%d_%d_%d", i, j, k),
				}
				udoc, err := db.UpdateDocument(dbID, colID, docID, updData, []string{gowrite.ReadAny, gowrite.WriteAny})
				if err != nil || udoc.ID != docID {
					t.Fatalf("UpdateDocument mismatch: %v", err)
				}
				log.Printf("Updated document - ID: %s", udoc.ID)
			}

			docs, err := db.ListDocuments(dbID, colID, []string{})
			if err != nil || len(docs) != docCount {
				t.Fatalf("ListDocuments mismatch: got %d want %d err=%v", len(docs), docCount, err)
			}

			cnt, err := db.CountDocuments(dbID, colID, []string{})
			if err != nil || cnt != docCount {
				t.Fatalf("CountDocuments mismatch: got %d want %d err=%v", cnt, docCount, err)
			}

			cols[j] = colInfo{id: colID, name: updColName, docCount: docCount}
		}

		colList, err := db.ListCollections(dbID)
		if err != nil || len(colList) != len(cols) {
			t.Fatalf("ListCollections mismatch: got %d want %d err=%v", len(colList), len(cols), err)
		}

		databases[i].id = dbID
		databases[i].name = dbName
		databases[i].newName = newName
		databases[i].collections = cols
	}

	listDb, err := db.ListDatabases()
	if err != nil || len(listDb) < dbCount {
		t.Fatalf("ListDatabases mismatch: got %d want >=%d err=%v", len(listDb), dbCount, err)
	}
}
