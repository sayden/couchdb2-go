package couchdb2_goclient

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var doDebug bool = false

type CouchDb2ConnDetails struct {
	Client   *http.Client
	Address  string
	Username string
	Password string
	protocol string
}

func (c *CouchDb2ConnDetails) GetConnection() *CouchDb2ConnDetails {
	return c
}

func (c *CouchDb2ConnDetails) bytesRequester(method, url string, reqBody io.Reader) (byt []byte, err error) {
	if c.Client == nil {
		return nil, errors.New("You must set an HTTP Client to make requests. Current client is nil")
	}

	if doDebug {
		fmt.Printf("%s://%s%s\n", c.protocol, c.Address, url)
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s%s", c.protocol, c.Address, url), reqBody)
	if err != nil {
		return nil, err
	}

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	httpRes, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpRes.Body.Close()
	byt, err = ioutil.ReadAll(httpRes.Body)
	return byt, err
}

func (c *CouchDb2ConnDetails) requester(method, url string, reqBody io.Reader, res interface{}) error {
	if c.Client == nil {
		return errors.New("You must set an HTTP Client to make requests. Current client is nil")
	}

	if doDebug {
		fmt.Printf("%s://%s%s\n", c.protocol, c.Address, url)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s://%s%s", c.protocol, c.Address, url), reqBody)
	if err != nil {
		return err
	}

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	httpRes, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if doDebug {
		var byt []byte
		byt, err = ioutil.ReadAll(httpRes.Body)
		defer httpRes.Body.Close()
		if err != nil {
			return err
		}

		//fmt.Printf("a %#v\n", string(byt))

		json.Unmarshal(byt, &res)
		if err != nil {
			return fmt.Errorf("Error parsing response: %s", err.Error())
		}

		//fmt.Printf("c %#v\n", res)

	} else {

		err = json.NewDecoder(httpRes.Body).Decode(&res)
		if err != nil {
			return fmt.Errorf("Error parsing response: %s", err.Error())
		}
	}

	return nil
}

func NewConnection(timeout time.Duration, addr string, user, pass string, secure bool) (conn *CouchDb2ConnDetails) {
	conn = &CouchDb2ConnDetails{
		Client: &http.Client{
			Timeout: timeout,
		},
		Address:  addr,
		Username: user,
		Password: pass,
		protocol: "http",
	}

	if secure {
		conn.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		conn.protocol = "https"
	}

	return conn
}
