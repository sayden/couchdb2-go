package couchdb2_goclient

import (
	"os"
	"testing"
	"time"

	"github.com/kr/pretty"
)

var testDatabase Database

func newDatabase() Database {
	if testDatabase == nil {
		testDatabase = NewDatabase(time.Minute,
			os.Getenv("COUCHDB_ADDRESS"),
			os.Getenv("COUCHDB_USER"),
			os.Getenv("COUCHDB_PASSWORD"),
			os.Getenv("SECURED") == "true",
		)
	}
	return testDatabase
}

func TestDatabasesClient_Changes(t *testing.T) {
	dT := newDatabase()

	d := dT.(*DatabasesClient)
	d.CouchDb2ConnDetails.Client.Timeout = time.Second * 20

	req := map[string]string{
		"timeout": "1",
		"feed":    "continuous",
	}

	res, err := d.Changes("test", req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func TestDatabasesClient_Changes_ForeverPolling(t *testing.T) {
	dT := newDatabase()

	d := dT.(*DatabasesClient)
	d.CouchDb2ConnDetails.Client.Timeout = 0

	req := map[string]string{
		"timeout":   "1000",
		"heartbeat": "10000",
		"feed":      "continuous",
		"since":     "0",
	}

	out, quit, err := d.ChangesContinuous("test", req, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var i int
	for result := range out {
		pretty.Printf("%# v\n", result)
		i++

		if i > 5 {
			quit <- struct{}{}
			break
		}
	}
}

func TestDatabasesClient_Changes_FullDocument(t *testing.T) {
	dT := newDatabase()

	d := dT.(*DatabasesClient)
	d.CouchDb2ConnDetails.Client.Timeout = 0

	req := map[string]string{
		"heartbeat": "10000",
		"feed":      "continuous",
		"since":     "0",
	}

	in, quit, err := d.ChangesContinuous("test", req, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	documents := DocumentsClient{
		CouchDb2ConnDetails: d.CouchDb2ConnDetails,
	}

	var i int

	//Take elements from channel of incoming ID'c
	for result := range in {
		t.Log("ASDFAsf")
		if result.ErrorResponse != nil {
			t.Fatalf("Test error: %s\n", result.ErrorResponse.Error())
		}

		doc, err := documents.Document("test", result.ID)
		if err != nil {
			t.Fatal(err)
		}

		pretty.Printf("%# v", string(doc))
		i++

		if i >= 100 {
			quit <- struct{}{}
			break
		}
	}
}
