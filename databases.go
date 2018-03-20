package couchdb2go

import (
	"encoding/json"
	"net/http"
)

type Database interface {
	Exists(string) (*bool, error)
	Meta(string) (*DbMetaResponse, error)
	CreateDb(string) (*OkKoResponse, error)
	DeleteDb(string) (*OkKoResponse, error)
	CreateDocument(db string, doc interface{})
	CreateDocumentExtra(db string, doc interface{}, batch bool, fullCommit bool)
	Documents(db string, req *AllDocsRequest) (*AllDocsResponse, error)
	DocumentsWithIDs(db string, req *DocsWithIDsRequest) (*AllDocsResponse, error)
	Bulk(db string, docs []interface{}, newEdits bool) (*BulkResponse, error)
	Find(db string, req *FindRequest) (*FindResponse, error)
	SetIndex(db string, req *SetIndexRequest) (*SetIndexResponse, error)
	Index(db string) (*IndexResponse, error)
	Delete(db string, designDoc string, name string) (*OkKoResponse, error)
	Explain(db string, req *FindRequest) (*ExplainResponse, error)
	Changes(db string, req map[string]string) (*ChangesResponse, error)
	//ChangesContinuous(db string, queryReq map[string]string, inCh chan *DbResult, quitCh chan struct{}) (chan *DbResult, chan<- struct{}, error)
	ChangesContinuousRaw(db string, queryReq map[string]string, inCh chan *DbResult, quitCh chan struct{}) (chan *DbResult, chan<- struct{}, error)
	ChangesContinuousRawWithHeartBeat(db string, queryReq map[string]string, inCh chan *DbResult, inHeartBeatCh chan *HeartBeatResult, quitCh chan struct{}) (chan *DbResult, chan *HeartBeatResult, chan<- struct{}, error)
	ChangesContinuousBytes(db string, queryReq map[string]string) (*http.Response, error)
	Compact(db string) (*OkKoResponse, error)
	CompactDesignDoc(db string, ddoc string) (*OkKoResponse, error)
	EnsureFullCommit(db string) (*EnsureFullCommitResponse, error)
	ViewCleanup(db string) (*OkKoResponse, error)
	Security(db string) (*SecurityResponse, error)
	SetSecurity(db string, req *SecurityRequest) (*OkKoResponse, error)
	DoPurge(db string, req map[string]interface{}) (*DoPurgeResponse, error)
	MissingKeys(db string, req map[string]interface{}) (*MissingKeysResponse, error)
	RevsDiff(db string, req *RevsDiffRequest) (*RevsDiffResponse, error)
	RevsLimit(db string) (int, error)
	SetRevsLimit(db string, n int) (*OkKoResponse, error)
}

type DbMetaResponse struct {
	ErrorResponse
	CommittedUpdateSeq int    `json:"committed_update_seq"`
	CompactRunning     bool   `json:"compact_running"`
	DataSize           int    `json:"data_size"`
	DbName             string `json:"db_name"`
	DiskFormatVersion  int    `json:"disk_format_version"`
	DiskSize           int    `json:"disk_size"`
	DocCount           int    `json:"doc_count"`
	DocDelCount        int    `json:"doc_del_count"`
	InstanceStartTime  string `json:"instance_start_time"`
	PurgeSeq           int    `json:"purge_seq"`
	UpdateSeq          int    `json:"update_seq"`
}

type OkKoResponse struct {
	ErrorResponse
	Ok bool `json:"ok"`
}

type CreateDocumentResponse struct {
	ErrorResponse
	ID  string `json:"id"`
	Ok  bool   `json:"ok"`
	Rev string `json:"rev, omitempty"`
}

type AllDocsRequest struct {
	Conflicts        bool   `json:"conflicts, omitempty"`
	Descending       bool   `json:"descending, omitempty"`
	EndKey           string `json:"endkey, omitempty"`
	End_key          string `json:"end_key, omitempty"`
	EndKey_DocID     string `json:"endkey_docid, omitempty"`
	End_Key_Doc_ID   string `json:"end_key_doc_id, omitempty"`
	IncludeDocs      bool   `json:"include_docs, omitempty"`
	InclusiveEnd     bool   `json:"inclusive_end, omitempty"`
	Key              string `json:"key, omitempty"`
	Keys             string `json:"keys, omitempty"`
	Limit            int    `json:"limit, omitempty"`
	Skip             int    `json:"skip, omitempty"`
	Stale            string `json:"stale, omitempty"`
	StartKey         string `json:"startkey, omitempty"`
	Start_Key        string `json:"start_key, omitempty"`
	StartKey_DocID   string `json:"startkey_docid, omitempty"`
	Start_Key_Doc_ID string `json:"start_key_doc_id, omitempty"`
	UpdateSeq        bool   `json:"update_seq, omitempty"`
}

type AllDocsResponse struct {
	ErrorResponse
	Offset int `json:"offset"`
	Rows []struct {
		ID  string `json:"id"`
		Key string `json:"key"`
		Value struct {
			Rev string `json:"rev"`
		} `json:"value"`
	} `json:"rows"`
	TotalRows int `json:"total_rows"`
}

type DocsWithIDsRequest struct {
	Keys []string `json:"keys"`
}

type BulkResponse []CreateDocumentResponse

type FindRequest struct {
	Selector map[string]interface{} `json:"selector"`
	Fields   []string               `json:"fields"`
	Sort []struct {
		Year string `json:"year"`
	} `json:"sort"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type FindResponse struct {
	ErrorResponse
	Docs []map[string]interface{} `json:"docs"`
}

type SetIndexRequest struct {
	Index map[string]interface{} `json:"index"`
	Ddoc  string                 `json:"ddocs"`
	Name  string                 `json:"name"`
}

type SetIndexResponse struct {
	ErrorResponse
	Result string `json:"result"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}

type IndexResponse struct {
	ErrorResponse
	TotalRows int `json:"total_rows"`
	Indexes []struct {
		Ddoc interface{} `json:"ddoc"`
		Name string      `json:"name"`
		Type string      `json:"type"`
		Def struct {
			Fields []struct {
				ID string `json:"_id"`
			} `json:"fields"`
		} `json:"def"`
	} `json:"indexes"`
}

type ExplainResponse struct {
	ErrorResponse
	Dbname string `json:"dbname"`
	Index struct {
		Ddoc string `json:"ddoc"`
		Name string `json:"name"`
		Type string `json:"type"`
		Def struct {
			Fields []struct {
				Year string `json:"year"`
			} `json:"fields"`
		} `json:"def"`
	} `json:"index"`
	Selector struct {
		Year struct {
			Gt int `json:"$gt"`
		} `json:"year"`
	} `json:"selector"`
	Opts struct {
		UseIndex []interface{} `json:"use_index"`
		Bookmark string        `json:"bookmark"`
		Limit    int           `json:"limit"`
		Skip     int           `json:"skip"`
		Sort struct {
		} `json:"sort"`
		Fields    []string `json:"fields"`
		R         []int    `json:"r"`
		Conflicts bool     `json:"conflicts"`
	} `json:"opts"`
	Limit  int      `json:"limit"`
	Skip   int      `json:"skip"`
	Fields []string `json:"fields"`
	Range struct {
		StartKey []int `json:"start_key"`
		EndKey []struct {
		} `json:"end_key"`
	} `json:"range"`
}

type ChangesRequest struct {
	DocsIds         []string `json:"docs_ids"`
	Conflicts       bool     `json:"conflicts"`
	Descending      bool     `json:"descending"`
	Feed            string   `json:"feed"`
	Filter          string   `json:"filter"`
	HeartBeat       int      `json:"heartbeat"`
	IncludeDocs     bool     `json:"include_docs"`
	Attachments     bool     `json:"attachments"`
	AttEncodingInfo bool     `json:"att_encoding_info"`
	LastEventId     int      `json:"last-event-id"`
	Limit           int      `json:"limit"`
	Since           int      `json:"since"`
	Style           string   `json:"string"`
	Timeout         int      `json:"timeout"`
	View            string   `json:"view"`
}

type LastSeq struct {
	Last_Seq *string `json:"last_seq"`
}

type DbResult struct {
	LastSeq
	*ErrorResponse
	Changes []struct {
		Rev string `json:"rev"`
	} `json:"changes"`
	ID      string                 `json:"id"`
	Seq     json.Number            `json:"seq,Number"`
	Deleted bool                   `json:"deleted,omitempty"`
	DbName  string                 `json:"database,omitempty"`
	Doc     map[string]interface{} `json:"doc,omitempty"`
}

type Deleted struct {
	Deleted bool   `json:"_deleted"`
	Id      string `json:"_id"`
	Rev     string `json:"_rev"`
}

type HeartBeatResult struct{}

type ChangesResponse struct {
	ErrorResponse
	LastSeq int         `json:"last_seq"`
	Results []*DbResult `json:"results"`
}

type EnsureFullCommitResponse struct {
	OkKoResponse
	InstanceStartTime string `json:"instance_start_time"`
}

type SecurityRequest struct {
	SecurityResponse
}

type SecurityResponse struct {
	ErrorResponse
	Admins struct {
		Names []string `json:"names"`
		Roles []string `json:"roles"`
	} `json:"admins"`
	Members struct {
		Names []string `json:"names"`
		Roles []string `json:"roles"`
	} `json:"members"`
}

type DoPurgeResponse struct {
	ErrorResponse
	PurgeSeq int `json:"purge_seq"`
	Purged struct {
		C6114C65E295552Ab1019E2B046B10E []string `json:"c6114c65e295552ab1019e2b046b10e"`
	} `json:"purged"`
}

type MissingKeysResponse struct {
	ErrorResponse
	MissedRevs []map[string]interface{} `json:"missed_revs"`
}

type RevsDiffRequest map[string][]string

type RevsDiffResponse map[string]struct {
	ErrorResponse
	Missing           []string `json:"missing"`
	PossibleAncestors []string `json:"possible_ancestors"`
}
