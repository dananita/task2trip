package client

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

// KeepAliveTransport drains the remaining body from a response
// so that go will reuse the TCP connections.
// This is not enabled by default because there are servers where
// the response never gets closed and that would make the code hang forever.
// So instead it's provided as a http client middleware that can be used to override
// any request.
func KeepAliveTransport(rt http.RoundTripper) http.RoundTripper {
	return &keepAliveTransport{wrapped: rt}
}

type keepAliveTransport struct {
	wrapped http.RoundTripper
}

func (k *keepAliveTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := k.wrapped.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	resp.Body = &drainingReadCloser{rdr: resp.Body}
	return resp, nil
}

type drainingReadCloser struct {
	rdr     io.ReadCloser
	seenEOF uint32
}

func (d *drainingReadCloser) Read(p []byte) (n int, err error) {
	n, err = d.rdr.Read(p)
	if err == io.EOF || n == 0 {
		atomic.StoreUint32(&d.seenEOF, 1)
	}
	return
}

func (d *drainingReadCloser) Close() error {
	// drain buffer
	if atomic.LoadUint32(&d.seenEOF) != 1 {
		//#nosec
		if _, err := io.Copy(ioutil.Discard, d.rdr); err != nil {
			return err
		}
	}
	return d.rdr.Close()
}
