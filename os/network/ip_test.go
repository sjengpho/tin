package network

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func fakeReadAllError(r io.Reader) ([]byte, error) {
	return nil, errors.New("read error")
}

func TestNewPublicIPLookup(t *testing.T) {
	want := &ipLookupper{}
	got := NewPublicIPLookup()

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("want %v, got %v", reflect.TypeOf(want), reflect.TypeOf(got))
	}
}

func TestIPLookupSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("127.0.0.1"))
	}))
	defer func() { testServer.Close() }()

	lookupper := ipLookupper{
		timeout: 10 * time.Second,
		sources: []string{testServer.URL},
	}

	_, err := lookupper.Lookup()
	if err != nil {
		t.Errorf("want %v, got %v", nil, err)
	}
}

func TestIPLookupResponseError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	}))
	defer func() { testServer.Close() }()

	lookupper := ipLookupper{
		timeout: time.Millisecond,
		sources: []string{testServer.URL},
	}

	want := ErrTimeout
	_, got := lookupper.Lookup()
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestIPLookupReadBodyError(t *testing.T) {
	readAll = fakeReadAllError
	defer func() { readAll = ioutil.ReadAll }()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()

	lookupper := ipLookupper{
		timeout: time.Millisecond,
		sources: []string{testServer.URL},
	}

	want := ErrTimeout
	_, got := lookupper.Lookup()
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestIPLookupParseError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("invalid-ip-address"))
	}))
	defer func() { testServer.Close() }()

	lookupper := ipLookupper{
		timeout: time.Millisecond,
		sources: []string{testServer.URL},
	}

	want := ErrTimeout
	_, got := lookupper.Lookup()
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
