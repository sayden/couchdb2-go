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
			false,
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
		"timeout":   "1000",
		"heartbeat": "10000",
		"feed":      "continuous",
		"since":     "0",
	}

	in, quit, err := d.ChangesContinuous("test", req, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	out := make(chan []byte, 5)
	documents := DocumentsClient{
		CouchDb2ConnDetails: d.CouchDb2ConnDetails,
	}

	go func(in <-chan *DbResult, out chan<- []byte) {
		var i int

		//Take elements from channel of incoming ID'c
		for result := range in {
			doc, err := documents.Document("test", result.ID)
			if err != nil {
				t.Fatal(err)
			}

			out <- doc
			i++

			if i >= 100 {
				quit <- struct{}{}
				break
			}
		}

		close(out)
	}(in, out)

	for doc := range out {
		pretty.Printf("%# v", string(doc))
	}

}
