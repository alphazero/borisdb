package web

import (
	"bytes"
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
		return nil, fmt.Errorf("host is zerovalue")
	}

	c := &Client{
		hostport: fmt.Sprintf("%s:%d", host, port),
	}
	return c, nil
}

func (p *Client) Put(v []byte) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("nil value")
	}
	if len(v) == 0 {
		return nil, fmt.Errorf("value must be atleast 1 bytes.")
	}

	buf := bytes.NewReader(v)

	uri := fmt.Sprintf("http://%s/set", p.hostport)
	resp, e := http.Post(uri, mimetype, buf)
	if e != nil {
		return nil, fmt.Errorf("%s", e)
	}
	defer resp.Body.Close()

	err := responseErrorIfAny(resp)

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		err = fmt.Errorf("%s - with error:%s", e)
	}
	return body, err
}

func (p *Client) Get(key string) ([]byte, error) {
	uri := fmt.Sprintf("http://%s/get/%s", p.hostport, key)
	return p.httpGet(uri)
}

func (p *Client) Del(key string) ([]byte, error) {
	uri := fmt.Sprintf("http://%s/del/%s", p.hostport, key)
	return p.httpGet(uri)
}

func (p *Client) Info() ([]byte, error) {
	uri := fmt.Sprintf("http://%s/info", p.hostport)
	return p.httpGet(uri)
}

func (p *Client) Shutdown() ([]byte, error) {
	uri := fmt.Sprintf("http://%s/shutdown", p.hostport)
	return p.httpGet(uri)
}

/* util */

func (p *Client) httpGet(uri string) ([]byte, error) {
	resp, e := http.Get(uri)
	if e != nil {
		return nil, fmt.Errorf("%s", e)
	}
	defer resp.Body.Close()

	err := responseErrorIfAny(resp)

	// results
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		err = fmt.Errorf("%s - with error:%s", e)
	}
	return body, err
}

func responseErrorIfAny(resp *http.Response) error {
	if resp.StatusCode > 299 {
		return fmt.Errorf("%s", resp.Status)
	}
	return nil
}
