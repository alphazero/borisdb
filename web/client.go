package web

import (
	"bytes"
	//	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

const mimetype = "application/binary"

type Client struct {
	hostport string
	setUri   string
}

func NewClient(host string, port int) (*Client, error) {
	if host == "" {
		return nil, fmt.Errorf("err - host is zerovalue")
	}

	c := &Client{
		hostport: fmt.Sprintf("%s:%d", host, port),
	}
	return c, nil
}

func (p *Client) Put(v []byte) (string, error) {
	if v == nil {
		return "", fmt.Errorf("err - nil value")
	}
	if len(v) == 0 {
		return "", fmt.Errorf("err - value must be atleast 1 bytes.")
	}

	buf := bytes.NewReader(v)

	uri := fmt.Sprintf("http://%s/set", p.hostport)
	resp, e := http.Post(uri, mimetype, buf)
	if e != nil {
		return "", fmt.Errorf("err - %s\n", e)
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return "", fmt.Errorf("err - %s\n", e)
	}
	return string(body), nil
}

func (p *Client) Get(key string) ([]byte, error) {
	/*
		// encode key
		oid, e := hex.DecodeString(key)
		if e != nil {
			return nil, fmt.Errorf("err - invalid key - %s", e)
		}
	*/
	// service request
	uri := fmt.Sprintf("http://%s/get/%s", p.hostport, key)
	resp, e := http.Get(uri)
	if e != nil {
		return nil, fmt.Errorf("err - %s\n", e)
	}
	defer resp.Body.Close()

	// results
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, fmt.Errorf("err - %s\n", e)
	}

	return body, nil
}
