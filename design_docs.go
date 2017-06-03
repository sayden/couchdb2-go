package couchdb2_goclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DesignDocuments interface {
	View(db string, req ViewRequest) (res ViewResponse, err error)
	ViewWithType(db string, req ViewRequest, t interface{}) (res ViewResponse, err error)
}

type DesignDocsClient struct {
	*CouchDb2ConnDetails
}

func (d *DesignDocsClient) View(db string, req ViewRequest) (res ViewResponse, err error) {
	var byt []byte

	byt, err = d.bytesRequester(http.MethodGet, fmt.Sprintf(`/%s/_design/%s/_view/%s?%s`,
		db, req.DesignDocName, req.ViewName, buildURLParams(req.QueryParams)), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(byt, &res)

	return
}

func (d *DesignDocsClient) ViewWithType(db string, req ViewRequest, t interface{}) (res ViewResponse, err error) {
	res, err = d.View(db, req)

	return
}

type ViewRequest struct {
	QueryParams   map[string]string
	DesignDocName string
	ViewName      string
}

type ViewResponse struct {
	TotalRows int `json:"total_rows"`
	Offset    int `json:"offset"`
	Rows      []map[string]interface{}
}

func NewDesignDocuments(timeout time.Duration, addr string, user, pass string, secured bool) (doc DesignDocuments) {
	doc = &DesignDocsClient{
		CouchDb2ConnDetails: NewConnection(timeout, addr, user, pass, secured),
	}

	return
}

func NewDesignDocumentsWithConnection(conn *CouchDb2ConnDetails) (doc DesignDocuments) {
	doc = &DesignDocsClient{
		CouchDb2ConnDetails: conn,
	}

	return
}
