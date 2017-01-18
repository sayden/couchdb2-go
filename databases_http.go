package couchdb2_goclient

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type DatabasesClient struct {
	*CouchDb2ConnDetails
}

func (d *DatabasesClient) GetConnection() *CouchDb2ConnDetails {
	return d.CouchDb2ConnDetails
}

func (d *DatabasesClient) Exists(string) (*bool, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Meta(string) (*DbMetaResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) CreateDb(string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) DeleteDb(string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) CreateDocument(db string, doc interface{}) {
	panic("not implemented")
}

func (d *DatabasesClient) CreateDocumentExtra(db string, doc interface{}, batch bool, fullCommit bool) {
	panic("not implemented")
}

func (d *DatabasesClient) Documents(db string, req *AllDocsRequest) (*AllDocsResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) DocumentsWithIDs(db string, req *DocsWithIDsRequest) (*AllDocsResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Bulk(db string, docs []interface{}, newEdits bool) (*BulkResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Find(db string, req *FindRequest) (*FindResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) SetIndex(db string, req *SetIndexRequest) (*SetIndexResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Index(db string) (*IndexResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Delete(db string, designDoc string, name string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Explain(db string, req *FindRequest) (*ExplainResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Changes(db string, queryReq map[string]string) (*ChangesResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) ChangesContinuous(db string, queryReq map[string]string, out chan *DbResult, quit chan struct{}) (chan *DbResult, chan<- struct{}, error) {
	var query string
	for k, v := range queryReq {
		query = fmt.Sprintf("%s%s=%v&", query, k, v)
	}

	//remove last '&' on query
	query = strings.TrimSuffix(query, "&")

	if d.Client == nil {
		return nil, nil, errors.New("You must set an HTTP Client to make requests. Current client is nil")
	}

	//fmt.Printf("%s://%s/_changes?%s\n", d.protocol, d.Address, query)
	fmt.Printf("Attempting connection to %s://%s/%s/_changes?%s\n", d.protocol, d.Address, db, query)
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
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Recovered error: %s\n", err)
				out <- &DbResult{
					ErrorResponse: &ErrorResponse{
						ErrorS: fmt.Sprintf("%s\n", err),
						Reason: "Error recovered",
					},
				}
			}
		}()

		for scanner.Scan() {
			line = scanner.Bytes()

			if string(line) == "" || string(line) == "\n" {
				fmt.Println("Got empty result or newline")
				continue
			} else if strings.Contains(string(line), "error") {
				fmt.Printf("ERROR: received from db: %s\n", string(line))

				var errRes ErrorResponse

				err = json.Unmarshal(line, &errRes)
				if err != nil {
					fmt.Printf("ERROR: received from db: %s\n", err.Error())
					out <- &DbResult{
						ErrorResponse: &ErrorResponse{
							ErrorS: err.Error(),
							Reason: "last_seq found. This usually means that the connections has finished",
						},
					}
				}

				out <- &DbResult{
					ErrorResponse:&errRes,
				}

				break
			} else if strings.Contains(string(line), "last_seq") {
				err = json.Unmarshal(line, &last)

				if err != nil {
					fmt.Printf("ERROR: reading json: %s\n", err.Error())
					out <- &DbResult{
						ErrorResponse: &ErrorResponse{
							ErrorS: err.Error(),
							Reason: "last_seq found. This usually means that the connections has finished",
						},
					}
				}

				break
			} else {
				var result DbResult

				err = json.Unmarshal(line, &result)
				if err != nil {
					fmt.Println("ERROR:", result.Error())
					out <- &DbResult{
						ErrorResponse: &ErrorResponse{
							ErrorS: err.Error(),
						},
					}
				}

				select {
				case <-quit:
					break
				case out <- &result:
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("ERROR: reading: %s\n", err.Error())
			out <- &DbResult{
				ErrorResponse: &ErrorResponse{
					ErrorS: "Error scanning input",
					Reason: err.Error(),
				},
			}
			close(out)
		}

		fmt.Println("Quitting")
	}()

	return out, quit, nil
}

func (d *DatabasesClient) Compact(db string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) CompactDesignDoc(db string, ddoc string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) EnsureFullCommit(db string) (*EnsureFullCommitResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) ViewCleanup(db string) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) Security(db string) (*SecurityResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) SetSecurity(db string, req *SecurityRequest) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) DoPurge(db string, req map[string]interface{}) (*DoPurgeResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) MissingKeys(db string, req map[string]interface{}) (*MissingKeysResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) RevsDiff(db string, req *RevsDiffRequest) (*RevsDiffResponse, error) {
	panic("not implemented")
}

func (d *DatabasesClient) RevsLimit(db string) (int, error) {
	panic("not implemented")
}

func (d *DatabasesClient) SetRevsLimit(db string, n int) (*OkKoResponse, error) {
	panic("not implemented")
}

func NewDatabase(timeout time.Duration, addr string, user, pass string, secure bool) (dat Database) {
	dat = &DatabasesClient{
		CouchDb2ConnDetails: NewConnection(timeout, addr, user, pass, secure),
	}

	return
}
