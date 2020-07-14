package controller

import (
	"github.com/stretchr/testify/assert"
	"gitwize-be/src/configuration"
	"net/http"
	"strconv"
	"testing"
)

func TestGetWeeklyImpact(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	from := int64(0)        // 1970
	to := int64(2524637115) // 2050
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/weekly-impact?date_from="+strconv.FormatInt(from, 10)+
		"&date_to="+strconv.FormatInt(to, 10),
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetWeeklyImpactBadRequest(t *testing.T) {
	configuration.CurConfiguration.Auth.AuthDisable = "true"
	w := performRequest(router, http.MethodGet, gwEndPointRepository+
		strconv.Itoa(1)+"/weekly-impact?date_from=abc&date_to=xyz",
		nil, header{Key: "AuthenticatedUser", Value: "tester@wizeline.com"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
