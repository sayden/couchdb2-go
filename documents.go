package couchdb2_goclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Documents interface {
	Headers(db string, id string) (size *int, rev string, err error)
	Document(db string, id string) ([]byte, error)
	DocumentWithType(db string, id string, t interface{}) error
	UpdateDocument(db string, id string, req map[string]interface{}) (*OkKoResponse, error)
	DeleteDocument(db string, id string) (*OkKoResponse, error)
}

type DocumentsClient struct {
	*CouchDb2ConnDetails
}

func (d *DocumentsClient) Headers(db string, id string) (size *int, rev string, err error) {
	panic("not implemented")
}

func (d *DocumentsClient) DocumentWithType(db string, id string, t interface{}) (err error) {
	var byt []byte
	if byt, err = d.bytesRequester(http.MethodGet, fmt.Sprintf("/%s/%s", db, id), nil); err != nil {
		return
	}

	var dbresult DbResult
	if err = json.Unmarshal(byt, &dbresult); err != nil {
		return
	}

	if dbresult.ErrorResponse != nil {
		err = dbresult.ErrorResponse
		return
	}

	err = json.Unmarshal(byt, &t)

	return
}

func (d *DocumentsClient) Document(db string, id string) ([]byte, error) {
	return d.bytesRequester(http.MethodGet, fmt.Sprintf("/%s/%s", db, id), nil)
}

func (d *DocumentsClient) UpdateDocument(db string, id string, req map[string]interface{}) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *DocumentsClient) DeleteDocument(db string, id string) (*OkKoResponse, error) {
	panic("not implemented")
}

func NewDocumentsWithConnection(conn *CouchDb2ConnDetails) (doc Documents) {
	doc = &DocumentsClient{conn}

	return
}

func NewDocuments(timeout time.Duration, addr string, user, pass string, secured bool) (doc Documents) {
	doc = &DocumentsClient{
		CouchDb2ConnDetails: NewConnection(timeout, addr, user, pass, secured),
	}

	return
}
