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
	ImpactPeriod     map[string]interface{}   `json:"period"`
	ImpactScore      map[string]interface{}   `json:"impactScore"`
	ActiveDays       map[string]interface{}   `json:"activeDays"`
	CommitsPerDay    map[string]interface{}   `json:"commitsPerDay"`
	MostChurnedFiles []map[string]interface{} `json:"mostChurnedFiles"`
	NewCodePercent   map[string]interface{}   `json:"newCodePercentage"`
	ChurnPercent     map[string]interface{}   `json:"churnPercentage"`
	LegacyPercent    map[string]interface{}   `json:"legacyPercentage"`
	FileChanged      map[string]interface{}   `json:"fileChanged"`
	InsertionPoints  map[string]interface{}   `json:"insertionPoints"`
	Additions        map[string]interface{}   `json:"additions"`
	Deletions        map[string]interface{}   `json:"deletions"`
}

func TestGetWeeklyImpact(t *testing.T) {
	// from := int64(1592784000) // 2020-06-22
	from := int64(1592870400) // 2020-06-23
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+endpoint+"?date_from="+strconv.FormatInt(from, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	resp := TestWeeklyData{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	log.Println("### Response GetWeeklyImpact: ", resp)
	assert.NotNil(t, resp.ImpactPeriod["date_from"])
	assert.NotNil(t, resp.ImpactPeriod["date_to"])
	assert.NotNil(t, resp.ImpactScore)
	assert.Equal(t, resp.ActiveDays["currentPeriod"], 1.0)
	assert.Equal(t, resp.ActiveDays["previousPeriod"], 1.0)
	assert.Equal(t, resp.CommitsPerDay["currentPeriod"], 1.0)
	assert.Equal(t, resp.MostChurnedFiles[0]["fileName"], "file2")
	assert.Equal(t, resp.NewCodePercent["currentPeriod"], 71.42857142857143)
	assert.Equal(t, resp.NewCodePercent["previousPeriod"], 25.0)
	assert.Equal(t, resp.ChurnPercent["currentPeriod"], 28.57142857142857)
	assert.Equal(t, resp.ChurnPercent["previousPeriod"], 75.0)
	assert.Equal(t, resp.LegacyPercent["currentPeriod"], 17.33)
	assert.Equal(t, resp.LegacyPercent["previousPeriod"], 23.44)
	assert.Equal(t, resp.FileChanged["currentPeriod"], 3.0)
	assert.Equal(t, resp.FileChanged["previousPeriod"], 1.0)
	assert.Equal(t, resp.InsertionPoints["currentPeriod"], 14.0)
	assert.Equal(t, resp.InsertionPoints["previousPeriod"], 3.0)
	assert.Equal(t, resp.Additions["currentPeriod"], 100.0)
	assert.Equal(t, resp.Additions["previousPeriod"], 32.0)
	assert.Equal(t, resp.Deletions["currentPeriod"], 90.0)
	assert.Equal(t, resp.Deletions["previousPeriod"], 20.0)
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
