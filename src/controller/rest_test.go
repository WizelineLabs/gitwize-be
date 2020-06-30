package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitwize-be/src/configuration"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

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

var router *gin.Engine

func init() {
	configuration.ReadConfiguration()
	router = Initialize()
}

func TestPostReposOK(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	posRequest := RepoInfoPost{
		Name: "Gitwize",
		Url:  "https://github.com/wizeline/gitwize-be",
		//User:     "tester",
		//Password: "",
	}
	expectedResult := "{\"id\":\\d+," +
		"\"name\":\"Gitwize\"," +
		"\"url\":\"https://github.com/wizeline/gitwize-be\"," +
		"\"status\":\"ONGOING\"," +
		"\"branches\":\\[.*\\]," +
		"\"last_updated\":\"[0-9:ZT\\+\\-\\.]+\"}"

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwRepoPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())
}

func TestPostRepos_BadRequest(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	posRequest := RepoInfoPost{
		Name: "Gitwize",
		//User: "tester",
	}

	b, err := json.Marshal(posRequest)
	if err != nil {
		t.Error(err.Error())
	}
	w := performRequest(router, http.MethodPost, gwRepoPost, bytes.NewReader(b))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRepo_Unauthorized(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "false"

	w := performRequest(router, http.MethodGet, gwEndPointRepository+"1", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetListRepos_PutRepos_GetRepos_GetStats_DelRepos_OK(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	repoInfos := make([]RepoInfoGet, 0)
	err := json.Unmarshal(w.Body.Bytes(), &repoInfos)
	if err != nil {
		t.Error(err.Error())
	}
	if len(repoInfos) > 0 {
		// TEST PutRepo
		updateRepo := repoInfos[len(repoInfos)-1]
		updateRepo.Status = "UPDATED"
		b, err := json.Marshal(updateRepo)
		if err != nil {
			t.Error(err.Error())
		}
		w = performRequest(router, http.MethodPut, gwEndPointRepository+strconv.Itoa(repoInfos[len(repoInfos)-1].ID), bytes.NewReader(b))
		assert.Equal(t, http.StatusOK, w.Code)

		// TEST GetRepo
		w := performRequest(router, http.MethodGet, gwEndPointRepository+strconv.Itoa(repoInfos[len(repoInfos)-1].ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)
		expectedResult := "{\"id\":\\d+," +
			"\"name\":\".*\"," +
			"\"url\":\".*\"," +
			"\"status\":\"UPDATED\"," +
			"\"branches\":\\[.*\\]," +
			"\"last_updated\":\"[0-9:ZT\\+\\-\\.]+\"}"
		assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())

		// TEST Stats
		to := time.Now().Unix()
		from := to - 7*24*3600
		w = performRequest(router, http.MethodGet, gwEndPointRepository+
			strconv.Itoa(repoInfos[0].ID)+"/stats?date_from="+strconv.FormatInt(from, 10)+
			"&date_to="+strconv.FormatInt(to, 10), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		expectedResult = "{\"id\":\\d+," +
			"\"name\":\".*\"," +
			"\"url\":\".*\"," +
			"\"status\":\".*\"," +
			"\"metric\":{\"commits\":\\[.*\\],\"lines_added\":\\[.*\\],\"lines_removed\":\\[.*\\]," +
			"\"loc\":\\[.*\\],\"prs_created\":\\[.*\\],\"prs_merged\":\\[.*\\],\"prs_rejected\":\\[.*\\]}}"
		assert.Regexp(t, regexp.MustCompile(expectedResult), w.Body.String())

		// TEST DelRepo
		w = performRequest(router, http.MethodDelete, gwEndPointRepository+strconv.Itoa(repoInfos[len(repoInfos)-1].ID), nil)
		assert.Equal(t, http.StatusNoContent, w.Code)
	}
}
