package couchdb2go

import (
	"os"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
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

type mockedDb struct {
}

func (*mockedDb) Exists(string) (*bool, error) {
	panic("implement me")
}

func (*mockedDb) Meta(string) (*DbMetaResponse, error) {
	panic("implement me")
}

func (*mockedDb) CreateDb(string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) DeleteDb(string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) CreateDocument(db string, doc interface{}) {
	panic("implement me")
}

func (*mockedDb) CreateDocumentExtra(db string, doc interface{}, batch bool, fullCommit bool) {
	panic("implement me")
}

func (*mockedDb) Documents(db string, req *AllDocsRequest) (*AllDocsResponse, error) {
	panic("implement me")
}

func (*mockedDb) DocumentsWithIDs(db string, req *DocsWithIDsRequest) (*AllDocsResponse, error) {
	panic("implement me")
}

func (*mockedDb) Bulk(db string, docs []interface{}, newEdits bool) (*BulkResponse, error) {
	panic("implement me")
}

func (*mockedDb) Find(db string, req *FindRequest) (*FindResponse, error) {
	panic("implement me")
}

func (*mockedDb) SetIndex(db string, req *SetIndexRequest) (*SetIndexResponse, error) {
	panic("implement me")
}

func (*mockedDb) Index(db string) (*IndexResponse, error) {
	panic("implement me")
}

func (*mockedDb) Delete(db string, designDoc string, name string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) Explain(db string, req *FindRequest) (*ExplainResponse, error) {
	panic("implement me")
}

func (*mockedDb) Changes(db string, req map[string]string) (*ChangesResponse, error) {
	panic("implement me")
}

func (*mockedDb) ChangesContinuousRaw(db string, queryReq map[string]string, inCh chan *DbResult, quitCh chan struct{}) (chan *DbResult, chan<- struct{}, error) {
	inCh <- &DbResult{
		ErrorResponse: &ErrorResponse{
			ErrorS: "forbidden",
			Reason: "Needs reader access",
		},
	}

	return inCh, quitCh, nil
}

func (*mockedDb) ChangesContinuousRawWithHeartBeat(db string, queryReq map[string]string, inCh chan *DbResult, inHeartBeatCh chan *HeartBeatResult, quitCh chan struct{}) (chan *DbResult, chan *HeartBeatResult, chan<- struct{}, error) {
	inHeartBeatCh <- &HeartBeatResult{}

	return inCh, inHeartBeatCh, quitCh, nil
}

func (*mockedDb) Compact(db string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) CompactDesignDoc(db string, ddoc string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) EnsureFullCommit(db string) (*EnsureFullCommitResponse, error) {
	panic("implement me")
}

func (*mockedDb) ViewCleanup(db string) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) Security(db string) (*SecurityResponse, error) {
	panic("implement me")
}

func (*mockedDb) SetSecurity(db string, req *SecurityRequest) (*OkKoResponse, error) {
	panic("implement me")
}

func (*mockedDb) DoPurge(db string, req map[string]interface{}) (*DoPurgeResponse, error) {
	panic("implement me")
}

func (*mockedDb) MissingKeys(db string, req map[string]interface{}) (*MissingKeysResponse, error) {
	panic("implement me")
}

func (*mockedDb) RevsDiff(db string, req *RevsDiffRequest) (*RevsDiffResponse, error) {
	panic("implement me")
}

func (*mockedDb) RevsLimit(db string) (int, error) {
	panic("implement me")
}

func (*mockedDb) SetRevsLimit(db string, n int) (*OkKoResponse, error) {
	panic("implement me")
}

func TestDatabasesClient_ChangesContinuousRaw(t *testing.T) {
	t.Run("Common test", func(t *testing.T) {
		dT := newDatabase()

		d := dT.(*DatabasesClient)
		d.CouchDb2ConnDetails.Client.Timeout = 0

		req := map[string]string{
			"heartbeat":    "500",
			"feed":         "continuous",
			"since":        "0",
			"include_docs": "true",
		}

		in, quit, err := d.ChangesContinuousRaw("test", req, nil, nil)
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
	})

	t.Run("Error returned", func(t *testing.T) {
		dT := &mockedDb{}

		in := make(chan *DbResult, 1)
		quit := make(chan struct{}, 0)

		if _, _, err := dT.ChangesContinuousRaw("test", make(map[string]string), in, quit); err != nil {
			t.Fatal(err)
		}
		close(quit)
		close(in)

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
		}
	})
}

func TestDatabasesClient_ChangesContinuousRawWithHeartBeat(t *testing.T) {
	t.Run("Heartbeat returned", func(t *testing.T) {
		dT := &mockedDb{}

		in := make(chan *DbResult, 1)
		inHeartbeat := make(chan *HeartBeatResult, 10)
		quit := make(chan struct{}, 0)

		if _, _, _, err := dT.ChangesContinuousRawWithHeartBeat("test", make(map[string]string), in, inHeartbeat, quit); err != nil {
			t.Fatal(err)
		}

		close(quit)
		close(in)
		close(inHeartbeat)

		result := <-inHeartbeat
		assert.Equal(t, &HeartBeatResult{}, result)
		assert.Equal(t, 0, len(inHeartbeat))
	})
}
