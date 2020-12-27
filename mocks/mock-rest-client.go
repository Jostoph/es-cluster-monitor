package mocks

import "net/http"

type Client struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var GetDoFunc func(req *http.Request) (*http.Response, error)

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}