package sprintly

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gorilla/schema"
	"github.com/kr/pretty"
)

// setup should be used to set up the testing environment.
// It returns a testing client, server and a mux associated with that server.
// The client is configured to send requests to the returned HTTP server.
// Let the end-to-end testing begin!
//
// The server should be closed at the end, so what you probably want to do is
//
//    client, server, mux := setup()
//    defer server.Close()
//
func setup() (*Client, *httptest.Server, *http.ServeMux) {
	// Set up the testing server.
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// Set up the testing client.
	client := NewClient("krtecek", "secret")
	client.SetBaseURL(server.URL)

	return client, server, mux
}

func ensureMethod(t *testing.T, r *http.Request, method string) {
	if m := r.Method; m != method {
		t.Errorf("Request method is %v, want %v", m, method)
	}
}

func ensureEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Objects not equal!\n\ngot = %# v\n\nwant = %# v \n",
			pretty.Formatter(got), pretty.Formatter(want))
	}
}

func decodeArgs(dst interface{}, r io.Reader) error {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(content))
	if err != nil {
		return err
	}

	return schema.NewDecoder().Decode(dst, values)
}
