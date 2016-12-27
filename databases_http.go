package couchdb2_goclient

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type databasesClient struct {
	*CouchDb2ConnDetails
	client *http.Client
}

func (d *databasesClient) Exists(string) (*bool, error) {
	panic("not implemented")
}

func (d *databasesClient) Meta(string) (*DbMetaResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) CreateDb(string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) DeleteDb(string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) CreateDocument(db string, doc interface{}) {
	panic("not implemented")
}

func (d *databasesClient) CreateDocumentExtra(db string, doc interface{}, batch bool, fullCommit bool) {
	panic("not implemented")
}

func (d *databasesClient) Documents(db string, req *AllDocsRequest) (*AllDocsResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) DocumentsWithIDs(db string, req *DocsWithIDsRequest) (*AllDocsResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Bulk(db string, docs []interface{}, newEdits bool) (*BulkResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Find(db string, req *FindRequest) (*FindResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) SetIndex(db string, req *SetIndexRequest) (*SetIndexResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Index(db string) (*IndexResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Delete(db string, designDoc string, name string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Explain(db string, req *FindRequest) (*ExplainResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Changes(db string, queryReq map[string]string) (*ChangesResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) ChangesContinuous(db string, queryReq map[string]string, out chan *DbResult, quit chan struct{}) (chan *DbResult, chan<- struct{}, error) {
	var query string
	for k, v := range queryReq {
		query = fmt.Sprintf("%s%s=%v&", query, k, v)
	}

	if d.Client == nil {
		return nil, nil, errors.New("You must set an HTTP Client to make requests. Current client is nil")
	}

	//fmt.Printf("%s://%s/_changes?%s\n", d.protocol, d.Address, query)

	reqHttp, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s/%s/_changes?%s", d.protocol, d.Address, db, query), nil)
	if err != nil {
		return nil, nil, err
	}

	if d.Username != "" && d.Password != "" {
		reqHttp.SetBasicAuth(d.Username, d.Password)
	}

	reqHttp.Header.Add("Accept", "application/json")
	reqHttp.Header.Add("Content-Type", "application/json")

	httpRes, err := d.Client.Do(reqHttp)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(httpRes.Body)

	var line []byte
	var last LastSeq

	if out == nil {
		out = make(chan *DbResult, 100)
	}

	if quit == nil {
		quit = make(chan struct{}, 1)
	}

	go func() {
		defer httpRes.Body.Close()

		for scanner.Scan() {
			line = scanner.Bytes()

			if string(line) == "" || string(line) == "\n" {
				continue
			} else if strings.Contains(string(line), "last_seq") {
				err = json.Unmarshal(line, &last)

				if err != nil {
					//TODO better error handling
					fmt.Println("ERROR:", err)
				}

				break
			} else {
				var result DbResult

				err = json.Unmarshal(line, &result)
				if result.ErrorS != "" {
					fmt.Println("ERROR:", result.Error())
				}

				result.DbName = db

				select {
				case <-quit:
					break
				case out <- &result:
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		fmt.Println("Quitting")
	}()

	return out, quit, nil
}

func (d *databasesClient) Compact(db string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) CompactDesignDoc(db string, ddoc string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) EnsureFullCommit(db string) (*EnsureFullCommitResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) ViewCleanup(db string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) Security(db string) (*SecurityResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) SetSecurity(db string, req *SecurityRequest) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) DoPurge(db string, req map[string]interface{}) (*DoPurgeResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) MissingKeys(db string, req map[string]interface{}) (*MissingKeysResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) RevsDiff(db string, req *RevsDiffRequest) (*RevsDiffResponse, error) {
	panic("not implemented")
}

func (d *databasesClient) RevsLimit(db string) (int, error) {
	panic("not implemented")
}

func (d *databasesClient) SetRevsLimit(db string, n int) (*OkKoResponse, error) {
	panic("not implemented")
}

func NewDatabase(t time.Duration, addr string, user, pass string) (dat Database) {
	dat = &databasesClient{
		CouchDb2ConnDetails: NewConnection(t, addr, user, pass, true),
	}

	return
}
