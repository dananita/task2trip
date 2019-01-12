package client

import (
	"context"
	"net/http"
	"net/http/httptest"
)

type Do interface {
	Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error)
}

func NewDO(h http.Handler) Do {
	return &doWithServeHTTP{h: h}
}

type doWithServeHTTP struct {
	h http.Handler
}

func (d *doWithServeHTTP) Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	d.h.ServeHTTP(resp, req.WithContext(ctx))
	return resp.Result(), nil
}
