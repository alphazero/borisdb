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

func (p *Client) Put(v []byte) (string, error) {
	if v == nil {
		return "", fmt.Errorf("nil value")
	}
	if len(v) == 0 {
		return "", fmt.Errorf("value must be atleast 1 bytes.")
	}

	buf := bytes.NewReader(v)

	uri := fmt.Sprintf("http://%s/set", p.hostport)
	resp, e := http.Post(uri, mimetype, buf)
	if e != nil {
		return "", fmt.Errorf("%s", e)
	}
	defer resp.Body.Close()

	err := responseErrorIfAny(resp)

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		err = fmt.Errorf("%s - with error:%s", e)
	}
	return string(body), err
}

func responseErrorIfAny(resp *http.Response) error {
	if resp.StatusCode > 299 {
		return fmt.Errorf("%s", resp.Status)
	}
	return nil
}

func (p *Client) Get(key string) ([]byte, error) {
	// service request
	uri := fmt.Sprintf("http://%s/get/%s", p.hostport, key)
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
