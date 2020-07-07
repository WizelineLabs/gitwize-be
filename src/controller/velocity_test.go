package controller

import (
	"encoding/json"
	"gitwize-be/src/configuration"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestCodeVelocity struct {
	Commits    map[string]string `json:"commits"`
	NetChanges map[string]string `json:"netChanges"`
	//NewCodeChanges map[string]string `json:"newCodeChanges"`
}

func Test_CodeVelocity_NetChanges(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"

	from := time.Date(2020, 5, 1, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.Local).Unix()
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(2)+"/code-velocity?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)

	// must contain data of June
	resp := TestCodeVelocity{
		Commits:    make(map[string]string),
		NetChanges: make(map[string]string),
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	log.Printf("### Response entity: %s", resp)
	assert.NotNil(t, resp.NetChanges["June"])
	assert.NotNil(t, resp.Commits["June"])
}
