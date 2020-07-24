package controller

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"gitwize-be/src/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

var endpoint = strings.TrimLeft(gwWeeklyImpact, "id:")

type TestWeeklyData struct {
	ImpactPeriod     map[string]interface{} `json:"period"`
	ImpactScore      map[string]interface{} `json:"impactScore"`
	ActiveDays       map[string]interface{} `json:"activeDays"`
	CommitsPerDay    map[string]interface{} `json:"commitsPerDay"`
	MostChurnedFiles []interface{}          `json:"mostChurnedFiles"`
	NewCodePercent   map[string]interface{} `json:"newCodePercentage"`
	ChurnPercent     map[string]interface{} `json:"churnPercentage"`
}

func TestGetWeeklyImpact(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+endpoint,
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	resp := TestWeeklyData{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	log.Println("### Response GetWeeklyImpact: ", resp)
	assert.NotNil(t, resp.ImpactPeriod["date_from"])
	assert.NotNil(t, resp.ImpactPeriod["date_to"])
	assert.NotNil(t, resp.ImpactScore)
	assert.NotNil(t, resp.ActiveDays)
	assert.NotNil(t, resp.CommitsPerDay)
	assert.NotNil(t, resp.MostChurnedFiles)
	assert.NotNil(t, resp.NewCodePercent)
	assert.NotNil(t, resp.ChurnPercent)
}

func TestGetWeeklyImpactNotFound(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(9999)+endpoint,
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetNewCodeAndChurnPercentage(t *testing.T) {
	stat := db.ModificationStat{Modifications: 5, Additions: 10, Deletions: 2}
	newCodePercent, churnPercent := getNewCodeAndChurnPercentage(stat)
	assert.Equal(t, float64(10)/15*100, newCodePercent)
	assert.Equal(t, float64(5)/15*100, churnPercent)
	stat = db.ModificationStat{Modifications: 0, Additions: 0, Deletions: 0}
	newCodePercent, churnPercent = getNewCodeAndChurnPercentage(stat)
	assert.Equal(t, 0.0, newCodePercent)
	assert.Equal(t, 0.0, churnPercent)
}
