package controller

import (
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_PullRequestSize_OK(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"

	from := time.Date(2020, 5, 5, 0, 0, 0, 0, time.Local).Unix()
	to := time.Date(2020, 5, 12, 23, 59, 59, 0, time.Local).Unix()
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/pullrequest-size?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)

	expected := "{" +
		"\"2020-05-05\":[]," +
		"\"2020-05-06\":[]," +
		"\"2020-05-07\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-08\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-09\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-10\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-11\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"merged\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-12\":[]}"
	assert.Equal(t, expected, w.Body.String())

	to = time.Date(2020, 5, 10, 23, 59, 59, 0, time.Local).Unix()
	w = performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/pullrequest-size?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)
	expected = "{" +
		"\"2020-05-05\":[]," +
		"\"2020-05-06\":[]," +
		"\"2020-05-07\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-08\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-09\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]," +
		"\"2020-05-10\":[{\"title\":\"GWZ-23 verifies access token\",\"size\":120,\"status\":\"opened\",\"review_time\":89," +
		"\"url\":\"https://github.com/wizeline/gitwize-be/pull/1\",\"created_date\":\"2020-05-07\",\"created_by\":\"ltvcuong\"}]}"
	assert.Equal(t, expected, w.Body.String())
}
