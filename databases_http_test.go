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

func TestDatabasesClient_ChangesContinuousRaw(t *testing.T) {
	dT := newDatabase()

	d := dT.(*DatabasesClient)
	d.CouchDb2ConnDetails.Client.Timeout = 0

	req := map[string]string{
		"heartbeat":    "500",
		"feed":         "continuous",
		"since":        "95418662",
		"include_docs": "true",
	}

	in, quit, err := d.ChangesContinuousRaw("maxi_billing", req, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var i int

	//Take elements from channel of incoming ID'c
	for result := range in {
		if result.ErrorResponse != nil {
			t.Logf("Error returned from db: %v\n", result.ErrorResponse.Error())
		}

		if result.Last_Seq != nil {
			t.Logf("Last seq reached: %d\n", result.Last_Seq)
			break
		}

		pretty.Printf("%# v", result.Doc)

		i++
		if i >= 1 {
			quit <- struct{}{}
			break
		}
	}

	if i <= 0 {
		t.Fail()
	}
}
