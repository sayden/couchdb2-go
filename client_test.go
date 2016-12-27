package couchdb2_goclient

import (
	"testing"
	"time"
	"strings"
	"os"
)

var testClient Client

func newClient() Client {
	if testClient == nil {
		testClient = NewClient(time.Second, os.Getenv("COUCHDB_ADDRESS"), os.Getenv("COUCHDB_USER"), os.Getenv("COUCHDB_PASSWORD"))
	}
	return testClient
}

var printDebug bool

func debug(i interface{}, t *testing.T) {
	if printDebug {
		t.Logf("%#v", i)
	}
}

func TestClient_Meta_And_Transport(t *testing.T) {
	t.Run("Reached correctly", func(t *testing.T) {
		c := newClient()

		res, err := c.Meta()
		if err != nil {
			t.Fatal(err)
		}

		if res.CouchDB == "" || res.Version == "" {
			t.Errorf("%#v", res)
		}
	})

	t.Run("Incorrect ip. Timeout", func(t *testing.T) {
		c := NewClient(time.Second, "0.0.0.1:213132", "", "")

		_, err := c.Meta()
		if err == nil {
			t.Fail()
		}
	})

	t.Run("Incorrect port. Timeout", func(t *testing.T) {
		c := NewClient(time.Second, "127.0.0.1:333333" , "", "")

		_, err := c.Meta()
		if err == nil {
			t.Fail()
		}
	})
}

func TestClient_ActiveTasks(t *testing.T) {
	c := newClient()

	res, err := c.ActiveTasks()
	if err != nil {
		t.Error(err)
	}

	if len(*res) > 0 {
		if (*res)[1].Type == "" {
			t.Fail()
		}
	}
}

func TestClient_AllDbs(t *testing.T) {
	c := newClient()

	res, err := c.AllDbs()
	if err != nil {
		t.Error(err)
	}

	if len(*res) > 0 {
		if (*res)[0] == "" {
			t.Fail()
		}
	}

	debug(res, t)
}

func TestClient_DbUpdates(t *testing.T) {
	//TODO
	c := newClient()

	req := &DbUpdatesRequest{
		Timeout:   1,
		Feed:      "",
		HeartBeat: true,
	}
	_, err := c.DbUpdates(req)
	if err == nil {
		t.Fail()
	}
}

func TestClient_Membership(t *testing.T) {
	//TODO
	c := newClient()

	res, err := c.Membership()
	if err != nil {
		t.Fatal(err)
	}

	if len(res.ClusterNodes) > 0 {
		if res.ClusterNodes[0] == "" {
			t.Fail()
		}
	}

	debug(res, t)
}

func TestClient_Log(t *testing.T) {
	c := newClient()

	req := &LogRequest{
		Bytes:  100,
		Offset: 0,
	}

	res, err := c.Log(req)
	if err != nil {
		t.Fatal(err)
	}

	if (*res) == "" {
		t.Fail()
	}

	debug(res, t)
}

func TestClient_Replicate(t *testing.T) {
	//TODO
	c := newClient()

	_, err := c.Replicate()
	if err == nil {
		t.Fatal()
	}
}

func TestClient_Restart(t *testing.T) {
	c := newClient()

	//No database is created so it must return an error
	res, err := c.Restart()
	if err == nil {
		t.Fail()
	}

	debug(res, t)
}

func TestClient_Stats(t *testing.T) {
	c := newClient()

	res, err := c.Stats()
	if err == nil {
		t.Fatal()
	}

	debug(res, t)
}

func TestClient_UUIDs(t *testing.T) {
	c := newClient()

	res, err := c.UUIDs(10)
	if err == nil {
		t.Fatal(err)
	}

	if !strings.Contains(err.Error(), "atabase does not"){
		t.Fail()
	}

	debug(res, t)
}

func TestClient_Config(t *testing.T) {
	c := newClient()

	res, err := c.Config()
	if err == nil {
		t.Fatal(err)
	}

	if !strings.Contains(err.Error(), "atabase does not"){
		t.Fail()
	}

	debug(res, t)
}

func TestClient_Section(t *testing.T) {
	c := newClient()

	res, err := c.Section("")
	if err == nil {
		t.Fatal(err)
	}

	if !strings.Contains(err.Error(), "atabase does not"){
		t.Fail()
	}

	debug(res, t)
}
