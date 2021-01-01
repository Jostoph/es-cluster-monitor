package rest

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient() *http.Client {
	return &http.Client{}
}