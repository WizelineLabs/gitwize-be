package controller

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

var endpoint = strings.TrimLeft(gwWeeklyImpact, "id:")

type TestWeeklyData struct {
	ImpactPeriod    map[string]interface{} `json:"period"`
	ImpactScore     map[string]interface{} `json:"impactScore"`
	ActiveDays      map[string]interface{} `json:"activeDays"`
	CommitsPerDay   map[string]interface{} `json:"commitsPerDay"`
	MostChurnedFile map[string]interface{} `json:"mostChurnedFile"`
}

func TestGetWeeklyImpact(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+endpoint,
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)

	resp := TestWeeklyData{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	log.Println("### Response entity: ", resp)
	assert.NotNil(t, resp.ImpactPeriod["date_from"])
	assert.NotNil(t, resp.ImpactPeriod["date_to"])
	assert.Equal(t, float64(184), resp.ImpactScore["currentPeriod"])
	assert.Equal(t, float64(10), resp.ImpactScore["previousPeriod"])
}

func TestGetWeeklyImpactNotFound(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(9999)+endpoint,
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusNotFound, w.Code)
}
