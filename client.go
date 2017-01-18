package couchdb2_goclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Meta() (*MetaResponse, error)
	ActiveTasks() (*ActiveTasksResponse, error)
	AllDbs() (*AllDbsResponse, error)
	DbUpdates(r *DbUpdatesRequest) (*DbUpdatesResponse, error)
	Membership() (*MembershipResponse, error)
	Log(*LogRequest) (*LogResponse, error)
	Replicate() (*ReplicateResponse, error)
	Restart() (*RestartResponse, error)
	Stats() (*StatsResponse, error)
	UUIDs(uint8) (*UUIDsResponse, error)
	Config() (*ConfigResponse, error)
	Config
}

type Config interface {
	Section(string) (*ConfigSectionResponse, error)
	// TODO Config interface
	//Key(s, k string)
	//SetKey(s, k string)
	//DeleteKey(s, k string)
}

type MetaResponse struct {
	ErrorResponse
	CouchDB string `json:"couchdb"`
	UUID    string `json:"uuid"`
	Vendor  struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"vendor"`
	Version string `json:"version"`
}

type ActiveTasksResponse []struct {
	ErrorResponse
	ChangesDone           int         `json:"changes_done,omitempty"`
	Database              string      `json:"database,omitempty"`
	Pid                   string      `json:"pid"`
	Progress              int         `json:"progress"`
	StartedOn             int         `json:"started_on"`
	TotalChanges          int         `json:"total_changes,omitempty"`
	Type                  string      `json:"type"`
	UpdatedOn             int         `json:"updated_on"`
	DesignDocument        string      `json:"design_document,omitempty"`
	CheckpointedSourceSeq int         `json:"checkpointed_source_seq,omitempty"`
	Continuous            bool        `json:"continuous,omitempty"`
	DocID                 interface{} `json:"doc_id,omitempty"`
	DocWriteFailures      int         `json:"doc_write_failures,omitempty"`
	DocsRead              int         `json:"docs_read,omitempty"`
	DocsWritten           int         `json:"docs_written,omitempty"`
	MissingRevisionsFound int         `json:"missing_revisions_found,omitempty"`
	ReplicationID         string      `json:"replication_id,omitempty"`
	RevisionsChecked      int         `json:"revisions_checked,omitempty"`
	Source                string      `json:"source,omitempty"`
	SourceSeq             int         `json:"source_seq,omitempty"`
	Target                string      `json:"target,omitempty"`
}

type AllDbsResponse []string

type DbUpdatesRequest struct {
	Feed      string `json:"feed"`
	Timeout   int    `json:"timeout"`
	HeartBeat bool   `json:"heartbeat"`
}

type DbUpdatesResponse struct {
	ErrorResponse
	DbName string `json:"db_name"`
	Ok     bool   `json:"ok"`
	Type   string `json:"type"`
}

type MembershipResponse struct {
	ErrorResponse
	AllNodes     []string `json:"all_nodes"`
	ClusterNodes []string `json:"cluster_nodes"`
}

type LogRequest struct {
	ErrorResponse
	Bytes  uint32 `json:"bytes"`
	Offset uint32 `json:"offset"`
}
type LogResponse string

type ReplicateRequest struct {
	Cancel       bool     `json:"cancel"`
	Continuous   bool     `json:"continuous"`
	CreateTarget bool     `json:"create_target"`
	DocIds       []string `json:"doc_ids"`
	Proxy        string   `json:"proxy"`
	Source       string   `json:"source"`
	Target       string   `json:"target"`
}

type ReplicateResponse struct {
	ErrorResponse
	History []struct {
		DocWriteFailures int    `json:"doc_write_failures"`
		DocsRead         int    `json:"docs_read"`
		DocsWritten      int    `json:"docs_written"`
		EndLastSeq       int    `json:"end_last_seq"`
		EndTime          string `json:"end_time"`
		MissingChecked   int    `json:"missing_checked"`
		MissingFound     int    `json:"missing_found"`
		RecordedSeq      int    `json:"recorded_seq"`
		SessionID        string `json:"session_id"`
		StartLastSeq     int    `json:"start_last_seq"`
		StartTime        string `json:"start_time"`
	} `json:"history"`
	Ok                   bool   `json:"ok"`
	ReplicationIDVersion int    `json:"replication_id_version"`
	SessionID            string `json:"session_id"`
	SourceLastSeq        int    `json:"source_last_seq"`
}

type RestartResponse struct {
	ErrorResponse
	Ok bool `json:"ok"`
}

type StatsResponse struct {
	ErrorResponse
	Couchdb struct {
		RequestTime struct {
			Current     float64 `json:"current"`
			Description string  `json:"description"`
			Max         float64 `json:"max"`
			Mean        float64 `json:"mean"`
			Min         float64 `json:"min"`
			Stddev      float64 `json:"stddev"`
			Sum         float64 `json:"sum"`
		} `json:"request_time"`
	} `json:"couchdb"`
}

type UUIDsResponse struct {
	ErrorResponse
	Uuids []string `json:"uuids"`
}

type ConfigResponse struct {
	ErrorResponse
	Attachments struct {
		CompressibleTypes string `json:"compressible_types"`
		CompressionLevel  string `json:"compression_level"`
	} `json:"attachments"`
	CouchHttpdAuth struct {
		AuthCacheSize          string `json:"auth_cache_size"`
		AuthenticationDb       string `json:"authentication_db"`
		AuthenticationRedirect string `json:"authentication_redirect"`
		RequireValidUser       string `json:"require_valid_user"`
		Timeout                string `json:"timeout"`
	} `json:"couch_httpd_auth"`
	Couchdb struct {
		DatabaseDir            string `json:"database_dir"`
		DelayedCommits         string `json:"delayed_commits"`
		MaxAttachmentChunkSize string `json:"max_attachment_chunk_size"`
		MaxDbsOpen             string `json:"max_dbs_open"`
		MaxDocumentSize        string `json:"max_document_size"`
		OsProcessTimeout       string `json:"os_process_timeout"`
		URIFile                string `json:"uri_file"`
		UtilDriverDir          string `json:"util_driver_dir"`
		ViewIndexDir           string `json:"view_index_dir"`
	} `json:"couchdb"`
	Daemons struct {
		AuthCache        string `json:"auth_cache"`
		DbUpdateNotifier string `json:"db_update_notifier"`
		ExternalManager  string `json:"external_manager"`
		Httpd            string `json:"httpd"`
		QueryServers     string `json:"query_servers"`
		StatsAggregator  string `json:"stats_aggregator"`
		StatsCollector   string `json:"stats_collector"`
		Uuids            string `json:"uuids"`
		ViewManager      string `json:"view_manager"`
	} `json:"daemons"`
	Httpd struct {
		AllowJsonp             string `json:"allow_jsonp"`
		AuthenticationHandlers string `json:"authentication_handlers"`
		BindAddress            string `json:"bind_address"`
		DefaultHandler         string `json:"default_handler"`
		MaxConnections         string `json:"max_connections"`
		Port                   string `json:"port"`
		SecureRewrites         string `json:"secure_rewrites"`
		VhostGlobalHandlers    string `json:"vhost_global_handlers"`
	} `json:"httpd"`
	HttpdDbHandlers struct {
		Changes     string `json:"_changes"`
		Compact     string `json:"_compact"`
		Design      string `json:"_design"`
		TempView    string `json:"_temp_view"`
		ViewCleanup string `json:"_view_cleanup"`
	} `json:"httpd_db_handlers"`
	HttpdDesignHandlers struct {
		Info    string `json:"_info"`
		List    string `json:"_list"`
		Rewrite string `json:"_rewrite"`
		Show    string `json:"_show"`
		Update  string `json:"_update"`
		View    string `json:"_view"`
	} `json:"httpd_design_handlers"`
	HttpdGlobalHandlers struct {
		NAMING_FAILED string `json:"/"`
		ActiveTasks   string `json:"_active_tasks"`
		AllDbs        string `json:"_all_dbs"`
		Config        string `json:"_config"`
		Log           string `json:"_log"`
		Oauth         string `json:"_oauth"`
		Replicate     string `json:"_replicate"`
		Restart       string `json:"_restart"`
		Session       string `json:"_session"`
		Stats         string `json:"_stats"`
		Utils         string `json:"_utils"`
		Uuids         string `json:"_uuids"`
		FaviconIco    string `json:"favicon.ico"`
	} `json:"httpd_global_handlers"`
	Log struct {
		File        string `json:"file"`
		IncludeSasl string `json:"include_sasl"`
		Level       string `json:"level"`
	} `json:"log"`
	QueryServerConfig struct {
		ReduceLimit string `json:"reduce_limit"`
	} `json:"query_server_config"`
	QueryServers struct {
		Javascript string `json:"javascript"`
	} `json:"query_servers"`
	Replicator struct {
		MaxHTTPPipelineSize string `json:"max_http_pipeline_size"`
		MaxHTTPSessions     string `json:"max_http_sessions"`
	} `json:"replicator"`
	Stats struct {
		Rate    string `json:"rate"`
		Samples string `json:"samples"`
	} `json:"stats"`
	Uuids struct {
		Algorithm string `json:"algorithm"`
	} `json:"uuids"`
}

type ConfigSectionResponse struct {
	ErrorResponse
	AllowJsonp             string `json:"allow_jsonp"`
	AuthenticationHandlers string `json:"authentication_handlers"`
	BindAddress            string `json:"bind_address"`
	DefaultHandler         string `json:"default_handler"`
	EnableCors             string `json:"enable_cors"`
	LogMaxChunkSize        string `json:"log_max_chunk_size"`
	Port                   string `json:"port"`
	SecureRewrites         string `json:"secure_rewrites"`
	VhostGlobalHandlers    string `json:"vhost_global_handlers"`
}

type client struct {
	*CouchDb2ConnDetails
}

func (s *client) Meta() (res *MetaResponse, err error) {
	err = s.requester(http.MethodGet, "/", nil, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

func (s *client) ActiveTasks() (*ActiveTasksResponse, error) {
	var res ActiveTasksResponse
	err := s.requester(http.MethodGet, "/_active_tasks", nil, &res)

	//TODO ActiveTasks call can return an array of strings or a JSON error. We aren't handling error case

	return &res, err
}

func (s *client) AllDbs() (*AllDbsResponse, error) {
	var res AllDbsResponse
	err := s.requester(http.MethodGet, "/_all_dbs", nil, &res)

	//TODO AllDbs call can return an array of strings or a JSON error. We aren't handling error case

	return &res, err
}

func (s *client) DbUpdates(r *DbUpdatesRequest) (res *DbUpdatesResponse, err error) {
	requestBytes, err := json.Marshal(r)
	if err != nil {
		err = fmt.Errorf("Error creating request JSON: %s", err.Error())
		return
	}
	requestReader := bytes.NewReader(requestBytes)
	err = s.requester(http.MethodGet, "/_db_updates", requestReader, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

func (s *client) Membership() (*MembershipResponse, error) {
	var res MembershipResponse
	err := s.requester(http.MethodGet, "/_membership", nil, &res)

	if res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return &res, err
}

func (s *client) Log(r *LogRequest) (res *LogResponse, err error) {
	if s.Client == nil {
		return nil, errors.New("You must set an HTTP Client to make requests. Current client is nil")
	}

	requestBytes, err := json.Marshal(r)
	if err != nil {
		return
	}
	requestReader := bytes.NewReader(requestBytes)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/_log", s.Address), requestReader)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	httpRes, err := s.Client.Do(req)
	if err != nil {
		return
	}

	byt, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return
	}
	defer httpRes.Body.Close()

	temp := LogResponse(byt)
	res = &temp

	return
}

func (s *client) Replicate() (*ReplicateResponse, error) {
	var res ReplicateResponse
	err := s.requester(http.MethodGet, "/_replicate", nil, &res)

	if res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return &res, err
}

func (s *client) Restart() (*RestartResponse, error) {
	var res RestartResponse
	err := s.requester(http.MethodPost, "/_restart", nil, &res)

	if res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return &res, err
}

func (s *client) Stats() (res *StatsResponse, err error) {
	err = s.requester(http.MethodGet, "/_restart", nil, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

func (s *client) UUIDs(c uint8) (res *UUIDsResponse, err error) {
	err = s.requester(http.MethodGet, "/_restart", nil, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

func (s *client) Config() (res *ConfigResponse, err error) {
	err = s.requester(http.MethodGet, "/_config", nil, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

func (s *client) Section(se string) (res *ConfigSectionResponse, err error) {
	err = s.requester(http.MethodGet, fmt.Sprintf("/_config/%s", se), nil, &res)

	if res != nil && res.ErrorS != "" {
		return nil, &ErrorResponse{
			ErrorS: res.ErrorS,
			Reason: res.Reason,
		}
	}

	return
}

//func (c *server) Key(se, k string) (res *ConfigSectionResponse, err error) {
//	err = c.requester(http.MethodGet, fmt.Sprintf("/_config/%s/%s", se, k), nil, res)
// TODO Key function on server type
//	return
//}
//
//func (c *server) SetKey(se, k string) (res *ConfigSectionResponse, err error) {
//	err = c.requester("PUT", fmt.Sprintf("/_config/%s/%s", se, k), nil, res)
// TODO SetKey function on server type
//	return
//}
//
//func (c *server) DeleteKey(se, k string) (res *ConfigSectionResponse, err error) {
//	err = c.requester("DELETE", fmt.Sprintf("/_config/%s/%s", se, k), nil, res)
// TODO DeleteKey function on server type
//	return
//}

func NewClient(t time.Duration, addr string, user, pass string) (c Client) {
	c = &client{
		CouchDb2ConnDetails: NewConnection(t, addr, user, pass, true),
	}

	return
}
