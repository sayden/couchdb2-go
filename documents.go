package couchdb2_goclient

import (
	"fmt"
	"net/http"
	"time"
)

type Documents interface {
	Headers(db string, id string) (size *int, rev string, err error)
	Document(db string, id string) ([]byte, error)
	UpdateDocument(db string, id string, req map[string]interface{}) (*OkKoResponse, error)
	DeleteDocument(db string, id string) (*OkKoResponse, error)
}

func NewDocumentsWithConnection(conn *CouchDb2ConnDetails) (doc Documents) {
	doc = &documents{conn}

	return
}

func NewDocuments(t time.Duration, addr string, user, pass string) (doc Documents) {
	doc = &documents{
		CouchDb2ConnDetails: NewConnection(t, addr, user, pass, true),
	}

	return
}

type documents struct {
	*CouchDb2ConnDetails
}

func (d *documents) Headers(db string, id string) (size *int, rev string, err error) {
	panic("not implemented")
}

func (d *documents) Document(db string, id string) ([]byte, error) {
	return d.bytesRequester(http.MethodGet, fmt.Sprintf("/%s/%s", db, id), nil)
}

func (d *documents) UpdateDocument(db string, id string, req map[string]interface{}) (*OkKoResponse, error) {
	panic("not implemented")
}

func (d *documents) DeleteDocument(db string, id string) (*OkKoResponse, error) {
	panic("not implemented")
}
