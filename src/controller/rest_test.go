package controller

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPostReposOK(t *testing.T) {
	router := Initialize()
	posRequest := RepoInfoPost{
		Name:   "Gitwize",
		Url:    "https://github.com/gitwize",
		Status: "ONGOING",
		User:   "tester",
	}
	expectedResult := "{\"id\": 1}"

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwEndPointPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedResult, w.Body.String())
}

func TestPostRepos_BadRequest(t *testing.T) {
	router := Initialize()
	posRequest := RepoInfoPost{
		Name:   "Gitwize",
		Url:    "https://github.com/gitwize",
		Status: "ONGOING",
		User:   "tester",
	}

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwEndPointPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
