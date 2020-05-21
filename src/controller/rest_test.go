package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
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
	os.Setenv("AUTH_DISABLED", "true")
	router := Initialize()
	posRequest := RepoInfoPost{
		Name:   "Gitwize",
		Url:    "https://github.com/gitwize",
		Status: "ONGOING",
		User:   "tester",
	}
	expectedResult := "{\"id\":\\d+}"

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwEndPointPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())
}

func TestPostRepos_BadRequest(t *testing.T) {
	os.Setenv("AUTH_DISABLED", "true")
	router := Initialize()
	posRequest := RepoInfoPost{
		Name: "Gitwize",
		Url:  "https://github.com/gitwize",
		User: "tester",
	}

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwEndPointPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRepo_Unauthorized(t *testing.T) {
	os.Setenv("AUTH_DISABLED", "false")
	router := Initialize()

	w := performRequest(router, http.MethodGet, "/api/v1/repositories/1", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
